package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	apikeys "github.com/relextm19/tracker.nvim/internal/apiKeys"
	sessions "github.com/relextm19/tracker.nvim/internal/sessions"
	"github.com/relextm19/tracker.nvim/internal/users"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	DB *sql.DB
}

type AggregatedTime struct {
	Name      string `json:"name"`
	TotalTime int    `json:"totalTime"` // Time in seconds
}

type DashboardData struct {
	ByLanguage []AggregatedTime `json:"byLanguage"`
	ByProject  []AggregatedTime `json:"byProject"`
	ByFile     []AggregatedTime `json:"byFile"`
	ByTime     []AggregatedTime `json:"byTime"`
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		DB: db,
	}
}

// InsertSession no point in making a new uuid from the header token so just pas it as string
func (s *Store) InsertSession(ses *sessions.Session, userID string) error {
	query := `
		INSERT INTO Sessions (UserID, FileName, ProjectName, LanguageName, StartTime, StartDate, EndTime, EndDate)
		SELECT ?, ?, ?, ?, ?, ?, ?, ?
	`

	_, err := s.DB.Exec(query,
		userID,
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

func (s *Store) CheckLoginAttempt(cub *users.RequestUserBody) bool {
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

func (s *Store) GetUserIDForToken(token string) (string, error) {
	var userID string

	query := `SELECT ID FROM Users WHERE Token = ?`

	err := s.DB.QueryRow(query, token).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (s *Store) GetUserIDForKeyHash(keyHash string) (string, error) {
	var userID string

	query := `SELECT UserID FROM APIKeys WHERE KeyHash = ?`

	err := s.DB.QueryRow(query, keyHash).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (s *Store) fetchCategoryAggragatedData(column, userID string) ([]AggregatedTime, error) {
	// the Sprintf call for db is fine cuz we only have 3 options so the db can still cache the query and shi
	query := fmt.Sprintf(`
			SELECT 
				%s, 
				SUM(EndTime - StartTime) as TotalTime 
			FROM Sessions 
			WHERE UserID = ?
			GROUP BY %s 
			ORDER BY TotalTime DESC;
		`, column, column)

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []AggregatedTime
	for rows.Next() {
		var item AggregatedTime
		if err := rows.Scan(&item.Name, &item.TotalTime); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

type TimePeriod struct {
	Name     string
	Modifier string
}

func (s *Store) fetchTimeAggregatedData(userID string) ([]AggregatedTime, error) {
	query := `
		SELECT 
			COALESCE(SUM(EndTime - StartTime) FILTER (WHERE StartDate >= date('now', 'start of day')), 0),
			COALESCE(SUM(EndTime - StartTime) FILTER (WHERE StartDate >= date('now', '-7 days')), 0),
			COALESCE(SUM(EndTime - StartTime) FILTER (WHERE StartDate >= date('now', '-1 month')), 0),
			COALESCE(SUM(EndTime - StartTime) FILTER (WHERE StartDate >= date('now', '-1 year')), 0)
		FROM Sessions 
		WHERE UserID = ?;
	`

	var day, week, month, year int

	err := s.DB.QueryRow(query, userID).Scan(&day, &week, &month, &year)
	if err != nil {
		return nil, err
	}

	results := []AggregatedTime{
		{Name: "day", TotalTime: day},
		{Name: "week", TotalTime: week},
		{Name: "month", TotalTime: month},
		{Name: "year", TotalTime: year},
	}

	return results, nil
}

var ErrAggregatingData = errors.New("error aggregating data")

func (s *Store) GetSessionDataForToken(userID string) (*DashboardData, error) {
	data := &DashboardData{}

	var err error

	data.ByLanguage, err = s.fetchCategoryAggragatedData("LanguageName", userID)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	data.ByProject, err = s.fetchCategoryAggragatedData("ProjectName", userID)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	data.ByFile, err = s.fetchCategoryAggragatedData("FileName", userID)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	data.ByTime, err = s.fetchTimeAggregatedData(userID)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}
	return data, nil
}

var ErrNoRowsAffected = errors.New("no rows affected")

func (s *Store) InsertAPIKey(userID string, ak *apikeys.APIKey) (int, int, error) {
	// FIXME: The db unique doesnt prevent multiple of the same api key cuz of the salt or sth of a hash
	query := `
        INSERT INTO ApiKeys (UserID, Name, KeyHash)
        VALUES(?,?,?)
    `

	res, err := s.DB.Exec(query, ak.Name, ak.KeyHash, userID)
	if err != nil {
		return 0, 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, 0, err
	}

	if rows == 0 {
		return 0, 0, ErrNoRowsAffected
	}

	newID, err := res.LastInsertId()
	if err != nil {
		return 0, 0, err
	}
	createdAt := 0
	err = s.DB.QueryRow("SELECT CreatedAt FROM ApiKeys WHERE ID = ?", newID).Scan(&createdAt)
	if err != nil {
		return 0, 0, err
	}

	return int(newID), createdAt, nil
}

func (s *Store) DeleteAPIKey(id, userID string) error {
	query := `DELETE FROM ApiKeys WHERE ID = ? AND UserID = ?`

	res, err := s.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return ErrNoRowsAffected
	}

	return nil
}

func (s *Store) GetAPIKeys(userID string) ([]apikeys.APIKey, error) {
	query := `SELECT ID, Name, CreatedAt, KeyHash FROM ApiKeys WHERE UserID = ?`
	res := []apikeys.APIKey{}

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		ak := apikeys.APIKey{}
		if err := rows.Scan(&ak.ID, &ak.Name, &ak.CreatedAt, &ak.KeyHash); err != nil {
			return nil, err
		}
		res = append(res, ak)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
