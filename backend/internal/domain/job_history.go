package domain

type JobHistory struct {
	ID             string
	Company        string
	StartDate      string
	EndDate        string
	EmploymentType string
	ProjectCount   int64
	SortOrder      int64
}

type JobHistoryInput struct {
	Company        string
	StartDate      string
	EndDate        string
	EmploymentType string
}
