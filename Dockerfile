FROM alpine:latest

# Install dependencies
RUN apk --no-cache add ca-certificates curl git

# Create non-root user
RUN adduser -D -s /bin/sh cloudpork

# Copy binary
COPY cloudpork /usr/local/bin/cloudpork

# Set permissions
RUN chmod +x /usr/local/bin/cloudpork

# Switch to non-root user
USER cloudpork

# Set working directory
WORKDIR /workspace

# Default command
ENTRYPOINT ["cloudpork"]
CMD ["--help"]