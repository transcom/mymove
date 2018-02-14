#!/bin/bash

set -eu

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DOCS_DIR="${DIR}/../docs"

yarn run markdown-toc -i "${DOCS_DIR}/frontend.md" --bullets='*'
#yarn run markdown-toc -i "${DOCS_DIR}/backend.md"
