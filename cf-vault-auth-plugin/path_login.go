package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathAuthLogin(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	encodedCertificate := d.Get("certificate").(string)
	encodedKey := d.Get("key").(string)

	b.Logger().Info("got login request")

	certs, _ := parseCertificate(encodedCertificate)

	b.Logger().Info("cert parsed")

	if len(certs) < 2 {
		return nil, errors.New("Login request must include container intermediate CA")
	}

	cert := certs[0]
	intermediate := certs[1]

	trustedCerts := b.loadTrustedCerts(ctx, req.Storage)

	err := b.validateConnState(trustedCerts, cert, intermediate)
	if err != nil {
		return nil, err
	}

	priv, err := parseKey(encodedKey)
	if err != nil {
		return nil, err
	}

	_, err = verifyCertificateAndKey(cert, priv)
	if err != nil {
		return nil, err
	}

	if len(cert.Subject.OrganizationalUnit) < 3 {
		return nil, errors.New("malformed OU field")
	}

	orgGUID := strings.Split(cert.Subject.OrganizationalUnit[0], ":")[1]
	spaceGUID := strings.Split(cert.Subject.OrganizationalUnit[1], ":")[1]
	appGUID := strings.Split(cert.Subject.OrganizationalUnit[2], ":")[1]

	appPolicies, err := b.AppsMap.Policies(ctx, req.Storage, appGUID)
	if err != nil {
		return nil, errors.New("app policies")
	}

	spacePolicies, err := b.SpacesMap.Policies(ctx, req.Storage, spaceGUID)
	if err != nil {
		return nil, errors.New("space policies")
	}

	orgPolicies, err := b.OrgsMap.Policies(ctx, req.Storage, orgGUID)
	if err != nil {
		return nil, errors.New("org policies")
	}

	policies := append(appPolicies, spacePolicies...)
	policies = append(policies, orgPolicies...)

	// Unique, since we want to remove duplicates and that will cause errors when
	// we compare policies later.
	uniq := map[string]struct{}{}
	for _, v := range policies {
		if _, ok := uniq[v]; !ok {
			uniq[v] = struct{}{}
		}
	}
	newPolicies := make([]string, 0, len(uniq))
	for k := range uniq {
		newPolicies = append(newPolicies, k)
	}
	policies = newPolicies

	return &logical.Response{
		Auth: &logical.Auth{
			Policies: policies,
			Metadata: map[string]string{
				"orgGuid":   orgGUID,
				"spaceGuid": spaceGUID,
				"appGuid":   appGUID,
			},
			LeaseOptions: logical.LeaseOptions{
				// Arbitrarily set the TTL to 5 minute (todo: allow customization)
				TTL: 5 * time.Minute,
				// Set the max TTL to the expiry time of the certificate we received
				MaxTTL:    time.Until(cert.NotAfter),
				Renewable: true,
			},
		},
	}, nil
}

func (b *backend) pathAuthRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	b.Logger().Info("got renew request")

	if req.Auth == nil {
		return nil, errors.New("request auth was nil")
	}

	b.Logger().Info("renew auth", "req", req.Auth, "maxTTL", req.Auth.MaxTTL)

	return framework.LeaseExtend(30*time.Second, req.Auth.MaxTTL, b.System())(ctx, req, d)
}

func parseCertificate(encoded string) ([]*x509.Certificate, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, errors.New("certificate decode error: " + err.Error())
	}

	certs := parsePEM(decoded)

	return certs, nil
}

// parsePEM parses a PEM encoded x509 certificate
func parsePEM(raw []byte) (certs []*x509.Certificate) {
	for len(raw) > 0 {
		var block *pem.Block
		block, raw = pem.Decode(raw)
		if block == nil {
			break
		}
		if (block.Type != "CERTIFICATE" && block.Type != "TRUSTED CERTIFICATE") || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}
		certs = append(certs, cert)
	}
	return
}

func parseKey(encoded string) (*rsa.PrivateKey, error) {
	var block *pem.Block

	decodedKey, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, errors.New("key decode error: " + err.Error())
	}

	block, _ = pem.Decode(decodedKey)
	if block == nil {
		return nil, nil
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func verifyCertificateAndKey(cert *x509.Certificate, priv *rsa.PrivateKey) (bool, error) {
	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		if pub.N.Cmp(priv.N) != 0 {
			return false, errors.New("tls: private key does not match public key")
		}
	default:
		return false, errors.New("tls: unknown public key algorithm")
	}

	return true, nil
}

func (b *backend) validateConnState(roots *x509.CertPool, cert *x509.Certificate, intermediate *x509.Certificate) error {
	pool := x509.NewCertPool()
	b.Logger().Info("adding cert to pool")
	pool.AddCert(intermediate)

	b.Logger().Info("created intermediate pool")

	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: pool,
		//KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	b.Logger().Info("created opts")

	chains, err := cert.Verify(opts)
	if err != nil {
		return err
	}

	b.Logger().Info("verify returned")

	if len(chains) > 0 {
		return nil
	}

	return errors.New("No certificate validation match")
}

func (b *backend) loadTrustedCerts(ctx context.Context, storage logical.Storage) (pool *x509.CertPool) {
	pool = x509.NewCertPool()

	names, err := storage.List(ctx, "cert/")
	if err != nil {
		b.Logger().Error("failed to list trusted certs", "error", err)
		return
	}
	for _, name := range names {
		b.Logger().Info("processing entry", "name", name)
		entry, err := b.Cert(ctx, storage, strings.TrimPrefix(name, "cert/"))
		if err != nil {
			b.Logger().Error("failed to load trusted cert", "name", name, "error", err)
			continue
		}
		parsed := parsePEM([]byte(entry.Certificate))
		if len(parsed) == 0 {
			b.Logger().Error("failed to parse certificate", "name", name)
			continue
		}
		if parsed[0].IsCA {
			b.Logger().Info("is ca")
			for _, p := range parsed {
				b.Logger().Error("adding cert to pool")
				pool.AddCert(p)
			}
		}
	}
	return
}
