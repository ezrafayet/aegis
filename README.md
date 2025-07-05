```
 █████╗ ███████╗ ██████╗ ██╗███████╗     ██████╗ ████████╗██╗  ██╗
██╔══██╗██╔════╝██╔════╝ ██║██╔════╝    ██╔═████╗╚══██╔══╝██║  ██║
███████║█████╗  ██║  ███╗██║███████╗    ██║██╔██║   ██║   ███████║
██╔══██║██╔══╝  ██║   ██║██║╚════██║    ████╔╝██║   ██║   ██╔══██║
██║  ██║███████╗╚██████╔╝██║███████║    ╚██████╔╝   ██║   ██║  ██║
╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝╚══════╝     ╚═════╝    ╚═╝   ╚═╝  ╚═╝
Drop-in auth service - no SaaS, no lock-in
```

# OAuth support:

- GitHub
- Discord

# Introduction

I found myself rewriting an authorization service each and every time on every project or constantly using the same platforms and tools (Auth0, Supabase, Firebase, Pocket Base), which comes with heavy vendor lock-in, way too many features (I don't want the Gorilla and the whole jungle), big ecosystems and a pretty significant cost.

I want to have just this: a simple DROP-IN auth service that I can just use in a docker for any project, with a single config file... Pretty much as one would use Nginx.

```
auth
|--- Dockerfile
|--- config.json
```

And that's it! Let's see if I can do that over night...
Spoiler alert: I did code it over night!
Another spoiler alert: I've been improving it since then

# How to use the auth service

In your project, just drop 2 files:

## 1. A Dockerfile

```Dockerfile
FROM ezrafayet/aegix:v0.2.0
COPY ./config.json /app/config.json
```

## 2. A config.json

```json
{
    "app": {
        "name": "MyApp",
        "url": "app.localhost:5000",
        "allowed_origins": ["http://app.localhost:5000"],
        "api_keys": ["xxxxxxxxxxxx"],
        "port": 5666
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
                "client_secret": "xxxxxxxxxxxx",
                "redirect_url": "http://app.localhost:5000/oauth/github/callback"
            },
            "discord": {
                "enabled": true,
                "app_name": "MyApp",
                "client_id": "xxxxxxxxxxxx",
                "client_secret": "xxxxxxxxxxxx",
                "redirect_url": "http://app.localhost:5000/oauth/github/callback"
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

# Project status
- Implemented GitHub OAuth to get started
- Finished all the basic flows and life cycyles of the access token and refresh token
- Must add tests (especially integration)

# Should be implemented
- Server 2 server checks
- creation of api tokens
- adding metadata
- More providers!
- A dashboard

Also it won't support passwords since it is bad practise.

# Doc to write
- How to get started under the same domain <= tested and works
- How to get started on different subdomains
- How to use it from client (Next.js etc) or server (Node.js, Python, Go...)

# Security

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
