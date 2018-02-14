#!/bin/bash

set -eu -o pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DOCS_DIR="${DIR}/../docs"

function generate_toc() {
  filename="$1"
  regen=$'\n\n_Regenerate with `bin/docs-toc.sh`_'

  # shellcheck disable=SC2016
  yarn run markdown-toc -i "${filename}" --bullets='*' --append="${regen}"
}

generate_toc "${DOCS_DIR}/frontend.md"
generate_toc "${DOCS_DIR}/backend.md"
