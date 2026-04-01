package database

import (
	"database/sql"
	"errors"
	"fmt"

	apikeys "github.com/relextm19/tracker.nvim/internal/apiKeys"
	"github.com/relextm19/tracker.nvim/internal/helpers"
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
	ByHour     []AggregatedTime `json:"byHour"`
	ByWeekday  []AggregatedTime `json:"byWeekday"`
	ByMonth    []AggregatedTime `json:"byMonth"`
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		DB: db,
	}
}

// InsertSession no point in making a new uuid from the header token so just pas it as string
func (s *Store) InsertSession(ses *sessions.Session, keyHash string) error {
	query := `
		INSERT INTO Sessions (KeyHash, FileName, ProjectName, LanguageName, StartTime, StartDate, EndTime, EndDate)
		SELECT ?, ?, ?, ?, ?, ?, ?, ?
	`

	_, err := s.DB.Exec(query,
		keyHash,
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
		INSERT INTO Users (Email, PasswordHash) 
		VALUES (?, ?)
	`

	_, err := s.DB.Exec(query,
		u.Email,
		u.PasswordHash,
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

func (s *Store) GetUserIDForToken(token string) (string, error) {
	var userID string
	tokenHash, err := helpers.GetHashFromUUID([]byte(token))
	if err != nil {
		return "", err
	}

	query := `SELECT UserID FROM Tokens WHERE TokenHash = ?`

	err = s.DB.QueryRow(query, tokenHash).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (s *Store) InsertToken(userID string, token string) error {
	tokenHash, err := helpers.GetHashFromUUID([]byte(token))
	if err != nil {
		return err
	}

	query := `
		INSERT INTO Tokens (UserID, TokenHash)
		VALUES (?, ?)
	`

	_, err = s.DB.Exec(query, userID, tokenHash)
	return err
}

func (s *Store) GetUserIDByEmail(email string) (string, error) {
	var userID string

	query := `SELECT ID FROM Users WHERE Email = ?`
	err := s.DB.QueryRow(query, email).Scan(&userID)
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

func (s *Store) fetchCategoryAggragatedData(column, keyHash string) ([]AggregatedTime, error) {
	// the Sprintf call for db is fine cuz we only have 3 options so the db can still cache the query and shi
	query := fmt.Sprintf(`
			SELECT 
				%s, 
				SUM(EndTime - StartTime) as TotalTime 
			FROM Sessions 
			WHERE KeyHash = ?
			GROUP BY %s 
			ORDER BY TotalTime DESC;
		`, column, column)

	rows, err := s.DB.Query(query, keyHash)
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

func (s *Store) fetchTimeAggregatedData(keyHash string, start string, end string) ([]AggregatedTime, error) {
	// consider sessions starting on one day and ending at another
	query := `
    SELECT 
        startDate,
        COALESCE(SUM(EndTime - StartTime), 0)
    FROM Sessions 
	    WHERE KeyHash = ? AND startDate >= date(?, ?)
    GROUP BY startDate;
`
	rows, err := s.DB.Query(query, keyHash, start, end)
	if err != nil {
		return nil, err
	}

	var result []AggregatedTime
	for rows.Next() {
		var item AggregatedTime
		if err := rows.Scan(&item.Name, &item.TotalTime); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *Store) fetchBucketAggregatedData(keyHash string, format string, labels []string) ([]AggregatedTime, error) {
	query := `
		SELECT CAST(strftime(?, datetime(StartTime, 'unixepoch')) AS INTEGER), COALESCE(SUM(EndTime - StartTime), 0)
		FROM Sessions
		WHERE KeyHash = ?
		GROUP BY 1
		ORDER BY 1;
	`

	rows, err := s.DB.Query(query, format, keyHash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]AggregatedTime, len(labels))
	for i, label := range labels {
		result[i] = AggregatedTime{Name: label, TotalTime: 0}
	}

	for rows.Next() {
		var bucket int
		var totalTime int
		if err := rows.Scan(&bucket, &totalTime); err != nil {
			return nil, err
		}
		if format == "%m" {
			bucket--
		}

		if bucket >= 0 && bucket < len(result) {
			result[bucket].TotalTime = totalTime
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

var ErrAggregatingData = errors.New("error aggregating data")

func (s *Store) GetSessionDataForKeyHash(keyHash string) (*DashboardData, error) {
	data := &DashboardData{}

	var err error

	data.ByLanguage, err = s.fetchCategoryAggragatedData("LanguageName", keyHash)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	data.ByProject, err = s.fetchCategoryAggragatedData("ProjectName", keyHash)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	data.ByFile, err = s.fetchCategoryAggragatedData("FileName", keyHash)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	data.ByTime, err = s.fetchTimeAggregatedData(keyHash, "now", "start of month")
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	hourLabels := []string{
		"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11",
		"12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23",
	}
	data.ByHour, err = s.fetchBucketAggregatedData(keyHash, "%H", hourLabels)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	weekdayLabels := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	data.ByWeekday, err = s.fetchBucketAggregatedData(keyHash, "%w", weekdayLabels)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	monthLabels := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	data.ByMonth, err = s.fetchBucketAggregatedData(keyHash, "%m", monthLabels)
	if err != nil {
		return nil, errors.Join(ErrAggregatingData, err)
	}

	return data, nil
}

func (s *Store) GetKeyHashes(userID string) ([]string, error) {
	query := `SELECT KeyHash FROM APIKeys WHERE UserID = ?`

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []string{}
	for rows.Next() {
		var keyHash string
		if err := rows.Scan(&keyHash); err != nil {
			return nil, err
		}
		result = append(result, keyHash)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Store) GetSessionDataGroupedByKeyHash(userID string) (map[string]DashboardData, error) {
	keyHashes, err := s.GetKeyHashes(userID)
	if err != nil {
		return nil, err
	}

	res := map[string]DashboardData{}
	for _, keyHash := range keyHashes {
		data, err := s.GetSessionDataForKeyHash(keyHash)
		if err != nil {
			return nil, err
		}
		res[keyHash] = *data
	}

	return res, nil
}

var ErrNoRowsAffected = errors.New("no rows affected")

func (s *Store) InsertAPIKey(userID string, ak *apikeys.APIKey) (int, int, error) {
	// FIXME: The db unique doesnt prevent multiple of the same api key cuz of the salt or sth of a hash
	query := `
        INSERT INTO ApiKeys (UserID, Name, KeyHash)
        VALUES(?,?,?)
    `

	res, err := s.DB.Exec(query, userID, ak.Name, ak.KeyHash)
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
