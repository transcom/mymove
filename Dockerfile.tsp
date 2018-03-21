# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.10.0 AS build

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

# These linker flags create a standalone binary that will run in scratch.
RUN go build -o /bin/tsp-award-queue -ldflags "-linkmode external -extldflags -static" ./cmd/tsp_award_queue

# This results in a single layer image
# https://github.com/GoogleCloudPlatform/distroless
# This google maintained image is scratch plus some basic necessities like a tmp dir and root certs.
FROM gcr.io/distroless/base
COPY --from=build /bin/tsp-award-queue /bin/tsp-award-queue
COPY --from=build /go/src/github.com/transcom/mymove/config /config
COPY --from=build /go/bin/chamber /bin/chamber
COPY /build /build
ENTRYPOINT ["/bin/tsp-award-queue"]
