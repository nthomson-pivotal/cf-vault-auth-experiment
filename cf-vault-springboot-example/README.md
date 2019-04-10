# Spring Boot CloudFoundry Vault Example

This project demonstrates an example Spring Boot application that obtains Vault tokens
via the `cf-vault-auth-plugin` mechanism. It does this by providing the CloudFoundry instance
identify certificate and private key as authentication credentials.

## Running in CloudFoundry

As prerequisites to running this example you must have:

1. A running Vault server with the `cf-vault-auth-plugin` installed
1. The Vault K/V backend mounted at `kv`
1. Added the Diego root CA of the CloudFoundry foundation you wish to run this in
to the Vault auth plugin
1. Installed the `cf-vault-cli` CF CLI plugin

First, write a value to the `kv` backend:

```
vault kv put kv/github github.oauth2.key=foobar
```

Build the application using Maven:

```
./mvnw package
```

Create a `manifest.yml` file that looks something like this:

```
---
applications:
- name: cf-vault-example
  memory: 1G
  random-route: true
  path: target/cf-vault-springboot-example-0.0.1-SNAPSHOT.jar
  env:
    SPRING_CLOUD_VAULT_SCHEME: http
    SPRING_CLOUD_VAULT_HOST: <your vault IP/host>
```

Now `push` the application and map its GUID to a Vault role:

```
cf push -i 0

cf vault apps set -a cf-vault-example -p some-vault-policy-name

cf scale cf-vault-example -i 1
```

If you access the CF route mapped to the app you should now see the value of your secret, 
securely displayed in clear text.
