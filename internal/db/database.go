package database

import (
	"database/sql"

	"github.com/relextm19/tracker.nvim/internal/session"
)

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		DB: db,
	}
}

func (store *Store) InsertSession(s *session.Session) error {
	query := `
		INSERT INTO Sessions (FileName, ProjectName, LanguageName, StartTime, StartDate, EndTime, EndDate) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := store.DB.Exec(query,
		s.FileName,
		s.ProjectName,
		s.LanguageName,
		s.StartTime,
		s.StartDate,
		s.EndTime,
		s.EndDate,
	)

	return err
}
