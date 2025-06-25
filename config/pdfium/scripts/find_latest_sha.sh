#!/usr/bin/env bash
# Looks up the pdfium repo on googlesource and prints the newest
# chromium/NNNN branch (highest NNNN)
set -euo pipefail


# Google repository. There is a GitHub mirror, but this is the source.
REPO_URL="https://pdfium.googlesource.com/pdfium.git"

latest_line=$(
  git ls-remote --heads "$REPO_URL"  |             # fetch <sha>\t<ref>
  grep -E 'refs/heads/chromium/[0-9]+$' |          # keep chromium/NNNN only
  sort -t/ -k4,4n |                                # numeric sort on NNNN
  tail -n1                                         # highest = newest
)

latest_sha=$(awk '{print $1}' <<<"$latest_line")
latest_branch=$(awk -F'refs/heads/' '{print $2}' <<<"$latest_line")

printf "Latest PDFium Chromium branch : %s\n"  "$latest_branch"
printf "Commit SHA (branch HEAD)      : %s\n"  "$latest_sha"