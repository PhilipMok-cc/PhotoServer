FROM golang:alpine AS builder

WORKDIR /app

# Copy the source code
COPY . .

# Build the application
RUN go build -o photo-server

# Use a smaller image for the final stage
FROM alpine:latest

# default values of UID and GID are 1000 and 100 respectively
ARG UID=1000
ARG GID=100

# Create a non-root user with the specified uid and gid
RUN adduser -D -s /bin/sh -u ${UID} -g ${GID} photouser

# Create the app directory
RUN mkdir /app 
RUN chown photouser:photouser /app 

# Create the photos directory inside the container (this is where the volume will be mounted)
RUN mkdir /photos
RUN chown photouser:photouser /photos
WORKDIR /app

USER photouser

# Copy the binary, static assets, and config
COPY --from=builder /app/photo-server . 
COPY --from=builder /app/static /app/static
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/config-docker.yaml /app/config.yaml

# Expose the port
EXPOSE 8080

# Define the command to run the application
ENTRYPOINT ["/app/photo-server"]