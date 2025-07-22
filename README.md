```
 █████╗ ███████╗ ██████╗ ██╗███████╗
██╔══██╗██╔════╝██╔════╝ ██║██╔════╝
███████║█████╗  ██║  ███╗██║███████╗
██╔══██║██╔══╝  ██║   ██║██║╚════██║
██║  ██║███████╗╚██████╔╝██║███████║
╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝╚══════╝
Drop-in auth service - no SaaS, no lock-in, just the banana.
```

# Introduction

I found myself rewriting authorization services on each and every side-project, or always re-using the same 3rd parties that come with heavy vendor lock-in, way too many features, big ecosystems and a pretty significant cost (Auth0, Supabase, Firebase, Pocket Base). This is very time consuming, and I don't want the Gorilla and the whole universe. Just the bannana.

I want just this: a simple DROP-IN auth service that I can just use  with a single config file... Pretty much like one would use Nginx:

```
auth
|--- Dockerfile
|--- config.json
```

This was a one night POC, became a project. Contibuting is more than welcome. See [CONTRIBUTING.md](./CONTRIBUTING.md)

Also it won't support passwords ever since they are a bad practise.

# The configuration files

In your project, create an 'auth' folder for the service and just create 2 files:

## Dockerfile

```Dockerfile
FROM ezrafayet/aegis:v0.11.1
COPY ./config.json /app/config.json
```

## And the config.json

```json
{
    "app": {
        "name": "MyApp",
        "url": "http://localhost:5000",
        "cors_allowed_origins": ["http://localhost:5000"],
        "early_adopters_only": false,
        "redirect_after_success": ["http://localhost:5000/login-success"],
        "redirect_after_error": "http://localhost:5000/auth/login-error",
        "internal_api_keys": ["${env:AEGIS_INTERNAL_API_KEY}"],
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
    "db": {
        "postgres_url": "${env:AEGIS_DPOSTGRES_URL}"
    },
    "jwt": {
        "secret": "${env:AEGIS_JWT_SECRET}",
        "access_token_expiration_minutes": 1,
        "refresh_token_expiration_days": 30
    },
    "auth": {
        "providers": {
            "github": {
                "enabled": true,
                "app_name": "MyApp",
                "client_id": "${env:AEGIS_GITHUB_CLIENT_ID}",
                "client_secret": "${env:AEGIS_GITHUB_CLIENT_SECRET}",
                "redirect_url": "http://localhost:5000/auth/github/callback"
            },
            "discord": {
                "enabled": true,
                "app_name": "MyApp",
                "client_id": "${env:AEGIS_DISCORD_CLIENT_ID}",
                "client_secret": "${env:AEGIS_DISCORD_CLIENT_SECRET}",
                "redirect_url": "http://localhost:5000/auth/discord/callback"
            }
        }
    },
    "cookies": {
        "domain": "localhost",
        "secure": false,
        "http_only": true,
        "same_site": 3,
        "path": "/"
    }
}
```

And now you have authentication & authorization.

# Architecture

You have multiple choices of architecture to use it:

## host on the same domain:

```
                       +----------+  
           +-----------> AEGIS    |  
           |           | :5666    |  
           |           +----------+  
    +------+---+      domain.com/auth
    | NGINX    |                     
    | :80      |                     
    +----------+                     
           |           +----------+  
           |           | CORE     |  
           +-----------> :8000    |  
                       +----------+  
                       domain.com
```

The cookies configuration for it:
```json
{
    "domain": "localhost",
    "secure": false,
    "http_only": true,
    "same_site": 3,
    "path": "/"
}
```

## host everything on the same subdomain:

```
                       +----------+  
           +-----------> AEGIS    |  
           |           | :5666    |  
           |           +----------+  
    +------+?          app.domain.com/auth
           |     
           |           +----------+  
           |           | CORE/    |  
           +-----------> :8000    |  
                       +----------+  
                       app.domain.com

?: can be nginx or even hosted independently
```

The cookies configuration for it:
```json
{
    "domain": "app.localhost",
    "secure": false,
    "http_only": true,
    "same_site": 3,
    "path": "/"
}
```

## host on different subdomains:

```
                       +----------+  
           +-----------> AEGIS    |  
           |           | :5666    |  
           |           +----------+  
    +------+?      auth.domain.com/auth
           |          
           |           +----------+  
           |           | CORE/    |  
           +-----------> :8000    |  
                       +----------+  
                       domain.com

?: can be nginx or even hosted independently
```

The cookies configuration for it:
```json
Not tested yet
```

# Providers

In order to have authentication (so the user can login with a provider), you need to have an app on those providers.

Currently, the supported providers are:
- Discord
- GitHub

Tutorials (to come):

- Setup GitHub auth (to come)
- Setup Discord auth (to come)

# Security

If you need a bullet-proof, battle-tested auth for production, do not use this service. Use Auth0, NextAuth, Supabase, Firebase, Work.os, but not this service.

## Implemented

### ✅ XSS (Cross-Site Scripting)

**Description**: Attackers inject malicious scripts into web pages to steal authentication tokens or user data.

**Prevention**: 
- **HTTP-only cookies**: Authentication tokens are stored in `HttpOnly` cookies that JavaScript cannot access

### ✅ CSRF (Cross-Site Request Forgery)

**Description**: Attackers trick authenticated users into performing unwanted actions on your application.

**Prevention**:
- **OAuth state parameter**: Random, unguessable state tokens prevent unauthorized OAuth callbacks
- **State expiration**: States expire after 3 minutes to limit attack window
- **One-time use**: States are deleted after verification to prevent replay attacks

### ✅ Session Fixation

**Description**: Attackers force users to use a known session ID, then hijack the session after authentication.

**Prevention**:
- **New tokens per login**: Fresh access and refresh tokens are generated on every OAuth login
- **Token rotation**: Existing refresh tokens are invalidated when new ones are issued
- **Device fingerprinting**: Tokens are tied to specific device fingerprints (needs improvement)

### ✅ Error Information Disclosure

**Description**: Detailed error messages can reveal system information to attackers.

**Prevention**:
- **Generic error messages**: Use consistent, non-leaking error responses
- **Logging separation**: Log detailed errors internally, return generic messages to users

## Needs Implementation

### ⚠️ Token Hijacking

**Description**: Attackers steal refresh tokens and use them to impersonate users from different devices.

**Current Risk**: 
- Refresh tokens can be used from any device
- No device validation during token refresh

**Prevention Needed**:
- **Device fingerprinting**: Generate unique device IDs from User-Agent, IP, and other headers
- **Device validation**: Verify device fingerprint matches during token refresh
- **Token binding**: Bind refresh tokens to specific device characteristics

### ⚠️ DDoS Protection

**Description**: Attackers overwhelm your service with requests to make it unavailable and inflate DB records.

**Current Protection**:
- **Basic rate limiting**: 20 requests per minute globally (too permissive)

**Prevention Needed**:
- **Per-endpoint rate limiting**: Different limits for different endpoints
- **IP-based limiting**: Track requests per IP address
- **Distributed rate limiting**: Use external storage or similar for persistent rate limiting

### ⚠️ JWT Secret Strength

**Description**: Weak JWT secrets can be brute-forced to forge tokens.

**Prevention Needed**:
- **Secret validation**: Ensure JWT secrets are at least 32 characters
- **Entropy checking**: Validate secret randomness
