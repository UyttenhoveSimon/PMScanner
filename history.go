package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"pmscanner/store"
)

// History keeps every scan as its own compressed JSON file in a directory,
// so the scans can live in the repository: files are only ever added, which
// keeps the git history small, unlike committing a growing database.

const historyTimeFormat = "2006-01-02T15-04-05Z"

type historyFile struct {
	At     time.Time     `json:"at"`
	Prices []store.Price `json:"prices"`
}

// importHistory adds to the database every scan file it does not have yet.
func importHistory(dir string, db *store.Store) error {
	names, err := filepath.Glob(filepath.Join(dir, "*.json.gz"))
	if err != nil {
		return err
	}
	have, err := scanTimes(db)
	if err != nil {
		return err
	}
	added := 0
	for _, name := range names { // Glob sorts, and names sort by time
		h, err := readHistoryFile(name)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		if have[h.At.Unix()] {
			continue
		}
		if _, err := db.Save(h.At, h.Prices); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		added++
	}
	if added > 0 {
		log.Printf("imported %d scans from %s", added, dir)
	}
	return nil
}

// exportHistory writes every scan of the database that has no file yet.
func exportHistory(dir string, db *store.Store) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	scans, err := db.Scans(1 << 30)
	if err != nil {
		return err
	}
	for _, sc := range scans {
		name := filepath.Join(dir, sc.At.UTC().Format(historyTimeFormat)+".json.gz")
		if _, err := os.Stat(name); err == nil {
			continue
		}
		at, prices, err := db.Prices(sc.ID)
		if err != nil {
			return err
		}
		if err := writeHistoryFile(name, historyFile{At: at, Prices: prices}); err != nil {
			return err
		}
		log.Printf("wrote %s", name)
	}
	return nil
}

func scanTimes(db *store.Store) (map[int64]bool, error) {
	scans, err := db.Scans(1 << 30)
	if err != nil {
		return nil, err
	}
	have := make(map[int64]bool, len(scans))
	for _, sc := range scans {
		have[sc.At.Unix()] = true
	}
	return have, nil
}

func readHistoryFile(name string) (historyFile, error) {
	var h historyFile
	f, err := os.Open(name)
	if err != nil {
		return h, err
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		return h, err
	}
	defer zr.Close()
	err = json.NewDecoder(zr).Decode(&h)
	return h, err
}

func writeHistoryFile(name string, h historyFile) error {
	tmp := name + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	zw := gzip.NewWriter(f)
	if err := json.NewEncoder(zw).Encode(h); err != nil {
		f.Close()
		return err
	}
	if err := zw.Close(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, name)
}
