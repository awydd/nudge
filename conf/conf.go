package conf

import (
	"path/filepath"
)

type Anniversary struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	OriginalDate     string `json:"original_date"`
	AdvanceDays      int    `json:"advance_days"`
	LastNotifiedYear int    `json:"last_notified_year"`
}

type Store struct {
	Anniversaries []Anniversary `json:"anniversaries"`
}

func Path() string {
	return filepath.Join("data", "anniversaries.json")
}
