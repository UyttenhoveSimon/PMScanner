package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Search engines should index the latest scan only, with its description.
func TestPageHeadForCrawlers(t *testing.T) {
	render := func(d pageData) string {
		var b strings.Builder
		if err := renderPage(&b, d); err != nil {
			t.Fatal(err)
		}
		return b.String()
	}
	latest := render(pageData{IsLatest: true, Canonical: siteURL, Shops: 3, Countries: []string{"BE", "FR"}, Scanned: "now"})
	for _, want := range []string{`<link rel="canonical" href="` + siteURL + `">`, "from 3 bullion dealers in 2 countries", `"@type": "Dataset"`, `class="intro"`} {
		if !strings.Contains(latest, want) {
			t.Errorf("latest page lacks %q", want)
		}
	}
	if strings.Contains(latest, "noindex") {
		t.Error("latest page is noindex")
	}
	if old := render(pageData{Scanned: "then"}); !strings.Contains(old, `content="noindex, follow"`) {
		t.Error("past scan is indexable")
	}
}

func TestWriteCrawlerFiles(t *testing.T) {
	dir := t.TempDir()
	if err := writeCrawlerFiles(dir, time.Date(2026, 11, 27, 14, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	sitemap, _ := os.ReadFile(filepath.Join(dir, "sitemap.xml"))
	if !strings.Contains(string(sitemap), "<lastmod>2026-11-27</lastmod>") {
		t.Errorf("sitemap = %s", sitemap)
	}
	llms, _ := os.ReadFile(filepath.Join(dir, "llms.txt"))
	if !strings.Contains(string(llms), siteURL+"api/prices.json") {
		t.Errorf("llms.txt = %s", llms)
	}
}
