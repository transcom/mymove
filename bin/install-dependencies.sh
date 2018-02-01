#! /bin/bash

set -eux -o pipefail

# Install utilities

# This is needed to use `psql` to test DB connectivity, until the app itself
# starts making database connections.
sudo apt-get -qq update
sudo apt-get -qq install -y postgresql-client-9.4 > /dev/null

# Install dep
go get -u github.com/golang/dep/cmd/dep
