package main

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/hashicorp/vault/logical/plugin"
)

func main() {
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)

	if err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: Factory,
		TLSProviderFunc:    tlsProviderFunc,
	}); err != nil {
		log.Fatal(err)
	}
}

func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(c)
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}
	return b, nil
}

type backend struct {
	*framework.Backend

	AppsMap   *framework.PolicyMap
	SpacesMap *framework.PolicyMap
	OrgsMap   *framework.PolicyMap
}

func Backend(c *logical.BackendConfig) *backend {
	var b backend

	b.AppsMap = &framework.PolicyMap{
		PathMap: framework.PathMap{
			Name: "apps",
			Schema: map[string]*framework.FieldSchema{

				"value": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Policies for the app GUID.",
				},
			},
		},
		DefaultKey: "default",
	}

	b.SpacesMap = &framework.PolicyMap{
		PathMap: framework.PathMap{
			Name: "spaces",
			Schema: map[string]*framework.FieldSchema{

				"value": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Policies for the space GUID.",
				},
			},
		},
		DefaultKey: "default",
	}

	b.OrgsMap = &framework.PolicyMap{
		PathMap: framework.PathMap{
			Name: "orgs",
			Schema: map[string]*framework.FieldSchema{

				"value": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Policies for the org GUID.",
				},
			},
		},
		DefaultKey: "default",
	}

	b.Backend = &framework.Backend{
		BackendType: logical.TypeCredential,
		AuthRenew:   b.pathAuthRenew,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{"login"},
		},
		Paths: func() []*framework.Path {
			var paths []*framework.Path

			paths = append(paths, &framework.Path{
				Pattern: "login",
				Fields: map[string]*framework.FieldSchema{
					"certificate": &framework.FieldSchema{
						Type: framework.TypeString,
					},
					"key": &framework.FieldSchema{
						Type: framework.TypeString,
					},
				},
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.pathAuthLogin,
				},
			})

			paths = append(paths, b.AppsMap.Paths()...)
			paths = append(paths, b.SpacesMap.Paths()...)
			paths = append(paths, b.OrgsMap.Paths()...)
			paths = append(paths, pathListCerts(&b), pathCerts(&b))

			return paths
		}(),
	}

	return &b
}
