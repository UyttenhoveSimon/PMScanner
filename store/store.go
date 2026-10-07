// Package store keeps every price scan in a small SQLite database, so past
// scans can be looked at again.
package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	_ "modernc.org/sqlite"
)

// ErrNotFound is returned when a scan does not exist.
var ErrNotFound = errors.New("scan not found")

// Scan is one run over all products.
type Scan struct {
	ID     int64
	At     time.Time
	Count  int
	Failed int
}

// Price is one product's result in a scan. Error is set when it could not be read.
type Price struct {
	Site        string
	Category    string
	Description string
	URL         string
	Price       decimal.Decimal
	Currency    string
	EURPerGram  decimal.Decimal
	Error       string
}

// Store is a SQLite database of scans.
type Store struct {
	db *sql.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS scans (
	id      INTEGER PRIMARY KEY,
	at      TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS prices (
	scan_id      INTEGER NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
	position     INTEGER NOT NULL,
	site         TEXT NOT NULL,
	category     TEXT NOT NULL,
	description  TEXT NOT NULL,
	url          TEXT NOT NULL,
	price        TEXT NOT NULL DEFAULT '',
	currency     TEXT NOT NULL DEFAULT '',
	eur_per_gram TEXT NOT NULL DEFAULT '',
	error        TEXT NOT NULL DEFAULT '',
	PRIMARY KEY (scan_id, position)
);
CREATE INDEX IF NOT EXISTS prices_url ON prices(url);
`

// Open opens or creates the database at path.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", "file:"+path+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("create schema: %w", err)
	}
	return &Store{db: db}, nil
}

// Close closes the database.
func (s *Store) Close() error { return s.db.Close() }

// Save stores a scan and its prices, keeping their order, and returns its ID.
func (s *Store) Save(at time.Time, prices []Price) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`INSERT INTO scans (at) VALUES (?)`, at.UTC().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare(`INSERT INTO prices
		(scan_id, position, site, category, description, url, price, currency, eur_per_gram, error)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	for i, p := range prices {
		price, perGram := "", ""
		if p.Error == "" {
			price, perGram = p.Price.String(), p.EURPerGram.String()
		}
		if _, err := stmt.Exec(id, i, p.Site, p.Category, p.Description, p.URL,
			price, p.Currency, perGram, p.Error); err != nil {
			return 0, err
		}
	}
	return id, tx.Commit()
}

// Scans lists the most recent scans first, at most limit of them.
func (s *Store) Scans(limit int) ([]Scan, error) {
	rows, err := s.db.Query(`
		SELECT s.id, s.at, COUNT(p.position), COALESCE(SUM(p.error != ''), 0)
		FROM scans s LEFT JOIN prices p ON p.scan_id = s.id
		GROUP BY s.id ORDER BY s.id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var scans []Scan
	for rows.Next() {
		var sc Scan
		var at string
		if err := rows.Scan(&sc.ID, &at, &sc.Count, &sc.Failed); err != nil {
			return nil, err
		}
		if sc.At, err = time.Parse(time.RFC3339, at); err != nil {
			return nil, fmt.Errorf("scan %d: %w", sc.ID, err)
		}
		scans = append(scans, sc)
	}
	return scans, rows.Err()
}

// Latest returns the ID of the most recent scan, or ErrNotFound.
func (s *Store) Latest() (int64, error) {
	var id int64
	err := s.db.QueryRow(`SELECT id FROM scans ORDER BY id DESC LIMIT 1`).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}
	return id, err
}

// Prices returns a scan's time and prices in the order they were saved.
func (s *Store) Prices(scanID int64) (time.Time, []Price, error) {
	var at string
	err := s.db.QueryRow(`SELECT at FROM scans WHERE id = ?`, scanID).Scan(&at)
	if errors.Is(err, sql.ErrNoRows) {
		return time.Time{}, nil, ErrNotFound
	}
	if err != nil {
		return time.Time{}, nil, err
	}
	when, err := time.Parse(time.RFC3339, at)
	if err != nil {
		return time.Time{}, nil, err
	}

	rows, err := s.db.Query(`
		SELECT site, category, description, url, price, currency, eur_per_gram, error
		FROM prices WHERE scan_id = ? ORDER BY position`, scanID)
	if err != nil {
		return time.Time{}, nil, err
	}
	defer rows.Close()
	var prices []Price
	for rows.Next() {
		var p Price
		var price, perGram string
		if err := rows.Scan(&p.Site, &p.Category, &p.Description, &p.URL,
			&price, &p.Currency, &perGram, &p.Error); err != nil {
			return time.Time{}, nil, err
		}
		if p.Error == "" {
			if p.Price, err = decimal.NewFromString(price); err != nil {
				return time.Time{}, nil, err
			}
			if p.EURPerGram, err = decimal.NewFromString(perGram); err != nil {
				return time.Time{}, nil, err
			}
		}
		prices = append(prices, p)
	}
	return when, prices, rows.Err()
}
