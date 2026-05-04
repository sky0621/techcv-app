package domain

type JobHistory struct {
	ID               string
	CompanyID        string
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
	CompanyID        string
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

type JobCompany struct {
	ID   string
	Name string
	URL  string
}

type JobHistoryOptions struct {
	EmploymentTypes []JobEmploymentType
	Companies       []JobCompany
}

type JobEmploymentTypeInput struct {
	ID        string
	Name      string
	SortOrder int64
}

type JobCompanyInput struct {
	ID   string
	Name string
	URL  string
}
