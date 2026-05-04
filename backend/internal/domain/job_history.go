package domain

type JobHistory struct {
	ID               string
	Company          string
	DisplayName      string
	StartYear        int64
	StartMonth       int64
	EndYear          *int64
	EndMonth         *int64
	EmploymentTypeID string
	EmploymentType   string
	ProjectCount     int64
}

type JobHistoryInput struct {
	Company          string
	DisplayName      string
	StartYear        int64
	StartMonth       int64
	EndYear          *int64
	EndMonth         *int64
	EmploymentTypeID string
}

type JobEmploymentType struct {
	ID        string
	Name      string
	SortOrder int64
}

type JobHistoryOptions struct {
	EmploymentTypes []JobEmploymentType
}

type JobEmploymentTypeInput struct {
	ID        string
	Name      string
	SortOrder int64
}
