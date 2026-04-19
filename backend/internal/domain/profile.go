package domain

import "time"

type Profile struct {
	ID                 string
	UserID             string
	FullName           string
	Nickname           string
	Location           string
	Email              string
	Phone              string
	Summary            string
	GitHubURL          string
	ZennURL            string
	QiitaURL           string
	WebsiteURL         string
	PreferredWorkStyle string
	VisibilitySettings map[string]any
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
