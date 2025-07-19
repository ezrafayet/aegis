# Contributing

All contributions are welcome, from all forms and all levels.

You can open a GitHub issue [here](https://github.com/ezrafayet/aegis/issues) or [engage on X](https://x.com/nohokthen) to align on what you want to do.

Feel free to pick something from the list below or to come up with your own.

# Brainstorming section

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
[🟢🟢] - Improve dev mode  
[🟢🔴] - Add support for external plugins  
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

# Pre-requisites

In order to run the project locally you need:
- Go (quite obviously)
- Docker (to run the dev environment and integration tests)
- A Postgres db running locally or in the cloud
- An app in one of the providers in order to test auth (GitHub is easy)

## Commands

Check the comprehensive list of commands in Makefile.

Run them with:

```
make command
```

# The development environment

I have added a development environment to make it easy to test authentication.

Is it perfect or finished ? No, but it is pretty cool.

Run it with:

```
make start
```

It runs Aegis as a service and the service under /dev that uses it. Note you will need a config.json in src for it to run.

A classic config.json:

```json
{
  "app": {
    "name": "Aegis",
    "url": "http://localhost:5000",
    "cors_allowed_origins": ["http://localhost:5000"],
    "early_adopters_only": true,
    "redirect_after_success": "http://localhost:5000/login-success",
    "redirect_after_error": "http://localhost:5000/auth/login-error",
    "internal_api_keys": ["dev-api-key-123"],
    "port": 5666
  },
  "login_page": {
    "enabled": true,
    "full_path": "/auth/login"
  },
  "error_page": {
    "enabled": true,
    "full_path": "/auth/login-error"
  },
  "404_page": {
    "enabled": true
  },
  "db": {
    "postgres_url": "your_pg_here"
  },
  "jwt": {
    "secret": "your-super-secret-jwt-key-change-this-in-production",
    "access_token_expiration_minutes": 15,
    "refresh_token_expiration_days": 7
  },
  "auth": {
    "providers": {
      "github": {
        "enabled": true,
        "app_name": "Scribylon",
        "client_id": "Ov23li1IQUp8Qu5k0PtA",
        "client_secret": "cfd60be9d519e219a2473217bffd12c1ba23ef0d",
        "redirect_url": "http://localhost:5000/auth/github/callback"
      },
      "discord": {
          "enabled": true,
          "app_name": "Scribylon",
          "client_id": "1289110322963152957",
          "client_secret": "sAJ113ESkC1aFPiweCUpweomq6R825GI",
          "redirect_url": "http://localhost:5000/auth/discord/callback"
      }
    }
  },
  "cookies": {
    "domain": "localhost",
    "secure": false,
    "http_only": true,
    "same_site": 1,
    "path": "/"
  },
  "user": {
    "roles": ["user", "platform_admin"]
  }
} 
```

# Advices

- Check the behavior of the API in the [integration tests](https://github.com/ezrafayet/aegis/blob/master/src/integration/integration_test.go)
- Run the dev mode
- Write something simple first to get to know how it works (example: add a provider)

# How to publish a new version on Docker Hub

```
sudo docker build -t ezrafayet/aegis:v0.3.0 .
sudo docker push ezrafayet/aegis:v0.3.0
```
