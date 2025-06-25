#!/usr/bin/env bash
# This will verify the anchored sha.txt to the latest sha from pdfium.
#   Expected layout:
#     config/pdfium/
#       - sha.txt
#       - scripts/
#           - find_latest_sha.sh
#           - check_sha.sh
set -euo pipefail

# Resolve paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
SHA_FILE="$ROOT_DIR/sha.txt"
FIND_SCRIPT="$SCRIPT_DIR/find_latest_sha.sh"

# Get the latest SHA from the helper script
latest_info="$("$FIND_SCRIPT")"
latest_sha="$(awk '/Commit SHA/{print $NF}' <<<"$latest_info")"

# Read the stored SHA
stored_sha="$(tr -d '[:space:]%' < "$SHA_FILE")"

# Compare
if [[ "$latest_sha" == "$stored_sha" ]]; then
  echo "SUCCESS: sha.txt is up-to-date -> $stored_sha"
  exit 0
else
  echo "FAIL: sha.txt ($stored_sha) is outdated; latest is $latest_sha"
  # TODO: automatically update sha.txt and commit / open PR
  exit 1
fi
