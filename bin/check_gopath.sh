#! /usr/bin/env bash

gpath=$GOPATH
if [ -z "$gpath" ]; then
  gpath=$HOME/go
fi

# Strip a trailing slash off of gpath if it exists (makes binpath below more robust)
gpath=${gpath%/}

goodpath=$gpath/src/github.com/transcom/mymove

if [ ! "$PWD" -ef "$goodpath" ]; then
  echo "In order to build the server, the project must be checked out into your gopath"
  echo "read more at https://golang.org/doc/code.html#Workspaces"
  echo "Expected: $goodpath"
  echo "Actual: $PWD"
  exit 1
fi

# shellcheck disable=SC2154
if [[ "$gpath" == *"$home" ]]; then
  if [[ ! ":$PATH:" == *":$gpath/bin:"* ]] && [[ ! ":$PATH:" == *":~${gpath#$HOME}/bin:"* ]]; then
    echo "In order for go dependencies to be runnable, \$GOPATH/bin must be in your \$PATH"
    echo "Please add $gpath/bin to your \$PATH in your .bash_profile"
    echo "Expected: $gpath/bin"
    echo "Actual: $PATH"
  fi
else
  if [[ ! ":$PATH:" == *":$gpath/bin:"* ]]; then
    echo "In order for go dependencies to be runnable, \$GOPATH/bin must be in your \$PATH"
    echo "Please add $gpath/bin to your \$PATH in your .bash_profile"
    echo "Expected: $gpath/bin"
    echo "Actual: $PATH"
  fi
fi
