package database

import (
	"database/sql"

	sessions "github.com/relextm19/tracker.nvim/internal/sessions"
	"github.com/relextm19/tracker.nvim/internal/users"
)

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		DB: db,
	}
}

func (store *Store) InsertSession(s *sessions.Session) error {
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

func (store *Store) InsertUser(u *users.User) error {
	query := `
		INSERT INTO Users (Email, PasswordHash, Token) 
		VALUES (?, ?, ?)
	`

	_, err := store.DB.Exec(query,
		u.Email,
		u.PasswordHash,
		u.Token,
	)

	return err
}
