# Configuration

## Container Definitions

The `*container-definition*` files define how the ECS containers are configured. They are primarily used to set non-secret environment variables.

## TLS cert/key (optional)

The `devlocal-https.*` files are a self-signed TLS cert/key pair. They are a [snake oil](https://en.wikipedia.org/wiki/Snake_oil_(cryptography)) certificate used to locally run the webserver during development. They are included as a convenience so engineers don't have to generate their own.
