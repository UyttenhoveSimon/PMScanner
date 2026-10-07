package scraper

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/shopspring/decimal"
)

// BrowserEnv names the environment variable that overrides the browser path.
const BrowserEnv = "PMSCANNER_BROWSER"

// browserCandidates are Chromium-based browsers, tried in order.
var browserCandidates = []string{
	"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
	"/Applications/Chromium.app/Contents/MacOS/Chromium",
	"/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
	"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
	"google-chrome", "google-chrome-stable", "chromium", "chromium-browser",
	"brave-browser", "microsoft-edge",
}

func findBrowser() (string, error) {
	if path := os.Getenv(BrowserEnv); path != "" {
		return path, nil
	}
	for _, c := range browserCandidates {
		if path, err := exec.LookPath(c); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no Chromium-based browser found; set %s", BrowserEnv)
}

const (
	browserTabs    = 8
	pageTimeout    = 30 * time.Second
	pricePollDelay = 2 * time.Second
	// settleDelay leaves bot challenges alone before the first DOM read;
	// Cloudflare tends to fail its check when DevTools touches the page early.
	settleDelay = 6 * time.Second
)

// scrapeBrowser fills results[i] for each index, loading pages in parallel
// tabs of one headless browser.
func scrapeBrowser(results []Result, indexes []int) {
	fail := func(err error) {
		for _, i := range indexes {
			results[i].Err = err
		}
	}
	path, err := findBrowser()
	if err != nil {
		fail(err)
		return
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(path),
		chromedp.UserAgent(userAgent),
		// Look less like an automated browser to bot challenges.
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
	)
	actx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	bctx, cancel := chromedp.NewContext(actx)
	defer cancel()
	if err := chromedp.Do(bctx); err != nil { // starts the browser
		fail(fmt.Errorf("start browser: %w", err))
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, browserTabs)
	for _, i := range indexes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			tctx, cancel := chromedp.NewContext(bctx)
			defer cancel()
			r := &results[i]
			r.Price, r.Currency, r.Err = scrapeTab(tctx, r.Product)
		}()
	}
	wg.Wait()
}

// scrapeTab loads the page and re-reads it until a price appears, since
// challenge pages and client-side rendering both delay the real content.
func scrapeTab(ctx context.Context, p Product) (decimal.Decimal, string, error) {
	if err := chromedp.Do(ctx); err != nil { // opens the tab
		return decimal.Zero, "", fmt.Errorf("open tab: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, pageTimeout)
	defer cancel()

	steps := []chromedp.Action[chromedp.Void]{}
	for name, value := range p.Cookies {
		steps = append(steps, setCookie(name, value, p.URL))
	}
	steps = append(steps, chromedp.Navigate(p.URL), chromedp.Sleep(settleDelay))
	if err := chromedp.Do(ctx, steps...); err != nil {
		return decimal.Zero, "", fmt.Errorf("load: %w", err)
	}
	for {
		html, err := chromedp.Run(ctx, chromedp.OuterHTML("html"))
		if err == nil {
			price, currency, err := extract([]byte(html), p)
			if !errors.Is(err, errNoPrice) {
				return price, currency, err
			}
		}
		select {
		case <-ctx.Done():
			return decimal.Zero, "", errNoPrice
		case <-time.After(pricePollDelay):
		}
	}
}

func setCookie(name, value, url string) chromedp.Action[chromedp.Void] {
	return func(ctx context.Context, t *chromedp.Target) (chromedp.Void, error) {
		_, err := cdp.Call(ctx, t, network.SetCookie, network.SetCookieParams{
			Name: name, Value: value, URL: url,
		})
		return chromedp.Void{}, err
	}
}
