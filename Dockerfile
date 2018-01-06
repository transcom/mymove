# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang AS build

# Install tools required to build the project
# We will need to run `docker build --no-cache .` to update those dependencies
RUN apt-get install git
RUN go get github.com/Masterminds/glide

# glide.yaml and glide.lock lists project dependencies
# These layers will only be re-built when Gopkg files are updated
COPY server/src/dp3/glide.lock server/src/dp3/glide.yaml /go/src/dp3/
WORKDIR /go/src/dp3/
# Install library dependencies
RUN glide install

# Copy all project and build it
# This layer will be rebuilt when ever a file has changed in the project directory
COPY server/src/dp3 /go/src/dp3/
WORKDIR /go/src/dp3/cmd/webserver/
# These linker flags create a standalone binary that will run in scratch.
RUN go build -o /bin/dp3-server -ldflags "-linkmode external -extldflags -static"

# This results in a single layer image
FROM scratch
COPY --from=build /bin/dp3-server /bin/dp3-server
COPY --from=build /go/src/dp3/config /server_config
COPY /client/build /app/client
ENTRYPOINT ["/bin/dp3-server"]
CMD ["-entry", "/app/client/index.html", \
     "-build", "/app/client/", \
     "-port", ":8080", \
     "-config-dir", "/server_config", \
     "-env", "prod", \
     "-debug_logging"]

EXPOSE 8080
