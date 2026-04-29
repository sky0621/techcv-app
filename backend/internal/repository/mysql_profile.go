package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sky0621/techcv-app/backend/internal/domain"
	dbgen "github.com/sky0621/techcv-app/backend/internal/repository/db"
)

type SQLiteProfileRepository struct {
	db      *sql.DB
	queries *dbgen.Queries
}

func NewSQLiteProfileRepository(ctx context.Context, dsn string, schemaPath string) (*SQLiteProfileRepository, error) {
	if err := ensureSQLiteDirectory(dsn); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(1)

	repository := &SQLiteProfileRepository{
		db:      db,
		queries: dbgen.New(db),
	}

	if err := repository.configure(ctx); err != nil {
		_ = repository.Close()
		return nil, err
	}

	if err := repository.applySchema(ctx, schemaPath); err != nil {
		_ = repository.Close()
		return nil, err
	}

	return repository, nil
}

func (r *SQLiteProfileRepository) Close() error {
	if r.db == nil {
		return nil
	}

	return r.db.Close()
}

func (r *SQLiteProfileRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *SQLiteProfileRepository) Get(ctx context.Context) (*domain.Profile, error) {
	row, err := r.queries.GetProfileByUserID(ctx, "user_01")
	if err == nil {
		qualifications, err := r.queries.ListQualificationsByProfileID(ctx, row.ID)
		if err != nil {
			return nil, fmt.Errorf("query profile qualifications: %w", err)
		}

		return toDomainProfile(row, qualifications)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("query profile: %w", err)
	}

	now := time.Now().UTC()
	profile := domain.Profile{
		ID:                 "profile_01",
		UserID:             "user_01",
		VisibilitySettings: map[string]any{"email": true, "location": true},
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	return r.Save(ctx, &profile)
}

func (r *SQLiteProfileRepository) Save(ctx context.Context, profile *domain.Profile) (*domain.Profile, error) {
	visibilitySettings := sanitizeVisibilitySettings(profile.VisibilitySettings)

	visibilityBytes, err := json.Marshal(visibilitySettings)
	if err != nil {
		return nil, fmt.Errorf("encode visibility settings: %w", err)
	}

	now := time.Now().UTC()
	createdAt := profile.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin profile transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	queries := r.queries.WithTx(tx)

	err = queries.UpsertProfile(ctx, dbgen.UpsertProfileParams{
		ID:                 profile.ID,
		UserID:             profile.UserID,
		FullName:           profile.FullName,
		Nickname:           profile.Nickname,
		Location:           profile.Location,
		Email:              profile.Email,
		Summary:            profile.Summary,
		GithubUrl:          profile.GitHubURL,
		ZennUrl:            profile.ZennURL,
		QiitaUrl:           profile.QiitaURL,
		WebsiteUrl:         profile.WebsiteURL,
		Occupation:         profile.Occupation,
		EmploymentType:     profile.EmploymentType,
		PreferredWorkStyle: profile.PreferredWorkStyle,
		VisibilitySettings: string(visibilityBytes),
		CreatedAt:          createdAt,
		UpdatedAt:          now,
	})
	if err != nil {
		return nil, fmt.Errorf("save profile: %w", err)
	}

	if err := queries.DeleteQualificationsByProfileID(ctx, profile.ID); err != nil {
		return nil, fmt.Errorf("delete profile qualifications: %w", err)
	}

	for index, qualification := range profile.Qualifications {
		if err := queries.InsertQualification(ctx, dbgen.InsertQualificationParams{
			ID:           qualificationID(qualification, index, now),
			ProfileID:    profile.ID,
			Name:         qualification.Name,
			AcquiredDate: qualification.AcquiredDate,
			Organization: qualification.Organization,
			Url:          qualification.URL,
			Memo:         qualification.Memo,
			SortOrder:    int64(index),
			CreatedAt:    now,
			UpdatedAt:    now,
		}); err != nil {
			return nil, fmt.Errorf("insert profile qualification: %w", err)
		}
	}

	row, err := queries.GetProfileByUserID(ctx, profile.UserID)
	if err != nil {
		return nil, fmt.Errorf("reload profile: %w", err)
	}

	qualifications, err := queries.ListQualificationsByProfileID(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("reload profile qualifications: %w", err)
	}

	result, err := toDomainProfile(row, qualifications)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit profile transaction: %w", err)
	}

	return result, nil
}

func (r *SQLiteProfileRepository) configure(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("enable foreign keys: %w", err)
	}
	if _, err := r.db.ExecContext(ctx, "PRAGMA busy_timeout = 5000"); err != nil {
		return fmt.Errorf("set busy timeout: %w", err)
	}

	return nil
}

func (r *SQLiteProfileRepository) applySchema(ctx context.Context, schemaPath string) error {
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("read schema file: %w", err)
	}

	statements := strings.Split(string(schemaBytes), ";")
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		if _, err := r.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("apply schema statement %q: %w", statement, err)
		}
	}

	if err := r.dropProfilePhoneColumn(ctx); err != nil {
		return err
	}
	if err := r.addProfileTextColumn(ctx, "occupation"); err != nil {
		return err
	}
	if err := r.addProfileTextColumn(ctx, "employment_type"); err != nil {
		return err
	}
	if err := r.addSkillCategoryIconColumn(ctx); err != nil {
		return err
	}
	if err := r.seedSkillOptions(ctx); err != nil {
		return err
	}
	if err := r.seedSkills(ctx); err != nil {
		return err
	}
	if err := r.seedJobEmploymentTypes(ctx); err != nil {
		return err
	}
	if err := r.migrateJobHistoryDateColumns(ctx); err != nil {
		return err
	}
	if err := r.seedJobHistories(ctx); err != nil {
		return err
	}
	if err := r.migrateProjectDateColumns(ctx); err != nil {
		return err
	}
	if err := r.seedProjects(ctx); err != nil {
		return err
	}
	if err := r.backfillEmptySkillCategoryIcons(ctx); err != nil {
		return err
	}

	return r.addQualificationURLColumn(ctx)
}

func (r *SQLiteProfileRepository) ListSkillOptions(ctx context.Context) (*domain.SkillOptions, error) {
	categories, err := r.queries.ListSkillCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("query skill categories: %w", err)
	}

	proficiencyLevels, err := r.queries.ListSkillProficiencyLevels(ctx)
	if err != nil {
		return nil, fmt.Errorf("query skill proficiency levels: %w", err)
	}

	return &domain.SkillOptions{
		Categories:        toDomainSkillCategoryOptions(categories),
		ProficiencyLevels: toDomainSkillProficiencyLevelOptions(proficiencyLevels),
	}, nil
}

func (r *SQLiteProfileRepository) CreateSkillCategory(ctx context.Context, input domain.SkillCategoryInput) (*domain.SkillOption, error) {
	row, err := r.queries.InsertSkillCategory(ctx, dbgen.InsertSkillCategoryParams{
		ID:   input.ID,
		Name: input.Name,
		Icon: input.Icon,
	})
	if err != nil {
		return nil, fmt.Errorf("insert skill category: %w", err)
	}

	return toDomainSkillCategoryOption(row), nil
}

func (r *SQLiteProfileRepository) UpdateSkillCategory(ctx context.Context, id string, input domain.SkillCategoryInput) (*domain.SkillOption, error) {
	row, err := r.queries.UpdateSkillCategory(ctx, dbgen.UpdateSkillCategoryParams{
		ID:   id,
		Name: input.Name,
		Icon: input.Icon,
	})
	if err != nil {
		return nil, fmt.Errorf("update skill category: %w", err)
	}

	return toDomainSkillCategoryOption(row), nil
}

func (r *SQLiteProfileRepository) ListSkills(ctx context.Context) ([]domain.Skill, error) {
	rows, err := r.queries.ListSkills(ctx)
	if err != nil {
		return nil, fmt.Errorf("query skills: %w", err)
	}

	return toDomainSkills(rows), nil
}

func (r *SQLiteProfileRepository) CreateSkill(ctx context.Context, input domain.SkillInput) (*domain.Skill, error) {
	now := time.Now().UTC()
	row, err := r.queries.InsertSkill(ctx, dbgen.InsertSkillParams{
		ID:                 skillID(now),
		Name:               input.Name,
		CategoryID:         input.CategoryID,
		Experience:         input.Experience,
		ProficiencyLevelID: input.ProficiencyLevelID,
	})
	if err != nil {
		return nil, fmt.Errorf("insert skill: %w", err)
	}

	created, err := r.queries.GetSkill(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("reload skill: %w", err)
	}

	return toDomainSkill(created), nil
}

func (r *SQLiteProfileRepository) UpdateSkill(ctx context.Context, id string, input domain.SkillInput) (*domain.Skill, error) {
	row, err := r.queries.UpdateSkill(ctx, dbgen.UpdateSkillParams{
		ID:                 id,
		Name:               input.Name,
		CategoryID:         input.CategoryID,
		Experience:         input.Experience,
		ProficiencyLevelID: input.ProficiencyLevelID,
	})
	if err != nil {
		return nil, fmt.Errorf("update skill: %w", err)
	}

	updated, err := r.queries.GetSkill(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("reload skill: %w", err)
	}

	return toDomainSkill(updated), nil
}

func (r *SQLiteProfileRepository) DeleteSkill(ctx context.Context, id string) error {
	rowsAffected, err := r.queries.DeleteSkill(ctx, id)
	if err != nil {
		return fmt.Errorf("delete skill: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SQLiteProfileRepository) ListJobHistories(ctx context.Context) ([]domain.JobHistory, error) {
	rows, err := r.queries.ListJobHistories(ctx)
	if err != nil {
		return nil, fmt.Errorf("query job histories: %w", err)
	}

	return toDomainJobHistories(rows), nil
}

func (r *SQLiteProfileRepository) CreateJobHistory(ctx context.Context, input domain.JobHistoryInput) (*domain.JobHistory, error) {
	now := time.Now().UTC()
	row, err := r.queries.InsertJobHistory(ctx, dbgen.InsertJobHistoryParams{
		ID:               jobHistoryID(now),
		Company:          input.Company,
		DisplayName:      input.DisplayName,
		StartYear:        input.StartYear,
		StartMonth:       input.StartMonth,
		EndYear:          toNullInt64(input.EndYear),
		EndMonth:         toNullInt64(input.EndMonth),
		EmploymentTypeID: input.EmploymentTypeID,
	})
	if err != nil {
		return nil, fmt.Errorf("insert job history: %w", err)
	}

	created, err := r.queries.GetJobHistory(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("reload job history: %w", err)
	}

	return toDomainJobHistory(created), nil
}

func (r *SQLiteProfileRepository) UpdateJobHistory(ctx context.Context, id string, input domain.JobHistoryInput) (*domain.JobHistory, error) {
	row, err := r.queries.UpdateJobHistory(ctx, dbgen.UpdateJobHistoryParams{
		ID:               id,
		Company:          input.Company,
		DisplayName:      input.DisplayName,
		StartYear:        input.StartYear,
		StartMonth:       input.StartMonth,
		EndYear:          toNullInt64(input.EndYear),
		EndMonth:         toNullInt64(input.EndMonth),
		EmploymentTypeID: input.EmploymentTypeID,
	})
	if err != nil {
		return nil, fmt.Errorf("update job history: %w", err)
	}

	updated, err := r.queries.GetJobHistory(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("reload job history: %w", err)
	}

	return toDomainJobHistory(updated), nil
}

func (r *SQLiteProfileRepository) DeleteJobHistory(ctx context.Context, id string) error {
	rowsAffected, err := r.queries.DeleteJobHistory(ctx, id)
	if err != nil {
		return fmt.Errorf("delete job history: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SQLiteProfileRepository) ListProjects(ctx context.Context) ([]domain.Project, error) {
	rows, err := r.queries.ListProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("query projects: %w", err)
	}

	return toDomainProjects(rows)
}

func (r *SQLiteProfileRepository) CreateProject(ctx context.Context, input domain.ProjectInput) (*domain.Project, error) {
	now := time.Now().UTC()
	technologies, err := encodeStringSlice(input.Technologies)
	if err != nil {
		return nil, err
	}
	phases, err := encodeStringSlice(input.Phases)
	if err != nil {
		return nil, err
	}

	row, err := r.queries.InsertProject(ctx, dbgen.InsertProjectParams{
		ID:           projectID(now),
		Name:         input.Name,
		Company:      input.Company,
		StartYear:    input.StartYear,
		StartMonth:   input.StartMonth,
		EndYear:      toNullInt64(input.EndYear),
		EndMonth:     toNullInt64(input.EndMonth),
		Description:  input.Description,
		Role:         input.Role,
		TeamSize:     input.TeamSize,
		Technologies: technologies,
		Phases:       phases,
		Achievements: input.Achievements,
		IsDraft:      boolToInt64(input.IsDraft),
	})
	if err != nil {
		return nil, fmt.Errorf("insert project: %w", err)
	}

	created, err := r.queries.GetProject(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("reload project: %w", err)
	}

	return toDomainProject(created)
}

func (r *SQLiteProfileRepository) UpdateProject(ctx context.Context, id string, input domain.ProjectInput) (*domain.Project, error) {
	technologies, err := encodeStringSlice(input.Technologies)
	if err != nil {
		return nil, err
	}
	phases, err := encodeStringSlice(input.Phases)
	if err != nil {
		return nil, err
	}

	row, err := r.queries.UpdateProject(ctx, dbgen.UpdateProjectParams{
		ID:           id,
		Name:         input.Name,
		Company:      input.Company,
		StartYear:    input.StartYear,
		StartMonth:   input.StartMonth,
		EndYear:      toNullInt64(input.EndYear),
		EndMonth:     toNullInt64(input.EndMonth),
		Description:  input.Description,
		Role:         input.Role,
		TeamSize:     input.TeamSize,
		Technologies: technologies,
		Phases:       phases,
		Achievements: input.Achievements,
		IsDraft:      boolToInt64(input.IsDraft),
	})
	if err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}

	updated, err := r.queries.GetProject(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("reload project: %w", err)
	}

	return toDomainProject(updated)
}

func (r *SQLiteProfileRepository) DeleteProject(ctx context.Context, id string) error {
	rowsAffected, err := r.queries.DeleteProject(ctx, id)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SQLiteProfileRepository) ListJobHistoryOptions(ctx context.Context) (*domain.JobHistoryOptions, error) {
	rows, err := r.queries.ListJobEmploymentTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("query job employment types: %w", err)
	}

	return &domain.JobHistoryOptions{
		EmploymentTypes: toDomainJobEmploymentTypes(rows),
	}, nil
}

func (r *SQLiteProfileRepository) CreateJobEmploymentType(ctx context.Context, input domain.JobEmploymentTypeInput) (*domain.JobEmploymentType, error) {
	row, err := r.queries.InsertJobEmploymentType(ctx, dbgen.InsertJobEmploymentTypeParams{
		ID:   input.ID,
		Name: input.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("insert job employment type: %w", err)
	}

	return toDomainJobEmploymentType(row), nil
}

func (r *SQLiteProfileRepository) UpdateJobEmploymentType(ctx context.Context, id string, input domain.JobEmploymentTypeInput) (*domain.JobEmploymentType, error) {
	row, err := r.queries.UpdateJobEmploymentType(ctx, dbgen.UpdateJobEmploymentTypeParams{
		ID:   id,
		Name: input.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("update job employment type: %w", err)
	}

	return toDomainJobEmploymentType(row), nil
}

func ensureSQLiteDirectory(dsn string) error {
	if dsn == ":memory:" || strings.HasPrefix(dsn, "file:") {
		return nil
	}

	dir := filepath.Dir(dsn)
	if dir == "." {
		return nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create sqlite data directory: %w", err)
	}

	return nil
}

func (r *SQLiteProfileRepository) dropProfilePhoneColumn(ctx context.Context) error {
	hasPhoneColumn, err := r.profileTableHasColumn(ctx, "phone")
	if err != nil {
		return err
	}
	if !hasPhoneColumn {
		return nil
	}

	if _, err := r.db.ExecContext(ctx, "ALTER TABLE profiles DROP COLUMN phone"); err != nil {
		return fmt.Errorf("drop profiles.phone column: %w", err)
	}

	return nil
}

func (r *SQLiteProfileRepository) addProfileTextColumn(ctx context.Context, columnName string) error {
	hasColumn, err := r.profileTableHasColumn(ctx, columnName)
	if err != nil {
		return err
	}
	if hasColumn {
		return nil
	}

	if _, err := r.db.ExecContext(ctx, "ALTER TABLE profiles ADD COLUMN "+columnName+" TEXT NOT NULL DEFAULT ''"); err != nil {
		return fmt.Errorf("add profiles.%s column: %w", columnName, err)
	}

	return nil
}

func (r *SQLiteProfileRepository) profileTableHasColumn(ctx context.Context, columnName string) (bool, error) {
	return r.tableHasColumn(ctx, "profiles", columnName)
}

func (r *SQLiteProfileRepository) qualificationTableHasColumn(ctx context.Context, columnName string) (bool, error) {
	return r.tableHasColumn(ctx, "profile_qualifications", columnName)
}

func (r *SQLiteProfileRepository) addQualificationURLColumn(ctx context.Context) error {
	hasURLColumn, err := r.qualificationTableHasColumn(ctx, "url")
	if err != nil {
		return err
	}
	if hasURLColumn {
		return nil
	}

	if _, err := r.db.ExecContext(ctx, "ALTER TABLE profile_qualifications ADD COLUMN url TEXT NOT NULL DEFAULT ''"); err != nil {
		return fmt.Errorf("add profile_qualifications.url column: %w", err)
	}

	return nil
}

func (r *SQLiteProfileRepository) addSkillCategoryIconColumn(ctx context.Context) error {
	hasIconColumn, err := r.tableHasColumn(ctx, "skill_categories", "icon")
	if err != nil {
		return err
	}
	if hasIconColumn {
		return nil
	}

	if _, err := r.db.ExecContext(ctx, "ALTER TABLE skill_categories ADD COLUMN icon TEXT NOT NULL DEFAULT ''"); err != nil {
		return fmt.Errorf("add skill_categories.icon column: %w", err)
	}

	return nil
}

func (r *SQLiteProfileRepository) migrateJobHistoryDateColumns(ctx context.Context) error {
	hasStartDateColumn, err := r.tableHasColumn(ctx, "job_histories", "start_date")
	if err != nil {
		return err
	}
	hasEmploymentTypeColumn, err := r.tableHasColumn(ctx, "job_histories", "employment_type")
	if err != nil {
		return err
	}
	hasDisplayNameColumn, err := r.tableHasColumn(ctx, "job_histories", "display_name")
	if err != nil {
		return err
	}
	if !hasStartDateColumn && !hasEmploymentTypeColumn && hasDisplayNameColumn {
		return nil
	}
	if hasEmploymentTypeColumn {
		if _, err := r.db.ExecContext(
			ctx,
			`INSERT OR IGNORE INTO job_employment_types (id, name, sort_order)
			SELECT
				'job_employment_type_' || lower(hex(randomblob(8))),
				employment_type,
				COALESCE((SELECT MAX(sort_order) FROM job_employment_types), 0) + row_number() OVER (ORDER BY employment_type)
			FROM (SELECT DISTINCT employment_type FROM job_histories WHERE employment_type <> '') AS old_types
			WHERE NOT EXISTS (
				SELECT 1 FROM job_employment_types WHERE job_employment_types.name = old_types.employment_type
			)`,
		); err != nil {
			return fmt.Errorf("backfill job employment types: %w", err)
		}
	}

	startYearExpression := "start_year"
	startMonthExpression := "start_month"
	endYearExpression := "end_year"
	endMonthExpression := "end_month"
	if hasStartDateColumn {
		startYearExpression = "CAST(substr(start_date, 1, 4) AS INTEGER)"
		startMonthExpression = "CAST(substr(start_date, 6, 2) AS INTEGER)"
		endYearExpression = "CASE WHEN end_date = '' OR end_date = '現在' THEN NULL ELSE CAST(substr(end_date, 1, 4) AS INTEGER) END"
		endMonthExpression = "CASE WHEN end_date = '' OR end_date = '現在' THEN NULL ELSE CAST(substr(end_date, 6, 2) AS INTEGER) END"
	}

	employmentTypeIDExpression := "employment_type_id"
	if hasEmploymentTypeColumn {
		employmentTypeIDExpression = "(SELECT id FROM job_employment_types WHERE job_employment_types.name = job_histories.employment_type LIMIT 1)"
	}
	displayNameExpression := "display_name"
	if !hasDisplayNameColumn {
		displayNameExpression = "company"
	}

	statements := []string{
		`DROP TABLE IF EXISTS job_histories_new`,
		`CREATE TABLE IF NOT EXISTS job_histories_new (
			id TEXT NOT NULL PRIMARY KEY,
			company TEXT NOT NULL,
			display_name TEXT NOT NULL,
			start_year INTEGER NOT NULL,
			start_month INTEGER NOT NULL,
			end_year INTEGER,
			end_month INTEGER,
			employment_type_id TEXT NOT NULL,
			project_count INTEGER NOT NULL DEFAULT 0,
			sort_order INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (employment_type_id) REFERENCES job_employment_types(id)
		)`,
		fmt.Sprintf(`INSERT INTO job_histories_new (
			id,
			company,
			display_name,
			start_year,
			start_month,
			end_year,
			end_month,
			employment_type_id,
			project_count,
			sort_order,
			created_at,
			updated_at
		)
		SELECT
			id,
			company,
			%s,
			%s,
			%s,
			%s,
			%s,
			%s,
			project_count,
			sort_order,
			created_at,
			updated_at
		FROM job_histories`, displayNameExpression, startYearExpression, startMonthExpression, endYearExpression, endMonthExpression, employmentTypeIDExpression),
		`DROP TABLE job_histories`,
		`ALTER TABLE job_histories_new RENAME TO job_histories`,
	}
	for _, statement := range statements {
		if _, err := r.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("migrate job_histories date columns: %w", err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) migrateProjectDateColumns(ctx context.Context) error {
	hasStartDateColumn, err := r.tableHasColumn(ctx, "projects", "start_date")
	if err != nil {
		return err
	}
	if !hasStartDateColumn {
		return nil
	}

	statements := []string{
		`DROP TABLE IF EXISTS projects_new`,
		`CREATE TABLE IF NOT EXISTS projects_new (
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			company TEXT NOT NULL,
			start_year INTEGER NOT NULL,
			start_month INTEGER NOT NULL,
			end_year INTEGER,
			end_month INTEGER,
			description TEXT NOT NULL,
			role TEXT NOT NULL,
			team_size TEXT NOT NULL,
			technologies TEXT NOT NULL,
			phases TEXT NOT NULL,
			achievements TEXT NOT NULL,
			is_draft INTEGER NOT NULL DEFAULT 0,
			sort_order INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`INSERT INTO projects_new (
			id,
			name,
			company,
			start_year,
			start_month,
			end_year,
			end_month,
			description,
			role,
			team_size,
			technologies,
			phases,
			achievements,
			is_draft,
			sort_order,
			created_at,
			updated_at
		)
		SELECT
			id,
			name,
			company,
			CAST(substr(start_date, 1, 4) AS INTEGER),
			CAST(substr(start_date, 6, 2) AS INTEGER),
			CASE WHEN end_date = '' OR end_date = '現在' THEN NULL ELSE CAST(substr(end_date, 1, 4) AS INTEGER) END,
			CASE WHEN end_date = '' OR end_date = '現在' THEN NULL ELSE CAST(substr(end_date, 6, 2) AS INTEGER) END,
			description,
			role,
			team_size,
			technologies,
			phases,
			achievements,
			is_draft,
			sort_order,
			created_at,
			updated_at
		FROM projects`,
		`DROP TABLE projects`,
		`ALTER TABLE projects_new RENAME TO projects`,
	}
	for _, statement := range statements {
		if _, err := r.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("migrate projects date columns: %w", err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) seedSkillOptions(ctx context.Context) error {
	categories := []domain.SkillOption{
		{ID: "skill_category_language", Name: "言語", Icon: "code", SortOrder: 1},
		{ID: "skill_category_framework", Name: "フレームワーク", Icon: "code", SortOrder: 2},
		{ID: "skill_category_database", Name: "データベース", Icon: "database", SortOrder: 3},
		{ID: "skill_category_infrastructure", Name: "インフラ", Icon: "cloud", SortOrder: 4},
		{ID: "skill_category_tool", Name: "ツール", Icon: "wrench", SortOrder: 5},
		{ID: "skill_category_other", Name: "その他", Icon: "wrench", SortOrder: 6},
	}
	for _, category := range categories {
		if _, err := r.db.ExecContext(
			ctx,
			"INSERT OR IGNORE INTO skill_categories (id, name, icon, sort_order) VALUES (?, ?, ?, ?)",
			category.ID,
			category.Name,
			category.Icon,
			category.SortOrder,
		); err != nil {
			return fmt.Errorf("seed skill category %s: %w", category.ID, err)
		}
	}

	proficiencyLevels := []domain.SkillOption{
		{ID: "skill_proficiency_beginner", Name: "初級", SortOrder: 1},
		{ID: "skill_proficiency_intermediate", Name: "中級", SortOrder: 2},
		{ID: "skill_proficiency_advanced", Name: "上級", SortOrder: 3},
		{ID: "skill_proficiency_expert", Name: "エキスパート", SortOrder: 4},
	}
	for _, level := range proficiencyLevels {
		if _, err := r.db.ExecContext(
			ctx,
			"INSERT OR IGNORE INTO skill_proficiency_levels (id, name, sort_order) VALUES (?, ?, ?)",
			level.ID,
			level.Name,
			level.SortOrder,
		); err != nil {
			return fmt.Errorf("seed skill proficiency level %s: %w", level.ID, err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) seedSkills(ctx context.Context) error {
	skills := []struct {
		id                 string
		name               string
		categoryID         string
		experience         string
		proficiencyLevelID string
		sortOrder          int64
	}{
		{id: "skill_typescript", name: "TypeScript", categoryID: "skill_category_language", experience: "3年", proficiencyLevelID: "skill_proficiency_advanced", sortOrder: 1},
		{id: "skill_react", name: "React", categoryID: "skill_category_framework", experience: "3年", proficiencyLevelID: "skill_proficiency_advanced", sortOrder: 2},
		{id: "skill_nodejs", name: "Node.js", categoryID: "skill_category_framework", experience: "2年", proficiencyLevelID: "skill_proficiency_intermediate", sortOrder: 3},
		{id: "skill_postgresql", name: "PostgreSQL", categoryID: "skill_category_database", experience: "2年", proficiencyLevelID: "skill_proficiency_intermediate", sortOrder: 4},
		{id: "skill_docker", name: "Docker", categoryID: "skill_category_infrastructure", experience: "2年", proficiencyLevelID: "skill_proficiency_intermediate", sortOrder: 5},
		{id: "skill_aws", name: "AWS", categoryID: "skill_category_infrastructure", experience: "1年", proficiencyLevelID: "skill_proficiency_beginner", sortOrder: 6},
		{id: "skill_git", name: "Git", categoryID: "skill_category_tool", experience: "4年", proficiencyLevelID: "skill_proficiency_advanced", sortOrder: 7},
	}
	for _, skill := range skills {
		if _, err := r.db.ExecContext(
			ctx,
			"INSERT OR IGNORE INTO skills (id, name, category_id, experience, proficiency_level_id, sort_order) VALUES (?, ?, ?, ?, ?, ?)",
			skill.id,
			skill.name,
			skill.categoryID,
			skill.experience,
			skill.proficiencyLevelID,
			skill.sortOrder,
		); err != nil {
			return fmt.Errorf("seed skill %s: %w", skill.id, err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) seedJobEmploymentTypes(ctx context.Context) error {
	employmentTypes := []domain.JobEmploymentType{
		{ID: "job_employment_type_full_time", Name: "正社員", SortOrder: 1},
		{ID: "job_employment_type_contract", Name: "契約社員", SortOrder: 2},
		{ID: "job_employment_type_freelance", Name: "業務委託", SortOrder: 3},
		{ID: "job_employment_type_dispatch", Name: "派遣", SortOrder: 4},
		{ID: "job_employment_type_part_time", Name: "アルバイト", SortOrder: 5},
	}
	for _, employmentType := range employmentTypes {
		if _, err := r.db.ExecContext(
			ctx,
			"INSERT OR IGNORE INTO job_employment_types (id, name, sort_order) VALUES (?, ?, ?)",
			employmentType.ID,
			employmentType.Name,
			employmentType.SortOrder,
		); err != nil {
			return fmt.Errorf("seed job employment type %s: %w", employmentType.ID, err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) seedJobHistories(ctx context.Context) error {
	jobHistories := []struct {
		id               string
		company          string
		displayName      string
		startYear        int64
		startMonth       int64
		endYear          *int64
		endMonth         *int64
		employmentTypeID string
		projectCount     int64
		sortOrder        int64
	}{
		{id: "job_history_company_a", company: "株式会社A", displayName: "株式会社A", startYear: 2023, startMonth: 1, endYear: nil, endMonth: nil, employmentTypeID: "job_employment_type_full_time", projectCount: 5, sortOrder: 1},
		{id: "job_history_company_b", company: "株式会社B", displayName: "株式会社B", startYear: 2021, startMonth: 4, endYear: int64Ptr(2022), endMonth: int64Ptr(12), employmentTypeID: "job_employment_type_full_time", projectCount: 4, sortOrder: 2},
		{id: "job_history_freelance", company: "フリーランス", displayName: "フリーランス", startYear: 2020, startMonth: 1, endYear: int64Ptr(2021), endMonth: int64Ptr(3), employmentTypeID: "job_employment_type_freelance", projectCount: 3, sortOrder: 3},
	}
	for _, jobHistory := range jobHistories {
		if _, err := r.db.ExecContext(
			ctx,
			"INSERT OR IGNORE INTO job_histories (id, company, display_name, start_year, start_month, end_year, end_month, employment_type_id, project_count, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			jobHistory.id,
			jobHistory.company,
			jobHistory.displayName,
			jobHistory.startYear,
			jobHistory.startMonth,
			toNullInt64(jobHistory.endYear),
			toNullInt64(jobHistory.endMonth),
			jobHistory.employmentTypeID,
			jobHistory.projectCount,
			jobHistory.sortOrder,
		); err != nil {
			return fmt.Errorf("seed job history %s: %w", jobHistory.id, err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) seedProjects(ctx context.Context) error {
	projects := []domain.Project{
		{
			ID:           "project_ec_renewal",
			Name:         "ECサイトリニューアル",
			Company:      "株式会社A",
			StartYear:    2024,
			StartMonth:   1,
			EndYear:      nil,
			EndMonth:     nil,
			Description:  "大手ECサイトのフロントエンド刷新プロジェクト",
			Role:         "フロントエンドエンジニア",
			TeamSize:     "8名",
			Technologies: []string{"React", "TypeScript", "Next.js", "Tailwind CSS"},
			Phases:       []string{"要件定義", "設計", "実装", "テスト"},
			Achievements: "ページ表示速度を50%改善、コンバージョン率15%向上",
			IsDraft:      false,
			SortOrder:    1,
		},
		{
			ID:           "project_business_system",
			Name:         "業務管理システム開発",
			Company:      "株式会社B",
			StartYear:    2023,
			StartMonth:   6,
			EndYear:      int64Ptr(2023),
			EndMonth:     int64Ptr(12),
			Description:  "社内業務効率化のためのWebアプリケーション開発",
			Role:         "フルスタックエンジニア",
			TeamSize:     "5名",
			Technologies: []string{"Vue.js", "Node.js", "PostgreSQL", "Docker"},
			Phases:       []string{"設計", "実装", "テスト", "運用保守"},
			Achievements: "業務時間を30%削減、ユーザー満足度90%以上",
			IsDraft:      false,
			SortOrder:    2,
		},
	}

	for _, project := range projects {
		technologies, err := encodeStringSlice(project.Technologies)
		if err != nil {
			return err
		}
		phases, err := encodeStringSlice(project.Phases)
		if err != nil {
			return err
		}
		if _, err := r.db.ExecContext(
			ctx,
			`INSERT OR IGNORE INTO projects (
				id,
				name,
				company,
				start_year,
				start_month,
				end_year,
				end_month,
				description,
				role,
				team_size,
				technologies,
				phases,
				achievements,
				is_draft,
				sort_order
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			project.ID,
			project.Name,
			project.Company,
			project.StartYear,
			project.StartMonth,
			toNullInt64(project.EndYear),
			toNullInt64(project.EndMonth),
			project.Description,
			project.Role,
			project.TeamSize,
			technologies,
			phases,
			project.Achievements,
			boolToInt64(project.IsDraft),
			project.SortOrder,
		); err != nil {
			return fmt.Errorf("seed project %s: %w", project.ID, err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) backfillEmptySkillCategoryIcons(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, "UPDATE skill_categories SET icon = 'wrench', updated_at = CURRENT_TIMESTAMP WHERE icon = ''"); err != nil {
		return fmt.Errorf("backfill empty skill category icons: %w", err)
	}

	return nil
}

func (r *SQLiteProfileRepository) tableHasColumn(ctx context.Context, tableName string, columnName string) (bool, error) {
	rows, err := r.db.QueryContext(ctx, "PRAGMA table_info("+tableName+")")
	if err != nil {
		return false, fmt.Errorf("inspect %s columns: %w", tableName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid          int
			name         string
			columnType   string
			notNull      int
			defaultValue sql.NullString
			primaryKey   int
		)
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return false, fmt.Errorf("scan %s column: %w", tableName, err)
		}
		if name == columnName {
			return true, nil
		}
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("iterate %s columns: %w", tableName, err)
	}

	return false, nil
}

func toDomainProfile(row dbgen.Profile, qualificationRows []dbgen.ProfileQualification) (*domain.Profile, error) {
	profile := &domain.Profile{
		ID:                 row.ID,
		UserID:             row.UserID,
		FullName:           row.FullName,
		Nickname:           row.Nickname,
		Location:           row.Location,
		Email:              row.Email,
		Summary:            row.Summary,
		GitHubURL:          row.GithubUrl,
		ZennURL:            row.ZennUrl,
		QiitaURL:           row.QiitaUrl,
		WebsiteURL:         row.WebsiteUrl,
		Occupation:         row.Occupation,
		EmploymentType:     row.EmploymentType,
		PreferredWorkStyle: row.PreferredWorkStyle,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}

	if row.VisibilitySettings != "" {
		if err := json.Unmarshal([]byte(row.VisibilitySettings), &profile.VisibilitySettings); err != nil {
			return nil, fmt.Errorf("decode visibility settings: %w", err)
		}
	}

	profile.VisibilitySettings = sanitizeVisibilitySettings(profile.VisibilitySettings)
	profile.Qualifications = toDomainQualifications(qualificationRows)

	return profile, nil
}

func toDomainQualifications(rows []dbgen.ProfileQualification) []domain.Qualification {
	qualifications := make([]domain.Qualification, 0, len(rows))
	for _, row := range rows {
		qualifications = append(qualifications, domain.Qualification{
			ID:           row.ID,
			Name:         row.Name,
			AcquiredDate: row.AcquiredDate,
			Organization: row.Organization,
			URL:          row.Url,
			Memo:         row.Memo,
		})
	}

	return qualifications
}

func toDomainSkillCategoryOptions(rows []dbgen.SkillCategory) []domain.SkillOption {
	options := make([]domain.SkillOption, 0, len(rows))
	for _, row := range rows {
		options = append(options, *toDomainSkillCategoryOption(row))
	}

	return options
}

func toDomainSkillCategoryOption(row dbgen.SkillCategory) *domain.SkillOption {
	return &domain.SkillOption{
		ID:        row.ID,
		Name:      row.Name,
		Icon:      row.Icon,
		SortOrder: row.SortOrder,
	}
}

func toDomainSkillProficiencyLevelOptions(rows []dbgen.SkillProficiencyLevel) []domain.SkillOption {
	options := make([]domain.SkillOption, 0, len(rows))
	for _, row := range rows {
		options = append(options, domain.SkillOption{
			ID:        row.ID,
			Name:      row.Name,
			SortOrder: row.SortOrder,
		})
	}

	return options
}

func toDomainSkills(rows []dbgen.ListSkillsRow) []domain.Skill {
	skills := make([]domain.Skill, 0, len(rows))
	for _, row := range rows {
		skills = append(skills, domain.Skill{
			ID:                 row.ID,
			Name:               row.Name,
			CategoryID:         row.CategoryID,
			CategoryName:       row.CategoryName,
			Experience:         row.Experience,
			ProficiencyLevelID: row.ProficiencyLevelID,
			ProficiencyName:    row.ProficiencyName,
			SortOrder:          row.SortOrder,
		})
	}

	return skills
}

func toDomainSkill(row dbgen.GetSkillRow) *domain.Skill {
	return &domain.Skill{
		ID:                 row.ID,
		Name:               row.Name,
		CategoryID:         row.CategoryID,
		CategoryName:       row.CategoryName,
		Experience:         row.Experience,
		ProficiencyLevelID: row.ProficiencyLevelID,
		ProficiencyName:    row.ProficiencyName,
		SortOrder:          row.SortOrder,
	}
}

func toDomainJobHistories(rows []dbgen.ListJobHistoriesRow) []domain.JobHistory {
	jobHistories := make([]domain.JobHistory, 0, len(rows))
	for _, row := range rows {
		jobHistories = append(jobHistories, domain.JobHistory{
			ID:               row.ID,
			Company:          row.Company,
			DisplayName:      row.DisplayName,
			StartYear:        row.StartYear,
			StartMonth:       row.StartMonth,
			EndYear:          fromNullInt64(row.EndYear),
			EndMonth:         fromNullInt64(row.EndMonth),
			EmploymentTypeID: row.EmploymentTypeID,
			EmploymentType:   row.EmploymentType,
			ProjectCount:     row.ProjectCount,
			SortOrder:        row.SortOrder,
		})
	}

	return jobHistories
}

func toDomainJobHistory(row dbgen.GetJobHistoryRow) *domain.JobHistory {
	return &domain.JobHistory{
		ID:               row.ID,
		Company:          row.Company,
		DisplayName:      row.DisplayName,
		StartYear:        row.StartYear,
		StartMonth:       row.StartMonth,
		EndYear:          fromNullInt64(row.EndYear),
		EndMonth:         fromNullInt64(row.EndMonth),
		EmploymentTypeID: row.EmploymentTypeID,
		EmploymentType:   row.EmploymentType,
		ProjectCount:     row.ProjectCount,
		SortOrder:        row.SortOrder,
	}
}

func toDomainProjects(rows []dbgen.Project) ([]domain.Project, error) {
	projects := make([]domain.Project, 0, len(rows))
	for _, row := range rows {
		project, err := toDomainProject(row)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *project)
	}

	return projects, nil
}

func toDomainProject(row dbgen.Project) (*domain.Project, error) {
	technologies, err := decodeStringSlice(row.Technologies)
	if err != nil {
		return nil, fmt.Errorf("decode project technologies: %w", err)
	}
	phases, err := decodeStringSlice(row.Phases)
	if err != nil {
		return nil, fmt.Errorf("decode project phases: %w", err)
	}

	return &domain.Project{
		ID:           row.ID,
		Name:         row.Name,
		Company:      row.Company,
		StartYear:    row.StartYear,
		StartMonth:   row.StartMonth,
		EndYear:      fromNullInt64(row.EndYear),
		EndMonth:     fromNullInt64(row.EndMonth),
		Description:  row.Description,
		Role:         row.Role,
		TeamSize:     row.TeamSize,
		Technologies: technologies,
		Phases:       phases,
		Achievements: row.Achievements,
		IsDraft:      row.IsDraft != 0,
		SortOrder:    row.SortOrder,
	}, nil
}

func toDomainJobEmploymentTypes(rows []dbgen.JobEmploymentType) []domain.JobEmploymentType {
	employmentTypes := make([]domain.JobEmploymentType, 0, len(rows))
	for _, row := range rows {
		employmentTypes = append(employmentTypes, *toDomainJobEmploymentType(row))
	}

	return employmentTypes
}

func toDomainJobEmploymentType(row dbgen.JobEmploymentType) *domain.JobEmploymentType {
	return &domain.JobEmploymentType{
		ID:        row.ID,
		Name:      row.Name,
		SortOrder: row.SortOrder,
	}
}

func qualificationID(qualification domain.Qualification, index int, now time.Time) string {
	if qualification.ID != "" {
		return qualification.ID
	}

	return fmt.Sprintf("qualification_%d_%d", now.UnixNano(), index+1)
}

func skillID(now time.Time) string {
	return fmt.Sprintf("skill_%d", now.UnixNano())
}

func jobHistoryID(now time.Time) string {
	return fmt.Sprintf("job_history_%d", now.UnixNano())
}

func projectID(now time.Time) string {
	return fmt.Sprintf("project_%d", now.UnixNano())
}

func int64Ptr(value int64) *int64 {
	return &value
}

func encodeStringSlice(values []string) (string, error) {
	if values == nil {
		values = []string{}
	}

	bytes, err := json.Marshal(values)
	if err != nil {
		return "", fmt.Errorf("encode string slice: %w", err)
	}

	return string(bytes), nil
}

func decodeStringSlice(value string) ([]string, error) {
	if value == "" {
		return []string{}, nil
	}

	var values []string
	if err := json.Unmarshal([]byte(value), &values); err != nil {
		return nil, err
	}

	return values, nil
}

func boolToInt64(value bool) int64 {
	if value {
		return 1
	}

	return 0
}

func toNullInt64(value *int64) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{}
	}

	return sql.NullInt64{Int64: *value, Valid: true}
}

func fromNullInt64(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}

	return &value.Int64
}

func sanitizeVisibilitySettings(values map[string]any) map[string]any {
	result := map[string]any{}
	for key, value := range values {
		if key == "phone" {
			continue
		}
		result[key] = value
	}

	return result
}
