#!/bin/bash

gpath=$GOPATH
if [ -z "$gpath" ]; then
	gpath=$HOME/go
fi

goodpath=$gpath/src/github.com/transcom/mymove

if [ ! "$PWD" -ef "$goodpath" ]; then
	echo "In order to build the server, the project must be checked out into your gopath"
	echo "read more at https://golang.org/doc/code.html#Workspaces"
	echo "Expected: $goodpath"
	echo "Actual: $PWD"
	exit 1
fi
