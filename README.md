# CloudFoundry Vault Auth Experiment

This repository contains code related to an experiment to allow CloudFoundry applications to obtain 
a token from Vault by using Diego instance identity certificates. This means that these CF applications
do not need to be explicitly provided with a token on startup via a service broker or some other separate
provisioning process.

![architecture](https://raw.githubusercontent.com/nthomson-pivotal/cf-vault-auth-experiment/master/docs/arch.png)

WARNING: Here be dragons

## Overview

This repository contains 3 sub-projects:

1. `cf-vault-auth-plugin` is a Vault auth plugin that allows CloudFoundry applications to authenticate
with Vault using CF instance identity certificates
1. `cf-vault-springboot-example` is an example Spring Boot application that obtains Vault tokens using
the above Vault auth plugin
1. `cf-vault-cli` is a CloudFoundry CLI plugin that provides convenience functions for mapping CF orgs,
spaces and applications to Vault roles in the Vault auth plugin

See the `README` files in each sub-project for more information.