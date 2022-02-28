#!/usr/bin/env bash

set -euo pipefail

if [ ! -v NIX_PROFILE ]; then
  echo "NIX_PROFILE not set, not installing globally"
  echo "Try running 'direnv allow'"
  exit 1
fi

# make sure this is set, as we unset it for most projects
export NIX_SSL_CERT_FILE=/nix/var/nix/profiles/default/etc/ssl/certs/ca-bundle.crt
# needed for go-1.17.1 on x86_64-darwin *sigh*
export NIXPKGS_ALLOW_UNSUPPORTED_SYSTEM=1

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
# install packages
nix-env -f "${DIR}" -i
# Store a hash of this file to the hash of the nix profile
# This way if the config changes, we can warn about it via direnv
# See the nix config in .envrc
config_hash=$(nix-hash "${DIR}")
store_hash=$(nix-store -q --hash "${NIX_PROFILE}")
echo "${config_hash}-${store_hash}" > "${DIR}/../.nix-hash"
