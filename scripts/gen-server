#! /usr/bin/env bash

set -eu -o pipefail

gendir=./pkg/gen

rm -rf $gendir
mkdir -p $gendir
./bin/swagger generate server -q -f swagger/internal.yaml -t $gendir --model-package internalmessages --server-package internalapi --api-package internaloperations --exclude-main -A mymove
./bin/swagger generate server -q -f swagger/api.yaml -t $gendir --model-package apimessages --server-package restapi --api-package apioperations --exclude-main -A mymove
./bin/swagger generate server -q -f swagger/orders.yaml -t $gendir --model-package ordersmessages --server-package ordersapi --api-package ordersoperations --exclude-main -A mymove
./bin/swagger generate server -q -f swagger/dps.yaml -t $gendir --model-package dpsmessages --server-package dpsapi --api-package dpsoperations --exclude-main -A mymove
