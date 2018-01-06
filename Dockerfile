# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy our server binary into the container
ADD server/bin/webserver-linux /app/server/webserver
ADD server/src/dp3/config /server_config

ADD client/build /app/client/

# Run the webserver by default when the container starts.
ENTRYPOINT /app/server/webserver \
  -entry /app/client/index.html \
  -build /app/client/ \
  -port :8080 \
  -config-dir /server_config \
  -env prod \
  -debug_logging

# Document that the service listens on port 8080.
EXPOSE 8080
