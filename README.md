```
 █████╗ ██╗   ██╗████████╗██╗  ██╗     █████╗ ███████╗ ██████╗ ██╗██╗  ██╗
██╔══██╗██║   ██║╚══██╔══╝██║  ██║    ██╔══██╗██╔════╝██╔════╝ ██║╚██╗██╔╝
███████║██║   ██║   ██║   ███████║    ███████║█████╗  ██║  ███╗██║ ╚███╔╝ 
██╔══██║██║   ██║   ██║   ██╔══██║    ██╔══██║██╔══╝  ██║   ██║██║ ██╔██╗ 
██║  ██║╚██████╔╝   ██║   ██║  ██║    ██║  ██║███████╗╚██████╔╝██║██╔╝ ██╗
╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝  ╚═╝    ╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝╚═╝  ╚═╝
Drop-in auth service - no SaaS, no lock-in
```

On every side project I'm just rewriting auth every time, or depending on Auth0, Supabase, Firebase, with heavy vendor lock-in and echosystems to learn.

I want to have just this: a simple drop-in auth service that I can just use in a docker for any project, with a single config file... Pretty much as one would add Nginx.

Also it won't support passwords since it is bad practise.

For now it will only support GitHub OAuth, and only a pg DB with basic config.

Let's see if I can do that over night...

It must handle:
- client login, request for session etc
- server requests to get a session from a jwt
- creation of api tokens
- role based auth
- adding metadata (just a stringified json is fine there)

In the future I should provide snippets for other projects to use (nextjs, go, node)

The user should be able to host it on a subdomain (auth.domain.com) or behind some gateway at domain.com/auth

Example of a working config:
```
put everything behind the same domain

"cookie": {
        "domain": "app.localhost:5000",
        "secure": false,
        "http_only": true,
        "same_site": "Lax",
        "path": "/"
    },
```

give more examples with different subdomain

give more examples with cross domains
