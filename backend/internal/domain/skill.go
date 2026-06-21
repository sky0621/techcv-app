package domain

type SkillOption struct {
	ID        string
	Name      string
	Icon      string
	SortOrder int64
}

type SkillOptions struct {
	Categories   []SkillOption
	SkillMasters []SkillMaster
}

type SkillCategoryInput struct {
	ID        string
	Name      string
	Icon      string
	SortOrder int64
}

type SkillMaster struct {
	ID           string
	Name         string
	CategoryID   string
	CategoryName string
	URL          string
}

type SkillMasterInput struct {
	Name       string
	CategoryID string
	URL        string
}

type Skill struct {
	ID            string
	SkillMasterID string
	Name          string
	CategoryID    string
	CategoryName  string
	Experience    int64
	SortOrder     int64
}

type SkillInput struct {
	SkillMasterID string
	Experience    int64
}
