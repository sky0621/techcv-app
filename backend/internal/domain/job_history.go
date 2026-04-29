package domain

type JobHistory struct {
	ID             string
	Company        string
	StartYear      int64
	StartMonth     int64
	EndYear        *int64
	EndMonth       *int64
	EmploymentType string
	ProjectCount   int64
	SortOrder      int64
}

type JobHistoryInput struct {
	Company        string
	StartYear      int64
	StartMonth     int64
	EndYear        *int64
	EndMonth       *int64
	EmploymentType string
}
