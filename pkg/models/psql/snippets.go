package psql

import (
	"database/sql"
	"errors"

	"krutolevichsemyon.life/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
			VALUES($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + $3 * INTERVAL '1' DAY) RETURNING id`
	lastInsertId := 0
	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets 
			WHERE expires > CURRENT_TIMESTAMP and id = $1;`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Lates() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets 
			WHERE expires > CURRENT_TIMESTAMP ORDER BY created DESC LIMIT 10;`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []*models.Snippet

	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
