package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	var id int
	stmt := `INSERT INTO snippets (title, content, created, expires)
             VALUES ($1, $2, CURRENT_TIMESTAMP AT TIME ZONE 'UTC', CURRENT_TIMESTAMP AT TIME ZONE 'UTC' + ($3 || ' days')::INTERVAL)
             RETURNING id`

	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id,title,content,created,expires FROM snippets WHERE expires > CURRENT_TIMESTAMP AT TIME ZONE 'UTC' AND id = $1 `
	var s Snippet

	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecords
		} else {
			return Snippet{}, err
		}
	}
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id,title,content,created,expires FROM SNIPPETS WHERE expires > CURRENT_TIMESTAMP AT TIME ZONE 'UTC' ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Snippet{}, ErrNoRecords
		} else {
			return []Snippet{}, err
		}
	}
	defer rows.Close()
	var snippets []Snippet
	for rows.Next() {
		var s Snippet
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return []Snippet{}, err
		}
		snippets = append(snippets, s)
	}
	if rows.Err(); err != nil {
		return []Snippet{}, err
	}
	return snippets, nil
}
