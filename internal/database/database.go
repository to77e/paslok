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
			create table if not exists paslok
			(
				id int primary key autoincrement,
				username text not null,
				password text not null,
				comment text,
				created_at datetime default current_timestamp not null,
				deleted_at datetime,
				service text default '' not null
			);`
	)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	_, err = db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("exec: %w", err)
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}
	return nil
}

func (d *Database) Create(req *models.CreatePasswordRequest) error {
	const (
		queryInsert = `
			insert into paslok (username, password, comment, service)
			values ($1, $2, $3, $4);`
	)
	d.Lock()
	defer d.Unlock()
	_, err := d.db.Exec(queryInsert, req.Username, req.Password, req.Comment, req.Service)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}

func (d *Database) Read(req *models.ReadPasswordRequest) (string, error) {
	const (
		query = `
			select 
				password 
			from paslok
			where deleted_at is null
				and ($1 = 0 or id = $1)
				and ($2 = '' or service = $2);`
	)
	d.Lock()
	defer d.Unlock()
	var password string
	if err := d.db.QueryRow(query, req.Id, req.Service).Scan(&password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrorNotFoundName
		}
		return "", fmt.Errorf("exec: %w", err)
	}
	return password, nil
}

func (d *Database) List(req *models.ListPasswordsRequest) ([]models.Resource, error) {
	const (
		query = `
			select
				id,
				service,
				username,
				comment,
				date(created_at) as created_at
			from paslok
			where deleted_at is null
				and ($1 = '' or lower(service) like format('%%%s%%', $1))
			order by service, id;`
	)
	d.Lock()
	defer d.Unlock()
	rows, err := d.db.Query(query, req.SearchTerm) //nolint:rowserrcheck
	if err != nil {
		return nil, fmt.Errorf("exec: %w", err)
	}
	defer rows.Close() //nolint: errcheck
	var resources []models.Resource
	for rows.Next() {
		var resource models.Resource
		if err = rows.Scan(
			&resource.Id,
			&resource.Service,
			&resource.Username,
			&resource.Comment,
			&resource.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

func (d *Database) Delete(req *models.DeletePasswordRequest) error {
	const (
		query = `
			update paslok 
			set deleted_at = current_timestamp
			where deleted_at is null
				and ($1 = 0 or id = $1)
				and ($2 = '' or service = $2);`
	)
	d.Lock()
	defer d.Unlock()
	_, err := d.db.Exec(query, req.Id, req.Service)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	return nil
}
