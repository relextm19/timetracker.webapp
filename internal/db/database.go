package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
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
func (s *Store) InsertSession(ses *sessions.Session, token string) error {
	query := `
		INSERT INTO Sessions (UserToken, FileName, ProjectName, LanguageName, StartTime, StartDate, EndTime, EndDate)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.DB.Exec(query,
		token,
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

func (s *Store) fetchCategoryAggragatedData(column, token string) ([]AggregatedTime, error) {
	// the Sprintf call for db is fine cuz we only have 3 options so the db can still cache the query and shi
	query := fmt.Sprintf(`
			SELECT 
				%s, 
				SUM(EndTime - StartTime) as TotalTime 
			FROM Sessions 
			WHERE UserToken = ? 
			GROUP BY %s 
			ORDER BY TotalTime DESC;
		`, column, column)

	rows, err := s.DB.Query(query, token)
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

func (s *Store) fetchTimeAggregatedData(token string) ([]AggregatedTime, error) {
	// the filter here is optimal it doesnt fetch new rows for each call it operates on already fetched ones
	query := `
		SELECT 
			COALESCE(SUM(EndTime - StartTime) FILTER (WHERE StartDate >= date('now', 'start of day')), 0),
			COALESCE(SUM(EndTime - StartTime) FILTER (WHERE StartDate >= date('now', '-7 days')), 0),
			COALESCE(SUM(EndTime - StartTime) FILTER (WHERE StartDate >= date('now', '-1 month')), 0),
			COALESCE(SUM(EndTime - StartTime) FILTER (WHERE StartDate >= date('now', '-1 year')), 0)
		FROM Sessions 
		WHERE UserToken = ?;
	`

	var day, week, month, year int

	err := s.DB.QueryRow(query, token).Scan(&day, &week, &month, &year)
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

func (s *Store) GetDataForToken(token string) (*DashboardData, error) {
	data := &DashboardData{}

	var err error

	data.ByLanguage, err = s.fetchCategoryAggragatedData("LanguageName", token)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	data.ByProject, err = s.fetchCategoryAggragatedData("ProjectName", token)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	data.ByFile, err = s.fetchCategoryAggragatedData("FileName", token)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	data.ByTime, err = s.fetchTimeAggregatedData(token)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}
	return data, nil
}
