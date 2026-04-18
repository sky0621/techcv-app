package profiledb

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
	GithubUrl          string
	ZennUrl            string
	QiitaUrl           string
	WebsiteUrl         string
	PreferredWorkStyle string
	VisibilitySettings []byte
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type UpsertProfileParams struct {
	ID                 string
	UserID             string
	FullName           string
	Nickname           string
	Location           string
	Email              string
	Phone              string
	Summary            string
	GithubUrl          string
	ZennUrl            string
	QiitaUrl           string
	WebsiteUrl         string
	PreferredWorkStyle string
	VisibilitySettings []byte
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
