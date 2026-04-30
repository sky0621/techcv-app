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
	SkillMasters      []SkillMaster
}

type SkillCategoryInput struct {
	ID   string
	Name string
	Icon string
}

type SkillMaster struct {
	ID           string
	Name         string
	CategoryID   string
	CategoryName string
	SortOrder    int64
}

type SkillMasterInput struct {
	ID         string
	Name       string
	CategoryID string
}

type Skill struct {
	ID                 string
	Name               string
	CategoryID         string
	CategoryName       string
	Experience         int64
	ProficiencyLevelID string
	ProficiencyName    string
	SortOrder          int64
}

type SkillInput struct {
	Name               string
	CategoryID         string
	Experience         int64
	ProficiencyLevelID string
}
