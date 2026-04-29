package domain

type Project struct {
	ID           string
	Name         string
	Company      string
	StartYear    int64
	StartMonth   int64
	EndYear      *int64
	EndMonth     *int64
	Description  string
	Role         string
	TeamSize     string
	Technologies []string
	Phases       []string
	Achievements string
	IsDraft      bool
	SortOrder    int64
}

type ProjectInput struct {
	Name         string
	Company      string
	StartYear    int64
	StartMonth   int64
	EndYear      *int64
	EndMonth     *int64
	Description  string
	Role         string
	TeamSize     string
	Technologies []string
	Phases       []string
	Achievements string
	IsDraft      bool
}
