package database

import (
	"database/sql"

	"github.com/google/uuid"
	sessions "github.com/relextm19/tracker.nvim/internal/sessions"
	"github.com/relextm19/tracker.nvim/internal/users"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		DB: db,
	}
}

func (s *Store) InsertSession(ses *sessions.Session) error {
	query := `
		INSERT INTO Sessions (FileName, ProjectName, LanguageName, StartTime, StartDate, EndTime, EndDate) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.DB.Exec(query,
		ses.FileName,
		ses.ProjectName,
		ses.LanguageName,
		ses.StartTime,
		ses.StartDate,
		ses.EndTime,
		ses.EndDate,
	)

	return err
}

func (s *Store) InsertUser(u *users.User) error {
	query := `
		INSERT INTO Users (Email, PasswordHash, Token) 
		VALUES (?, ?, ?)
	`

	_, err := s.DB.Exec(query,
		u.Email,
		u.PasswordHash,
		u.Token,
	)

	return err
}

func (s *Store) CheckLoginAttempt(cub *users.ClientUserBody) bool {
	query := `SELECT PasswordHash FROM Users WHERE Email = ?`
	var storedHash string
	err := s.DB.QueryRow(query, cub.Email).Scan(&storedHash)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(cub.Password))
	return err == nil
}

func (s *Store) GetUserToken(email string) (uuid.UUID, error) {
	var token uuid.UUID
	query := `SELECT Token FROM Users WHERE Email = ?`
	err := s.DB.QueryRow(query, email).Scan(&token)
	if err != nil {
		return uuid.Nil, err
	}

	return token, nil
}
