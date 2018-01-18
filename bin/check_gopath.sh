#!/bin/bash

gpath=$GOPATH
if [ -z "$gpath" ]; then
	gpath=$HOME/go
fi

goodpath=$gpath/src/github.com/transcom/mymove

if [ ! "$PWD" -ef "$goodpath" ]; then
	echo "In order to build the server, the project must be checked out into the correct path"
	echo "Expected: $PWD"
	echo "Actual: $goodpath"
	exit 1
fi
