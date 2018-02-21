#!/bin/bash

gendir=./pkg/gen

rm -rf $gendir
mkdir -p $gendir
./bin/swagger generate server -f swagger/internal.yaml -t $gendir --model-package internalmessages --server-package internalapi --api-package internaloperations --exclude-main -A mymove
./bin/swagger generate server -f swagger/api.yaml -t $gendir --model-package apimessages --server-package restapi --api-package apioperations --exclude-main -A mymove
