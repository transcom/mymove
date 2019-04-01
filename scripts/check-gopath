#! /usr/bin/env bash

gpath=$GOPATH
if [ -z "$gpath" ]; then
  gpath=$HOME/go
fi

# Strip a trailing slash off of gpath if it exists (makes binpath below more robust)
gpath=${gpath%/}

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
