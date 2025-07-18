# Contributing

Thanks for openning this file!

All levels are welcome to engage with me and contribute to this project.

You can open an issue or engage on X ([https://x.com/nohokthen](https://x.com/nohokthen)).

## Pre-requisites

In order to develop new fetures you will need:

- Go installed
- Docker (to run the dev environment and integration tests), without the need of sudo to run it
- A postgres db running locally or in the cloud (try AWS, Neon, they have free tiers)
- An app in one of the providers in order to test auth

## Commands

I will not go over build, fmt, vet, test, since they are pretty obvious.

The commands for development are:

```
# To  start the dev environemnet
make start

# It will launch the mock project inside of /dev, that is using the local Aegis for auth
# You can access this website at localhost:5000

# Kill it
make kill
```

# For development

## Knowing the code

## Running dev mode

## How to add a provider

## Ideas that may reach master

Pick from them if you want to help
- First dot = 🟢 low priority / 🔴 high priority
- Second dot =  🟢 easy (or short) / 🔴 hard (or long)

[🔴🟢] - Add Google provider
[🔴🟢] - Add Facebook provider
[🔴🟢] - Support API tokens
[🔴🟢] - Add proper logs + verbose logs, add a request ID for trace if none
[🔴🟢] - Add opt-in rate limiting
[🔴🟢] - Write better doc and tutorials
[🔴🟢] - Add validation of the config file
[🔴🟢] - Add env vars support in the config file, probably like {{env:SOME_VAR}}
[🔴🔴] - Enhance Device Fingerprinting
[🟢🟢] - Add a built in 404 page
[🟢🟢] - Add Microsoft provider
[🟢🟢] - Add Apple provider
[🟢🟢] - Test the redirects to custom paths from config in integration tests
[🟢🟢] - Remove dependency of package cookies on entities
[🟢🟢] - Database encoding (add a secret to encode/decode the db)
[🟢🟢] - Deploy a new version on DockerHub on push of a new tag
[🟢🔴] - Write SDKs for other languages
[🟢🔴] - Add an opt-in admin page
[🟢🔴] - Add opt-in statistics
[🟢🔴] - Create fake providers for tests
[🟢🔴] - Use fake providers for tests, and write e2e tests
[🟢🔴] - Improve integration tests by creating containers once
[🟢🔴] - Investigate for support of magic links / magic texts
[🟢🔴] - Investigate support for SAML 2, OpenID
[🟢🔴] - Investigate support for WebAuthn support
[🟢🔴] - Investigate support of a fine resource access system




## Publish a new version on Docker Hub

```
sudo docker build -t ezrafayet/aegis:v0.3.0 .
sudo docker push ezrafayet/aegis:v0.3.0
```
