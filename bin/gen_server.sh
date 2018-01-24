#!/bin/bash

swagger generate server -f swagger.yaml -t pkg/ --server-package genserver --model-package messages --exclude-main -A mymove
