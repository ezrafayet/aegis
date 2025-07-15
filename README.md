```
 █████╗ ███████╗ ██████╗ ██╗███████╗
██╔══██╗██╔════╝██╔════╝ ██║██╔════╝
███████║█████╗  ██║  ███╗██║███████╗
██╔══██║██╔══╝  ██║   ██║██║╚════██║
██║  ██║███████╗╚██████╔╝██║███████║
╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝╚══════╝
Drop-in auth service - no SaaS, no lock-in
```

# Introduction

I found myself either rewriting an authorization service each and every time on every project, or constantly re-using the same platforms and tools (Auth0, Supabase, Firebase, Pocket Base), which comes with heavy vendor lock-in, way too many features, big ecosystems and a pretty significant cost. I don't want the Gorilla and the whole universe.

I want to have just this: a simple DROP-IN auth service that I can just use in a docker for any project, with a single config file... Pretty much as one would use Nginx.

```
auth
|--- Dockerfile
|--- config.json
```

Let's see if I can do that over night...  
Spoiler alert: I did code it over night.  
Another spoiler alert: I've been improving it since then

Supported auth methods (to date):

- GitHub
- Discord

Also it won't support passwords ever since it should be considered a bad practise.

# How to use the auth service

See the architecture you need, the configuration and tutorials to set up auth providers (GitHub, Discord...).

## Architecture

You have 2 choices of architecture to use it:

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
     domain.com        +----------+  
           |           | CORE     |  
           +-----------> :8000    |  
                       +----------+  
```

## host on subdomains:

```
                       +----------+  
           +-----------> AEGIS    |  
           |           | :5666    |  
           |           +----------+  
    +------+---+      domain.com/auth
    | NGINX    |                     
    | :80      |                     
    +----------+                     
     domain.com        +----------+  
           |           | CORE     |  
           +-----------> :8000    |  
                       +----------+  
```

## Setup

In your project, just drop 2 files:

```
auth
|--- Dockerfile
|--- config.json
```

### Dockerfile

```Dockerfile
FROM ezrafayet/aegix:v0.6.0
COPY ./config.json /app/config.json
```

### config.json

```json
{
    "app": {
        "name": "MyApp",
        "url": "http://app.localhost:5000",
        "cors_allowed_origins": ["http://app.localhost:5000"],
        "early_adopters_only": false,
        "redirect_after_success": ["http://app.localhost:5000/login-success"],
        "redirect_after_error": "http://app.localhost:5000/login-error",
        "internal_api_keys": ["xxxxxxxxxxxx"],
        "endpoints_prefix": "/auth",
        "port": 5666
    },
    "rate-limiting": {
        "enabled": true
    },
    "statistics": {
        "enabled": true,
        "retention_months": 24
    },
    "admin_panel": {
        "enabled": true,
        "full_path": "/auth/admin"
    },
    "login_page": {
        "enabled": true,
        "full_path": "/auth/login"
    },
    "db": {
        "postgres_url": "xxxxxxxxxxxx"
    },
    "jwt": {
        "secret": "xxxxxxxxxxxx",
        "access_token_expiration_minutes": 1,
        "refresh_token_expiration_days": 30
    },
    "auth": {
        "providers": {
            "github": {
                "enabled": true,
                "app_name": "MyApp",
                "client_id": "xxxxxxxxxxxx",
                "client_secret": "xxxxxxxxxxxx"
            },
            "discord": {
                "enabled": true,
                "app_name": "MyApp",
                "client_id": "xxxxxxxxxxxx",
                "client_secret": "xxxxxxxxxxxx"
            }
        }
    },
    "cookie": {
        "domain": "app.localhost",
        "secure": false,
        "http_only": true,
        "same_site": 1,
        "path": "/"
    },
    "user": {
        "roles": ["platform_admin", "user"]
    }
}
```

## Tutorials

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

# Contributing

Contibuting is more than welcome. See CONTRIBUTING.md