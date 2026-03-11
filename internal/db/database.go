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
	ByDay      int              `json:"byDay"`
	ByWeek     int              `json:"byWeek"`
	ByMonth    int              `json:"byMonth"`
	ByYear     int              `json:"byYear"`
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

var (
	PeriodDay   = TimePeriod{Name: "day", Modifier: "start of day"}
	PeriodWeek  = TimePeriod{Name: "week", Modifier: "-7 days"}
	PeriodMonth = TimePeriod{Name: "month", Modifier: "-1 month"}
	PeriodYear  = TimePeriod{Name: "year", Modifier: "-1 year"}
)

func (s *Store) fetchTimeAggregatedData(timePeriod TimePeriod, token string) (int, error) {
	query := fmt.Sprintf(`
			SELECT 
			COALESCE(SUM(EndTime - StartTime), 0) as TotalTime 
			FROM Sessions 
			WHERE UserToken = ? 
			AND StartDate >= date('now', '%s');
		`, timePeriod.Modifier)

	result := 0
	err := s.DB.QueryRow(query, token).Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
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

	data.ByDay, err = s.fetchTimeAggregatedData(PeriodDay, token)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	data.ByWeek, err = s.fetchTimeAggregatedData(PeriodWeek, token)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	data.ByMonth, err = s.fetchTimeAggregatedData(PeriodMonth, token)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	data.ByYear, err = s.fetchTimeAggregatedData(PeriodYear, token)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	return data, nil
}
