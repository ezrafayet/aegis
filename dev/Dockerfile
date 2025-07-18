# Use the official nginx image as base
FROM nginx:alpine

# Install curl for health checks
RUN apk add --no-cache curl

# Copy the HTML files to nginx's default serving directory
COPY html/login-success.html /usr/share/nginx/html/
COPY html/home.html /usr/share/nginx/html/
COPY html/404.html /usr/share/nginx/html/

# Copy the custom nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"] 