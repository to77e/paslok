package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/to77e/paslok/internal/models"
)

type Database struct {
	sync.Mutex
	db *sql.DB
}

func New(path string) (*Database, error) {
	const (
		query = `
			CREATE TABLE IF NOT EXISTS paslok (
				id         INTEGER PRIMARY KEY AUTOINCREMENT,
				name       TEXT     NOT NULL,
				password   TEXT     NOT NULL,
				comment    TEXT,
				created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
				deleted_at DATETIME);`
	)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("exec create: %w", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}
	return nil
}

func (d *Database) Create(name, password, comment string) error {
	const (
		queryInsert = `
			INSERT INTO paslok (name, password, comment)
			VALUES (?, ?, ?);`

		queryCheck = `
			SELECT COUNT(*) FROM paslok WHERE name = ? AND deleted_at IS NULL;`
	)
	d.Lock()
	defer d.Unlock()

	var count int
	err := d.db.QueryRow(queryCheck, name).Scan(&count)
	if err != nil {
		return fmt.Errorf("exec select: %w", err)
	}
	if count > 0 {
		return models.ErrorAlreadyExistsName
	}

	_, err = d.db.Exec(queryInsert, name, password, comment)
	if err != nil {
		return fmt.Errorf("exec insert: %w", err)
	}
	return nil
}

func (d *Database) Read(name string) (string, error) {
	const (
		query = `
			SELECT password FROM paslok
			WHERE name = ? AND deleted_at IS NULL;`
	)
	d.Lock()
	defer d.Unlock()
	var password string
	err := d.db.QueryRow(query, name).Scan(&password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrorNotFoundName
		}
		return "", fmt.Errorf("exec select: %w", err)
	}
	return password, nil
}

func (d *Database) List() ([]models.Resource, error) {
	const (
		query = `
			SELECT name, comment FROM paslok
			WHERE deleted_at IS NULL
			ORDER BY id;`
	)
	d.Lock()
	defer d.Unlock()
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("exec select: %w", err)
	}
	defer rows.Close() //nolint: errcheck

	var resources []models.Resource
	for rows.Next() {
		var name, comment string
		if err := rows.Scan(&name, &comment); err != nil {
			return nil, fmt.Errorf("scan select: %w", err)
		}
		resources = append(resources, models.Resource{
			Name:    name,
			Comment: comment,
		})
	}
	return resources, nil
}

func (d *Database) Update(name, password, comment string) error {
	const (
		queryUpdate = `
			UPDATE paslok SET deleted_at = CURRENT_TIMESTAMP
			WHERE name = ? AND deleted_at IS NULL;
			INSERT INTO paslok (name, password, comment) VALUES (?, ?, ?);`

		queryCheck = `
			SELECT COUNT(*) FROM paslok WHERE name = ? AND deleted_at IS NULL;`
	)
	d.Lock()
	defer d.Unlock()

	var count int
	err := d.db.QueryRow(queryCheck, name).Scan(&count)
	if err != nil {
		return fmt.Errorf("exec select: %w", err)
	}
	if count == 0 {
		return models.ErrorNotFoundName
	}
	_, err = d.db.Exec(queryUpdate, name, name, password, comment)
	if err != nil {
		return fmt.Errorf("exec update: %w", err)
	}
	return nil
}

func (d *Database) Delete(name string) error {
	const (
		query = `
			UPDATE paslok SET deleted_at = CURRENT_TIMESTAMP
			WHERE name = ? AND deleted_at IS NULL;`
	)
	d.Lock()
	defer d.Unlock()
	_, err := d.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("exec delete: %w", err)
	}
	return nil
}
