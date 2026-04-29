package domain

type SkillOption struct {
	ID        string
	Name      string
	Icon      string
	SortOrder int64
}

type SkillOptions struct {
	Categories        []SkillOption
	ProficiencyLevels []SkillOption
}

type SkillCategoryInput struct {
	ID   string
	Name string
	Icon string
}

type Skill struct {
	ID                 string
	Name               string
	CategoryID         string
	CategoryName       string
	Experience         string
	ProficiencyLevelID string
	ProficiencyName    string
	SortOrder          int64
}

type SkillInput struct {
	Name               string
	CategoryID         string
	Experience         string
	ProficiencyLevelID string
}
