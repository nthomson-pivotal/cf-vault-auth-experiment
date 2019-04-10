# CloudFoundry Auth Method

The `cf` auth method can be used to authenticate with Vault using
CloudFoundry instance identity certificates. This method of authentication makes it easy to
introduce a Vault token into a CloudFoundry application.

## Authentication

### Via the API

The default endpoint is `auth/cf/login`. If this auth method was enabled
at a different path, use that value instead of `cf`.

```shell
$ curl \
    --request POST \
    --data '{"jwt": "your_service_account_jwt", "role": "demo"}' \
    http://127.0.0.1:8200/v1/auth/cf/login
```

The response will contain a token at `auth.client_token`:

```json
{
  "auth": {
    "client_token": "38fe9691-e623-7238-f618-c94d4e7bc674",
    "accessor": "78e87a38-84ed-2692-538f-ca8b9f400ab3",
    "policies": [
      "default"
    ],
    "lease_duration": 2764800,
    "renewable": true
  }
}
```

## Configuration

Auth methods must be configured in advance before applications can
authenticate. These steps are usually completed by an operator or configuration
management tool.


1. Enable the CloudFoundry auth method:

    ```text
    $ vault auth enable cf
    ```

1. Use the `/certs` endpoint to configure Vault to accept instance identity
certificates signed by your CloudFoundry foundation by adding the Diego root
CA certificate public key.

    ```text
    $ vault write auth/cf/certs/default \
        display_name=default \
        certificate=@diego-ca-cert.pem
    ```

1. Map an application GUID to an existing Vault policy by name:

    ```text
    vault write auth/cf/map/apps/5f4e225f-87cb-4f29-a21a-28b72717168d \
        value=admin-policy
    ```

    This role authorizes the given application GUID to generate token with 
    the policy Vault policy `admin-policy`

    For the complete list of configuration options, please see the API
    documentation.