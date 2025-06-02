```
 █████╗ ██╗   ██╗████████╗██╗  ██╗     █████╗ ███████╗ ██████╗ ██╗██╗  ██╗
██╔══██╗██║   ██║╚══██╔══╝██║  ██║    ██╔══██╗██╔════╝██╔════╝ ██║╚██╗██╔╝
███████║██║   ██║   ██║   ███████║    ███████║█████╗  ██║  ███╗██║ ╚███╔╝ 
██╔══██║██║   ██║   ██║   ██╔══██║    ██╔══██║██╔══╝  ██║   ██║██║ ██╔██╗ 
██║  ██║╚██████╔╝   ██║   ██║  ██║    ██║  ██║███████╗╚██████╔╝██║██╔╝ ██╗
╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝    ╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝╚═╝  ╚═╝
Drop-in auth service - no SaaS, no lock-in
```

I found myself rewriting an authorization service each and every time on every project or reusing the same platforms and tools (Auth0, Supabase, Firebase, Pocket Base), which comes with heavy vendor lock-in, way too many features (I don't want the Gorilla and the whole jungle), big ecosystems and a pretty significant cost.

I want to have just this: a simple DROP-IN auth service that I can just use in a docker for any project, with a single config file... Pretty much as one would use Nginx.

```
auth
|--- Dockerfile
|--- config.json
```

And that's it! Let's see if I can do that over night...

Spoiler alert: I did! Having a working version for all basic flows took me 2 nights

# Project status
- Implemented GitHub OAuth to get started
- Finished all the basic flows and life cycyles of the access token and refresh token
- Must add tests
- Must host on Docker (that's the whole point of the project)
- Code is not so good - repetitions etc, must be arranged - but the overall architecture is ready

# Should be implemented
- Roles for RBA
- Server 2 server checks
- creation of api tokens
- adding metadata
- blocking and deleting users (already handled)
- More providers!

Also it won't support passwords since it is bad practise.

# Doc to write
- How to get started under the same domain <= tested and works
- How to get started on different subdomains
- How to get started on 2 different domains completly
- How to use it from client (Next.js etc) or server (Node.js, Python, Go...)

# Example of a working config
```
{
    "app": {
        "name": "App Name",
        "url": "app.localhost:5000",
        "env": "development",
        "log_level": "debug",
        "api_keys": ["xxx"]
    },
    "db": {
        "postgres_url": "xxx"
    },
    "jwt": {
        "secret": "xxx",
        "access_token_expiration_minutes": 15,
        "refresh_token_expiration_days": 30
    },
    "auth": {
        "providers": {
            "github": {
                "enabled": true,
                "client_id": "xxx",
                "client_secret": "xxx"
            }
        },
        "allowed_origins": [
            "http://app.localhost:5000"
        ]
    },
    "cookie": {
        "domain": "app.localhost:5000",
        "secure": false,
        "http_only": true,
        "same_site": "Lax",
        "path": "/"
    },
    "user": {
        "roles": ["platform_admin", "user"],
        "metadata": {
            "foo": {
                "type": "string",
                "default": "bar",
                "enum": ["bar", "baz"]
            }
        }
    }
}
```
