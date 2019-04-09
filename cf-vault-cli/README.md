# CF CLI Vault Auth Plugin

This is a plugin the CloudFoundry CLI that is intended to provide convenience functions for 
mapping orgs, spaces and applications to Vault policies via the Vault CloudFoundry auth plugin.

A simple examples of usage would be:

```
$ cf apps
Getting apps in org my-org / space my-space as admin...
OK

name               requested state   instances   memory   disk   urls
cf-vault-example   started           1/1         1G       1G     cf-vault-example-comedic-emu.example.com

$ cf vault apps set -a cf-vault-example -p my-policy

$ cf vault apps get -a cf-vault-example
Retrieving mapping information from Vault...

App: 		cf-vault-example (091233ae-c5f3-43c3-8b2f-940c8931059d)

Mapped policies: 		admin
Inherited space policies:
Inherited org policies:
```

This would map the application `cf-vault-example` to the existing Vault policy `my-policy`. This 
plugin does NOT take responsibility for provisioning Vault policies themselves, just mapping.

## Installation

Build and install the plugin like so:

```
go build .

chmod +x cf-vault-cli

cf install-plugin "$(pwd)/cf-vault-cli"

cf vault
```

## Usage

Usage overview of the plugin:

```
NAME:
   vault - Manage how CloudFoundry orgs, spaces and applications map to Vault policies

USAGE:
   vault [global options] command [command options] [arguments...]

VERSION:
   0.0.1

DESCRIPTION:
   Manage how CloudFoundry orgs, spaces and applications map to Vault policies

COMMANDS:
     apps     operations related to apps
     spaces   operations related to spaces
     orgs     operations related to orgs
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --mount value, -m value  name of the Vault mount where the CF auth plugin is installed (default: "cf")
   --help, -h               show help
   --version, -v            print the version
```