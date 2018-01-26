# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang AS build

# Install tools required to build the project
# We will need to run `docker build --no-cache .` to update those dependencies
RUN apt-get install git
RUN go get github.com/golang/dep/cmd/dep

# Copy all project and build it
# This layer will be rebuilt when ever a file has changed in the project directory
COPY ./ /go/src/github.com/transcom/mymove/
WORKDIR /go/src/github.com/transcom/mymove/
RUN rm .*.stamp
RUN make server_deps
RUN make server_generate
# These linker flags create a standalone binary that will run in scratch.
RUN go build -o /bin/mymove-server -ldflags "-linkmode external -extldflags -static" ./cmd/webserver

# This results in a single layer image
FROM gcr.io/distroless/base
COPY --from=build /bin/mymove-server /bin/mymove-server
COPY --from=build /go/src/github.com/transcom/mymove/config /server_config
COPY --from=build /go/src/github.com/transcom/mymove/swagger.yaml /app/swagger.yaml
COPY /build /app/client
ENTRYPOINT ["/bin/mymove-server"]
CMD ["-build", "/app/client/", \
     "-swagger", "/app/swagger.yaml", \
     "-port", ":8080", \
     "-config-dir", "/server_config", \
     "-env", "prod", \
     "-debug_logging"]

EXPOSE 8080
