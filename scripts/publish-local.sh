#!/usr/bin/env bash
# Scans prices from this machine and publishes them to GitHub Pages.
#
# Some shops block GitHub's servers, so a scan run here reads more of them.
# The scan is committed to data/scans on main; pushing it makes the Pages
# workflow republish the site without scanning again.
#
# Requires: go and push access to the repository.
set -euo pipefail

cd "$(dirname "$0")/.."
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

echo "Updating from GitHub…"
git pull -q --ff-only origin main

echo "Building…"
go build -o "$work/pmscanner" .

echo "Scanning…"
# Exits non-zero only when every product failed; publish what we have anyway.
"$work/pmscanner" -db "$work/pmscanner.db" -history data/scans > "$work/scan.txt" || true
echo "Failures:"
grep 'error:' "$work/scan.txt" || echo "none"

echo "Publishing…"
git add data/scans
git commit -q -m "Add prices scanned locally" -- data/scans
git push -q origin HEAD:main
echo "Done. GitHub republishes the site in a minute or two."
