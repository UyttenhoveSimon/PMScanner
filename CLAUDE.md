# PMScanner — notes for AI agents

- Read [README.md](README.md) (options, taxes, "For shop owners") and
  [docs/NOTES.md](docs/NOTES.md) (per-shop findings, VAT checks, open issues)
  before changing products or scraping logic. Update docs/NOTES.md with what
  you learn.
- Commit and push directly on `main`. CI pushes "Add scanned prices" commits
  daily, so `git fetch` and rebase on `origin/main` first.
- Check: `gofmt -l . && go vet ./... && go test ./...`
- Prices must be exact. Before switching a product to an API, or changing its
  `VAT`, confirm the number against the visible page price (incl. VAT).
- Be gentle with shops: keep at most 2 concurrent requests per shop/RateGroup
  (also in the browser); don't hammer bot-walled shops while testing.
- CI runs in the US: shops that localise by IP (e.g. suissegold) may show other
  currencies or prices without VAT there (`ExportCurrency`).
