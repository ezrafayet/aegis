# Aegis Development Pages

This directory contains the login success and error pages for the Aegis authentication system, served via nginx in a Docker container.

## Files

- `login-success.html` - Page shown after successful login
- `login-error.html` - Page shown after failed login
- `nginx.conf` - Nginx configuration
- `Dockerfile` - Docker image definition
- `docker-compose.yml` - Docker Compose configuration

## Quick Start

### Using Docker Compose (Recommended)

1. Navigate to the dev directory:
   ```bash
   cd dev
   ```

2. Build and start the container:
   ```bash
   docker-compose up -d
   ```

3. Access the pages:
   - Success page: http://localhost:5000/login-success
   - Error page: http://localhost:5000/login-error
   - Health check: http://localhost:5000/health

4. Stop the container:
   ```bash
   docker-compose down
   ```

### Using Docker directly

1. Build the image:
   ```bash
   docker build -t aegis-dev-pages .
   ```

2. Run the container:
   ```bash
   docker run -d -p 5000:80 --name aegis-dev-pages aegis-dev-pages
   ```

3. Stop and remove the container:
   ```bash
   docker stop aegis-dev-pages
   docker rm aegis-dev-pages
   ```

## Development

### Live Reload (Optional)

To enable live reload during development, uncomment the volume mounts in `docker-compose.yml`:

```yaml
volumes:
  - ./login-success.html:/usr/share/nginx/html/login-success.html:ro
  - ./login-error.html:/usr/share/nginx/html/login-error.html:ro
```

Then restart the container:
```bash
docker-compose down
docker-compose up -d
```

### Customizing Pages

The HTML pages include:
- Modern, responsive design
- Aegis branding
- Auto-redirect functionality (success page)
- Error parameter handling (error page)
- Smooth animations and transitions

### Configuration

The nginx configuration includes:
- Gzip compression
- Security headers
- Health check endpoint
- Error page handling
- Proper MIME types

## Integration with Aegis

These pages are designed to work with the Aegis authentication system. The URLs match the configuration in `config.json`:

- `redirect_after_success`: "http://localhost:5000/login-success"
- `redirect_after_error`: "http://localhost:5000/login-error"

## Troubleshooting

### Check container status:
```bash
docker-compose ps
```

### View logs:
```bash
docker-compose logs aegis-dev-pages
```

### Health check:
```bash
curl http://localhost:5000/health
```

### Rebuild after changes:
```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
``` 