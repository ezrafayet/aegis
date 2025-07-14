# Use the official nginx image as base
FROM nginx:alpine

# Install curl for health checks
RUN apk add --no-cache curl

# Copy the HTML files to nginx's default serving directory
COPY login-success.html /usr/share/nginx/html/
COPY login-error.html /usr/share/nginx/html/
COPY home.html /usr/share/nginx/html/
COPY 404.html /usr/share/nginx/html/

# Copy the custom nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Create a simple index page that redirects to login-success
RUN echo '<!DOCTYPE html><html><head><meta http-equiv="refresh" content="0;url=/login-success"></head><body>Redirecting...</body></html>' > /usr/share/nginx/html/index.html

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"] 