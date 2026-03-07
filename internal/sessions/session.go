package session

import (
	"errors"

	"github.com/relextm19/tracker.nvim/internal/helpers"
)

type Session struct {
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	StartTime    uint64 `json:"startTime"`
	EndTime      uint64 `json:"endTime"`
	LanguageName string `json:"languageName"`
	ProjectName  string `json:"projectName"`
	FileName     string `json:"fileName"`
}

func NewSession() *Session {
	return &Session{}
}

func (s *Session) Valid() error {
	if !helpers.ValidStringField(s.StartDate) {
		return errors.New("startDate is required")
	}
	if !helpers.ValidStringField(s.EndDate) {
		return errors.New("endDate is required")
	}
	if !helpers.ValidStringField(s.ProjectName) {
		return errors.New("projectName is required")
	}
	if !helpers.ValidStringField(s.FileName) {
		return errors.New("fileName is required")
	}
	if !helpers.ValidStringField(s.LanguageName) {
		return errors.New("languageName is required")
	}

	if s.StartTime == 0 {
		return errors.New("startTime cannot be zero")
	}
	if s.EndTime == 0 {
		return errors.New("endTime cannot be zero")
	}

	if s.EndTime < s.StartTime {
		return errors.New("endTime cannot be before startTime")
	}

	return nil
}
