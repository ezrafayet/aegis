# This was vibe coded, must change it but it works on my machine
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y \
        ca-certificates \
        wget && \
    rm -rf /var/lib/apt/lists/*

RUN groupadd -r authuser && useradd -r -g authuser authuser

WORKDIR /app

COPY --chown=authuser:authuser auth/main /app/main
COPY --chown=authuser:authuser auth/config.json /app/config.json

RUN chmod +x /app/main

USER authuser

EXPOSE 5666

CMD ["/app/main"]