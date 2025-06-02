# Use debian slim for glibc compatibility
FROM debian:bookworm-slim
# FROM alpine:3.18

# FROM alpine:3.18

# # Install ca-certificates for HTTPS requests
# RUN apk --no-cache add ca-certificates

# # Create a non-root user for security
# RUN addgroup -g 1001 authuser && adduser -u 1001 -G authuser 

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y \
        ca-certificates \
        wget && \
    rm -rf /var/lib/apt/lists/*

# Create a non-root user for security
RUN groupadd -r authuser && useradd -r -g authuser authuser

# Set working directory
WORKDIR /app

# Copy the compiled binary and configuration with explicit permissions
COPY --chown=authuser:authuser auth/main /app/main
COPY --chown=authuser:authuser auth/config.json /app/config.json

# Ensure the binary is executable
RUN chmod +x /app/main

# Switch to non-root user
USER authuser

# Expose the port
EXPOSE 5666

# Health check (optional - remove if /health endpoint doesn't exist)
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#     CMD wget --no-verbose --tries=1 --spider http://localhost:5666/health || exit 1

# Run the auth service using absolute path
CMD ["/app/main"] 