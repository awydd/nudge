package conf

import (
	"path/filepath"
)

type ChannelType string

const (
	ChannelConsole ChannelType = "console"
	ChannelEmail   ChannelType = "email"
)

func (ct ChannelType) Valid() bool {
	switch ct {
	case ChannelConsole, ChannelEmail:
		return true
	default:
		return false
	}
}

type Anniversary struct {
	ID               string      `json:"id"`
	Title            string      `json:"title"`
	Description      string      `json:"description"`
	OriginalDate     string      `json:"original_date"`
	AdvanceDays      int         `json:"advance_days"`
	LastNotifiedYear int         `json:"last_notified_year"`
	Channel          ChannelType `json:"channel"` // console / email / ...
}

type Store struct {
	Anniversaries []Anniversary `json:"anniversaries"`
}

func Path() string {
	return filepath.Join(DataDir, "anniversaries.json")
}
