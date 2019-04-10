# CloudFoundry Vault Auth Experiment

This repository contains code related to an experiment to integrate CloudFoundry applications to Vault
for authentication purposes.

WARNING: The code is pretty messy and hacky.

## Overview

This repository contains 3 sub-projects:

1. `cf-vault-auth-plugin` is a Vault auth plugin that allows CloudFoundry applications to authenticate
with Vault using CF instance identity certificates
1. `cf-vault-springboot-example` is an example Spring Boot application that obtains Vault tokens using
the above Vault auth plugin
1. `cf-vault-cli` is a CloudFoundry CLI plugin that provides convenience functions for mapping CF orgs,
spaces and applications to Vault roles in the Vault auth plugin

See the README files in each sub-project for more information.