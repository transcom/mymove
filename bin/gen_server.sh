#!/bin/bash

mkdir -p pkg/gen
swagger generate server -f swagger.yaml -t pkg/gen/ --server-package genserver --model-package messages --exclude-main -A mymove
