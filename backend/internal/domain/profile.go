package domain

import "time"

type Profile struct {
	ID                 string
	UserID             string
	FullName           string
	Nickname           string
	Location           string
	Email              string
	Summary            string
	GitHubURL          string
	ZennURL            string
	QiitaURL           string
	WebsiteURL         string
	Occupation         string
	EmploymentType     string
	PreferredWorkStyle string
	VisibilitySettings map[string]any
	Qualifications     []Qualification
	Links              []ProfileLink
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Qualification struct {
	ID           string
	Name         string
	AcquiredDate string
	Organization string
	URL          string
	Memo         string
}

type ProfileLinkMaster struct {
	ID          string
	Key         string
	Name        string
	Icon        string
	Placeholder string
	SortOrder   int64
}

type ProfileLinkMasterInput struct {
	ID          string
	Key         string
	Name        string
	Icon        string
	Placeholder string
	SortOrder   int64
}

type ProfileLink struct {
	ID           string
	LinkMasterID string
	Key          string
	Name         string
	Icon         string
	Placeholder  string
	URL          string
	SortOrder    int64
}
