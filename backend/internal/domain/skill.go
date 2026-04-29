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
