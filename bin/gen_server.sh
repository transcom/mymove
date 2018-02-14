#!/bin/bash

gendir=./pkg/gen

rm -rf $gendir
mkdir -p $gendir
./bin/swagger generate server -f swagger/swagger.yaml -t $gendir --model-package messages --exclude-main -A mymove
