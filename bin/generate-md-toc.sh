#!/bin/bash
# Wrapper script to generate table of contents on Markdown files.
# Markdown files must have the `<!-- toc -->` tag in them in order to receive
# a TOC.
#
# This is intended to be called by a pre-commit hook to reduce overhead.
#
# Generated with https://github.com/jonschlinkert/markdown-toc

set -eu -o pipefail

function generate_toc() {
  filename="$1"
  regen=$'\n\n_Regenerate with `bin/generate-md-toc.sh`_'

  # shellcheck disable=SC2016
  yarn run markdown-toc -i "${filename}" --bullets='*' --append="${regen}"
}

filename="${1:-}"

if [ -z "$filename" ]; then
  echo "error: must supply a markdown file to generate table of contents."
  exit 1
fi

generate_toc "$filename"
