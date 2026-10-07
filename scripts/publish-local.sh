#!/usr/bin/env bash
# Scans prices from this machine and publishes them to GitHub Pages.
#
# Some shops block GitHub's servers, so a scan run here reads more of them.
# The scan is added to the published database (the `data` branch); pushing it
# makes the Pages workflow republish the site without scanning again.
#
# Requires: go and git push access to the repository.
set -euo pipefail

cd "$(dirname "$0")/.."
remote=$(git remote get-url origin)
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

echo "Fetching the published database…"
if git fetch -q origin data && git show FETCH_HEAD:pmscanner.db > "$work/pmscanner.db"; then
	echo "Adding to the published history."
else
	rm -f "$work/pmscanner.db"
	echo "No published database yet; starting one."
fi

echo "Building…"
go build -o "$work/pmscanner" .

echo "Scanning…"
# Exits non-zero only when every product failed; publish what we have anyway.
"$work/pmscanner" -db "$work/pmscanner.db" > "$work/scan.txt" || true
echo "Failures:"
grep 'error:' "$work/scan.txt" || echo "none"

echo "Publishing the database to the data branch…"
(
	cd "$work"
	git init -q -b data
	git add pmscanner.db
	git commit -q -m "Price database"
	git push -q -f "$remote" data
)

echo "Done. GitHub republishes the site in a minute or two."
