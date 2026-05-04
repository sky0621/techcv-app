package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

		links, err := r.listProfileLinks(ctx, row.ID)
		if err != nil {
			return nil, err
		}

		return toDomainProfile(row, qualifications, links)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("query profile: %w", err)
	}

	now := time.Now().UTC()
	profile := domain.Profile{
		ID:                 "1",
		UserID:             "user_01",
		VisibilitySettings: map[string]any{"email": true, "location": true},
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	return r.Save(ctx, &profile)
}

func (r *SQLiteProfileRepository) Save(ctx context.Context, profile *domain.Profile) (*domain.Profile, error) {
	if _, err := strconv.ParseInt(profile.ID, 10, 64); err != nil {
		profile.ID = "1"
	}
	syncLegacyProfileLinkColumns(profile)

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

	if err := r.saveProfileLinks(ctx, tx, profile.ID, profile.Links); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit profile transaction: %w", err)
	}

	queries = r.queries

	row, err := queries.GetProfileByUserID(ctx, profile.UserID)
	if err != nil {
		return nil, fmt.Errorf("reload profile: %w", err)
	}

	qualifications, err := queries.ListQualificationsByProfileID(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("reload profile qualifications: %w", err)
	}

	links, err := r.listProfileLinks(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	result, err := toDomainProfile(row, qualifications, links)
	if err != nil {
		return nil, err
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
	if err := r.seedProfileLinkMasters(ctx); err != nil {
		return err
	}
	if err := r.migrateProfileLinks(ctx); err != nil {
		return err
	}
	if err := r.seedSkillOptions(ctx); err != nil {
		return err
	}
	if err := r.seedSkillMasters(ctx); err != nil {
		return err
	}
	if err := r.migrateSkillExperienceColumn(ctx); err != nil {
		return err
	}
	if err := r.migrateSkillMasterReference(ctx); err != nil {
		return err
	}
	if err := r.addQualificationURLColumn(ctx); err != nil {
		return err
	}
	if err := r.migrateIntegerPrimaryKeys(ctx); err != nil {
		return err
	}
	if err := r.seedSkills(ctx); err != nil {
		return err
	}
	if err := r.seedJobEmploymentTypes(ctx); err != nil {
		return err
	}
	if err := r.seedJobCompanies(ctx); err != nil {
		return err
	}
	if err := r.migrateJobHistoryDateColumns(ctx); err != nil {
		return err
	}
	if err := r.migrateJobHistoryNullableNames(ctx); err != nil {
		return err
	}
	if err := r.migrateJobHistorySortOrderColumn(ctx); err != nil {
		return err
	}
	if err := r.backfillJobCompanies(ctx); err != nil {
		return err
	}
	if err := r.migrateJobHistoryCompanyID(ctx); err != nil {
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

	return nil
}

func (r *SQLiteProfileRepository) ListProfileLinkMasters(ctx context.Context) ([]domain.ProfileLinkMaster, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, key, name, icon, placeholder, sort_order
		FROM profile_link_masters
		ORDER BY sort_order ASC, name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query profile link masters: %w", err)
	}
	defer rows.Close()

	masters := []domain.ProfileLinkMaster{}
	for rows.Next() {
		var master domain.ProfileLinkMaster
		if err := rows.Scan(&master.ID, &master.Key, &master.Name, &master.Icon, &master.Placeholder, &master.SortOrder); err != nil {
			return nil, fmt.Errorf("scan profile link master: %w", err)
		}
		masters = append(masters, master)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate profile link masters: %w", err)
	}

	return masters, nil
}

func (r *SQLiteProfileRepository) CreateProfileLinkMaster(ctx context.Context, input domain.ProfileLinkMasterInput) (*domain.ProfileLinkMaster, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO profile_link_masters (id, key, name, icon, placeholder, sort_order)
		SELECT ?, ?, ?, ?, ?, COALESCE(NULLIF(?, 0), COALESCE(MAX(sort_order), 0) + 1)
		FROM profile_link_masters
		RETURNING id, key, name, icon, placeholder, sort_order
	`, input.ID, input.Key, input.Name, input.Icon, input.Placeholder, input.SortOrder)

	var master domain.ProfileLinkMaster
	if err := row.Scan(&master.ID, &master.Key, &master.Name, &master.Icon, &master.Placeholder, &master.SortOrder); err != nil {
		return nil, fmt.Errorf("insert profile link master: %w", err)
	}

	return &master, nil
}

func (r *SQLiteProfileRepository) UpdateProfileLinkMaster(ctx context.Context, id string, input domain.ProfileLinkMasterInput) (*domain.ProfileLinkMaster, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE profile_link_masters
		SET key = ?, name = ?, icon = ?, placeholder = ?, sort_order = COALESCE(NULLIF(?, 0), sort_order), updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING id, key, name, icon, placeholder, sort_order
	`, input.Key, input.Name, input.Icon, input.Placeholder, input.SortOrder, id)

	var master domain.ProfileLinkMaster
	if err := row.Scan(&master.ID, &master.Key, &master.Name, &master.Icon, &master.Placeholder, &master.SortOrder); err != nil {
		return nil, fmt.Errorf("update profile link master: %w", err)
	}

	return &master, nil
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
	skillMasters, err := r.queries.ListSkillMasters(ctx)
	if err != nil {
		return nil, fmt.Errorf("query skill masters: %w", err)
	}

	return &domain.SkillOptions{
		Categories:        toDomainSkillCategoryOptions(categories),
		ProficiencyLevels: toDomainSkillProficiencyLevelOptions(proficiencyLevels),
		SkillMasters:      toDomainSkillMasters(skillMasters),
	}, nil
}

func (r *SQLiteProfileRepository) CreateSkillCategory(ctx context.Context, input domain.SkillCategoryInput) (*domain.SkillOption, error) {
	row, err := r.queries.InsertSkillCategory(ctx, dbgen.InsertSkillCategoryParams{
		ID:        input.ID,
		Name:      input.Name,
		Icon:      input.Icon,
		SortOrder: input.SortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("insert skill category: %w", err)
	}

	return toDomainSkillCategoryOption(row), nil
}

func (r *SQLiteProfileRepository) UpdateSkillCategory(ctx context.Context, id string, input domain.SkillCategoryInput) (*domain.SkillOption, error) {
	row, err := r.queries.UpdateSkillCategory(ctx, dbgen.UpdateSkillCategoryParams{
		ID:        id,
		Name:      input.Name,
		Icon:      input.Icon,
		SortOrder: input.SortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("update skill category: %w", err)
	}

	return toDomainSkillCategoryOption(row), nil
}

func (r *SQLiteProfileRepository) UpdateSkillProficiencyLevel(ctx context.Context, id string, input domain.SkillProficiencyLevelInput) (*domain.SkillOption, error) {
	row, err := r.queries.UpdateSkillProficiencyLevel(ctx, dbgen.UpdateSkillProficiencyLevelParams{
		ID:        id,
		Name:      input.Name,
		SortOrder: input.SortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("update skill proficiency level: %w", err)
	}

	return toDomainSkillProficiencyLevelOption(row), nil
}

func (r *SQLiteProfileRepository) CreateSkillMaster(ctx context.Context, input domain.SkillMasterInput) (*domain.SkillMaster, error) {
	row, err := r.queries.InsertSkillMaster(ctx, dbgen.InsertSkillMasterParams{
		ID:         input.ID,
		Name:       input.Name,
		CategoryID: input.CategoryID,
		SortOrder:  input.SortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("insert skill master: %w", err)
	}

	created, err := r.queries.GetSkillMaster(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("reload skill master: %w", err)
	}

	return toDomainSkillMaster(created), nil
}

func (r *SQLiteProfileRepository) UpdateSkillMaster(ctx context.Context, id string, input domain.SkillMasterInput) (*domain.SkillMaster, error) {
	row, err := r.queries.UpdateSkillMaster(ctx, dbgen.UpdateSkillMasterParams{
		ID:         id,
		Name:       input.Name,
		CategoryID: input.CategoryID,
		SortOrder:  input.SortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("update skill master: %w", err)
	}

	updated, err := r.queries.GetSkillMaster(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("reload skill master: %w", err)
	}

	return toDomainSkillMaster(updated), nil
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
		SkillMasterID:      input.SkillMasterID,
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
		SkillMasterID:      input.SkillMasterID,
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
	row, err := r.queries.InsertJobHistory(ctx, dbgen.InsertJobHistoryParams{
		CompanyID:        input.CompanyID,
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
		CompanyID:        input.CompanyID,
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
	companies, err := r.listJobCompanies(ctx)
	if err != nil {
		return nil, err
	}

	return &domain.JobHistoryOptions{
		EmploymentTypes: toDomainJobEmploymentTypes(rows),
		Companies:       companies,
	}, nil
}

func (r *SQLiteProfileRepository) CreateJobEmploymentType(ctx context.Context, input domain.JobEmploymentTypeInput) (*domain.JobEmploymentType, error) {
	row, err := r.queries.InsertJobEmploymentType(ctx, dbgen.InsertJobEmploymentTypeParams{
		ID:        input.ID,
		Name:      input.Name,
		SortOrder: input.SortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("insert job employment type: %w", err)
	}

	return toDomainJobEmploymentType(row), nil
}

func (r *SQLiteProfileRepository) UpdateJobEmploymentType(ctx context.Context, id string, input domain.JobEmploymentTypeInput) (*domain.JobEmploymentType, error) {
	row, err := r.queries.UpdateJobEmploymentType(ctx, dbgen.UpdateJobEmploymentTypeParams{
		ID:        id,
		Name:      input.Name,
		SortOrder: input.SortOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("update job employment type: %w", err)
	}

	return toDomainJobEmploymentType(row), nil
}

func (r *SQLiteProfileRepository) CreateJobCompany(ctx context.Context, input domain.JobCompanyInput) (*domain.JobCompany, error) {
	if input.ID == "" {
		return r.insertJobCompany(ctx, "INSERT INTO job_companies (name, url) VALUES (?, ?) RETURNING id, name, url", input.Name, input.URL)
	}

	return r.insertJobCompany(ctx, "INSERT INTO job_companies (id, name, url) VALUES (?, ?, ?) RETURNING id, name, url", input.ID, input.Name, input.URL)
}

func (r *SQLiteProfileRepository) UpdateJobCompany(ctx context.Context, id string, input domain.JobCompanyInput) (*domain.JobCompany, error) {
	return r.insertJobCompany(ctx, `
		UPDATE job_companies
		SET name = ?, url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING id, name, url
	`, input.Name, input.URL, id)
}

func (r *SQLiteProfileRepository) insertJobCompany(ctx context.Context, query string, args ...any) (*domain.JobCompany, error) {
	var company domain.JobCompany
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&company.ID, &company.Name, &company.URL); err != nil {
		return nil, fmt.Errorf("save job company: %w", err)
	}

	return &company, nil
}

func (r *SQLiteProfileRepository) listJobCompanies(ctx context.Context) ([]domain.JobCompany, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, url
		FROM job_companies
		ORDER BY name ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query job companies: %w", err)
	}
	defer rows.Close()

	companies := make([]domain.JobCompany, 0)
	for rows.Next() {
		var company domain.JobCompany
		if err := rows.Scan(&company.ID, &company.Name, &company.URL); err != nil {
			return nil, fmt.Errorf("scan job company: %w", err)
		}
		companies = append(companies, company)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate job companies: %w", err)
	}

	return companies, nil
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
			company TEXT,
			display_name TEXT,
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

func (r *SQLiteProfileRepository) migrateJobHistoryNullableNames(ctx context.Context) error {
	hasCompanyColumn, err := r.tableHasColumn(ctx, "job_histories", "company")
	if err != nil {
		return err
	}
	if !hasCompanyColumn {
		return nil
	}

	companyNotNull, err := r.tableColumnNotNull(ctx, "job_histories", "company")
	if err != nil {
		return err
	}
	displayNameNotNull, err := r.tableColumnNotNull(ctx, "job_histories", "display_name")
	if err != nil {
		return err
	}
	if !companyNotNull && !displayNameNotNull {
		return nil
	}

	statements := []string{
		`PRAGMA foreign_keys = OFF`,
		`DROP TABLE IF EXISTS job_histories_new`,
		`CREATE TABLE job_histories_new (
			id INTEGER PRIMARY KEY,
			company TEXT,
			display_name TEXT,
			start_year INTEGER NOT NULL,
			start_month INTEGER NOT NULL,
			end_year INTEGER,
			end_month INTEGER,
			employment_type_id INTEGER NOT NULL,
			project_count INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (employment_type_id) REFERENCES job_employment_types(id)
		)`,
		`INSERT INTO job_histories_new
		SELECT
			id,
			NULLIF(company, ''),
			NULLIF(display_name, ''),
			start_year,
			start_month,
			end_year,
			end_month,
			employment_type_id,
			project_count,
			sort_order,
			created_at,
			updated_at
		FROM job_histories`,
		`DROP TABLE job_histories`,
		`ALTER TABLE job_histories_new RENAME TO job_histories`,
		`PRAGMA foreign_keys = ON`,
	}
	for _, statement := range statements {
		if _, err := r.db.ExecContext(ctx, statement); err != nil {
			_, _ = r.db.ExecContext(ctx, "PRAGMA foreign_keys = ON")
			return fmt.Errorf("migrate job_histories nullable names: %w", err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) migrateJobHistorySortOrderColumn(ctx context.Context) error {
	hasSortOrderColumn, err := r.tableHasColumn(ctx, "job_histories", "sort_order")
	if err != nil {
		return err
	}
	if !hasSortOrderColumn {
		return nil
	}

	statements := []string{
		`PRAGMA foreign_keys = OFF`,
		`DROP TABLE IF EXISTS job_histories_new`,
		`CREATE TABLE job_histories_new (
			id INTEGER PRIMARY KEY,
			company TEXT,
			display_name TEXT,
			start_year INTEGER NOT NULL,
			start_month INTEGER NOT NULL,
			end_year INTEGER,
			end_month INTEGER,
			employment_type_id INTEGER NOT NULL,
			project_count INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (employment_type_id) REFERENCES job_employment_types(id)
		)`,
		`INSERT INTO job_histories_new
		SELECT
			id,
			company,
			display_name,
			start_year,
			start_month,
			end_year,
			end_month,
			employment_type_id,
			project_count,
			created_at,
			updated_at
		FROM job_histories`,
		`DROP TABLE job_histories`,
		`ALTER TABLE job_histories_new RENAME TO job_histories`,
		`PRAGMA foreign_keys = ON`,
	}
	for _, statement := range statements {
		if _, err := r.db.ExecContext(ctx, statement); err != nil {
			_, _ = r.db.ExecContext(ctx, "PRAGMA foreign_keys = ON")
			return fmt.Errorf("migrate job_histories sort_order column: %w", err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) migrateJobHistoryCompanyID(ctx context.Context) error {
	hasCompanyColumn, err := r.tableHasColumn(ctx, "job_histories", "company")
	if err != nil {
		return err
	}
	if !hasCompanyColumn {
		return nil
	}

	statements := []string{
		`PRAGMA foreign_keys = OFF`,
		`DROP TABLE IF EXISTS job_histories_new`,
		`CREATE TABLE job_histories_new (
			id INTEGER PRIMARY KEY,
			company_id INTEGER,
			display_name TEXT,
			start_year INTEGER NOT NULL,
			start_month INTEGER NOT NULL,
			end_year INTEGER,
			end_month INTEGER,
			employment_type_id INTEGER NOT NULL,
			project_count INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (company_id) REFERENCES job_companies(id),
			FOREIGN KEY (employment_type_id) REFERENCES job_employment_types(id)
		)`,
		`INSERT INTO job_histories_new (
			id,
			company_id,
			display_name,
			start_year,
			start_month,
			end_year,
			end_month,
			employment_type_id,
			project_count,
			created_at,
			updated_at
		)
		SELECT
			job_histories.id,
			job_companies.id,
			job_histories.display_name,
			job_histories.start_year,
			job_histories.start_month,
			job_histories.end_year,
			job_histories.end_month,
			job_histories.employment_type_id,
			job_histories.project_count,
			job_histories.created_at,
			job_histories.updated_at
		FROM job_histories
		LEFT JOIN job_companies ON job_companies.name = job_histories.company`,
		`DROP TABLE job_histories`,
		`ALTER TABLE job_histories_new RENAME TO job_histories`,
		`PRAGMA foreign_keys = ON`,
	}
	for _, statement := range statements {
		if _, err := r.db.ExecContext(ctx, statement); err != nil {
			_, _ = r.db.ExecContext(ctx, "PRAGMA foreign_keys = ON")
			return fmt.Errorf("migrate job_histories company_id: %w", err)
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

func (r *SQLiteProfileRepository) migrateSkillExperienceColumn(ctx context.Context) error {
	columnType, err := r.tableColumnType(ctx, "skills", "experience")
	if err != nil {
		return err
	}
	if strings.EqualFold(columnType, "INTEGER") {
		return nil
	}

	statements := []string{
		`DROP TABLE IF EXISTS skills_new`,
		`CREATE TABLE IF NOT EXISTS skills_new (
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			category_id TEXT NOT NULL,
			experience INTEGER NOT NULL,
			proficiency_level_id TEXT NOT NULL,
			sort_order INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (category_id) REFERENCES skill_categories(id),
			FOREIGN KEY (proficiency_level_id) REFERENCES skill_proficiency_levels(id)
		)`,
		`INSERT INTO skills_new (
			id,
			name,
			category_id,
			experience,
			proficiency_level_id,
			sort_order,
			created_at,
			updated_at
		)
		SELECT
			id,
			name,
			category_id,
			CAST(REPLACE(experience, '年', '') AS INTEGER),
			proficiency_level_id,
			sort_order,
			created_at,
			updated_at
		FROM skills`,
		`DROP TABLE skills`,
		`ALTER TABLE skills_new RENAME TO skills`,
	}
	for _, statement := range statements {
		if _, err := r.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("migrate skills experience column: %w", err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) migrateSkillMasterReference(ctx context.Context) error {
	hasSkillMasterIDColumn, err := r.tableHasColumn(ctx, "skills", "skill_master_id")
	if err != nil {
		return err
	}
	if hasSkillMasterIDColumn {
		if _, err := r.db.ExecContext(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS idx_skills_skill_master_id_unique ON skills(skill_master_id)`); err != nil {
			return fmt.Errorf("ensure unique skill master ids: %w", err)
		}

		return nil
	}

	statements := []string{
		`DROP TABLE IF EXISTS skills_new`,
		`CREATE TABLE IF NOT EXISTS skills_new (
			id TEXT NOT NULL PRIMARY KEY,
			skill_master_id TEXT NOT NULL UNIQUE,
			experience INTEGER NOT NULL,
			proficiency_level_id TEXT NOT NULL,
			sort_order INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (skill_master_id) REFERENCES skill_masters(id),
			FOREIGN KEY (proficiency_level_id) REFERENCES skill_proficiency_levels(id)
		)`,
		`INSERT OR IGNORE INTO skills_new (
			id,
			skill_master_id,
			experience,
			proficiency_level_id,
			sort_order,
			created_at,
			updated_at
		)
		SELECT
			skills.id,
			skill_masters.id,
			skills.experience,
			skills.proficiency_level_id,
			skills.sort_order,
			skills.created_at,
			skills.updated_at
		FROM skills
		JOIN skill_masters ON skill_masters.name = skills.name
		ORDER BY skills.sort_order ASC, skills.name ASC`,
		`DROP TABLE skills`,
		`ALTER TABLE skills_new RENAME TO skills`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_skills_skill_master_id_unique ON skills(skill_master_id)`,
	}
	for _, statement := range statements {
		if _, err := r.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("migrate skills skill_master_id reference: %w", err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) migrateIntegerPrimaryKeys(ctx context.Context) error {
	tables := []string{
		"profiles",
		"profile_qualifications",
		"skill_categories",
		"skill_proficiency_levels",
		"skill_masters",
		"skills",
		"job_employment_types",
		"job_histories",
		"projects",
	}
	needsMigration := false
	for _, table := range tables {
		idType, err := r.tableColumnType(ctx, table, "id")
		if err != nil {
			return err
		}
		if !strings.EqualFold(idType, "INTEGER") {
			needsMigration = true
			break
		}
	}
	if !needsMigration {
		return nil
	}

	jobHistoryIDMapStatement := `CREATE TEMP TABLE _job_history_id_map AS
		SELECT id AS old_id, ROW_NUMBER() OVER (ORDER BY start_year DESC, start_month DESC, id ASC) AS new_id
		FROM job_histories`
	hasJobHistorySortOrderColumn, err := r.tableHasColumn(ctx, "job_histories", "sort_order")
	if err != nil {
		return err
	}
	if hasJobHistorySortOrderColumn {
		jobHistoryIDMapStatement = `CREATE TEMP TABLE _job_history_id_map AS
			SELECT id AS old_id, ROW_NUMBER() OVER (ORDER BY sort_order ASC, id ASC) AS new_id
			FROM job_histories`
	}
	jobHistoryCompanyIDExpression := "NULL"
	hasJobHistoryCompanyIDColumn, err := r.tableHasColumn(ctx, "job_histories", "company_id")
	if err != nil {
		return err
	}
	if hasJobHistoryCompanyIDColumn {
		jobHistoryCompanyIDExpression = "job_histories.company_id"
	}

	statements := []string{
		`PRAGMA foreign_keys = OFF`,
		`DROP TABLE IF EXISTS profile_qualifications_new`,
		`DROP TABLE IF EXISTS profiles_new`,
		`DROP TABLE IF EXISTS skills_new`,
		`DROP TABLE IF EXISTS skill_masters_new`,
		`DROP TABLE IF EXISTS skill_categories_new`,
		`DROP TABLE IF EXISTS skill_proficiency_levels_new`,
		`DROP TABLE IF EXISTS job_histories_new`,
		`DROP TABLE IF EXISTS job_employment_types_new`,
		`DROP TABLE IF EXISTS projects_new`,
		`DROP TABLE IF EXISTS _profile_id_map`,
		`DROP TABLE IF EXISTS _qualification_id_map`,
		`DROP TABLE IF EXISTS _skill_category_id_map`,
		`DROP TABLE IF EXISTS _skill_proficiency_id_map`,
		`DROP TABLE IF EXISTS _skill_master_id_map`,
		`DROP TABLE IF EXISTS _skill_id_map`,
		`DROP TABLE IF EXISTS _job_employment_type_id_map`,
		`DROP TABLE IF EXISTS _job_history_id_map`,
		`DROP TABLE IF EXISTS _project_id_map`,
		`CREATE TEMP TABLE _profile_id_map AS
		SELECT id AS old_id, ROW_NUMBER() OVER (ORDER BY created_at ASC, id ASC) AS new_id
		FROM profiles`,
		`CREATE TEMP TABLE _qualification_id_map AS
		SELECT id AS old_id, ROW_NUMBER() OVER (ORDER BY sort_order ASC, created_at ASC, id ASC) AS new_id
		FROM profile_qualifications`,
		`CREATE TEMP TABLE _skill_category_id_map AS
		SELECT id AS old_id, ROW_NUMBER() OVER (ORDER BY sort_order ASC, name ASC, id ASC) AS new_id
		FROM skill_categories`,
		`CREATE TEMP TABLE _skill_proficiency_id_map AS
		SELECT id AS old_id, ROW_NUMBER() OVER (ORDER BY sort_order ASC, name ASC, id ASC) AS new_id
		FROM skill_proficiency_levels`,
		`CREATE TEMP TABLE _skill_master_id_map AS
		SELECT id AS old_id, ROW_NUMBER() OVER (ORDER BY sort_order ASC, name ASC, id ASC) AS new_id
		FROM skill_masters`,
		`CREATE TEMP TABLE _skill_id_map AS
		SELECT id AS old_id, ROW_NUMBER() OVER (ORDER BY sort_order ASC, id ASC) AS new_id
		FROM skills`,
		`CREATE TEMP TABLE _job_employment_type_id_map AS
		SELECT id AS old_id, ROW_NUMBER() OVER (ORDER BY sort_order ASC, name ASC, id ASC) AS new_id
		FROM job_employment_types`,
		jobHistoryIDMapStatement,
		`CREATE TEMP TABLE _project_id_map AS
		SELECT id AS old_id, ROW_NUMBER() OVER (ORDER BY sort_order ASC, id ASC) AS new_id
		FROM projects`,
		`CREATE TABLE profiles_new (
			id INTEGER PRIMARY KEY,
			user_id TEXT NOT NULL UNIQUE,
			full_name TEXT NOT NULL,
			nickname TEXT NOT NULL,
			location TEXT NOT NULL,
			email TEXT NOT NULL,
			summary TEXT NOT NULL,
			github_url TEXT NOT NULL,
			zenn_url TEXT NOT NULL,
			qiita_url TEXT NOT NULL,
			website_url TEXT NOT NULL,
			occupation TEXT NOT NULL,
			employment_type TEXT NOT NULL,
			preferred_work_style TEXT NOT NULL,
			visibility_settings TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`INSERT INTO profiles_new
		SELECT
			_profile_id_map.new_id,
			profiles.user_id,
			profiles.full_name,
			profiles.nickname,
			profiles.location,
			profiles.email,
			profiles.summary,
			profiles.github_url,
			profiles.zenn_url,
			profiles.qiita_url,
			profiles.website_url,
			profiles.occupation,
			profiles.employment_type,
			profiles.preferred_work_style,
			profiles.visibility_settings,
			profiles.created_at,
			profiles.updated_at
		FROM profiles
		JOIN _profile_id_map ON _profile_id_map.old_id = profiles.id`,
		`CREATE TABLE profile_qualifications_new (
			id INTEGER PRIMARY KEY,
			profile_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			acquired_date TEXT NOT NULL,
			organization TEXT NOT NULL,
			url TEXT NOT NULL,
			memo TEXT NOT NULL,
			sort_order INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
		)`,
		`INSERT INTO profile_qualifications_new
		SELECT
			_qualification_id_map.new_id,
			_profile_id_map.new_id,
			profile_qualifications.name,
			profile_qualifications.acquired_date,
			profile_qualifications.organization,
			profile_qualifications.url,
			profile_qualifications.memo,
			profile_qualifications.sort_order,
			profile_qualifications.created_at,
			profile_qualifications.updated_at
		FROM profile_qualifications
		JOIN _qualification_id_map ON _qualification_id_map.old_id = profile_qualifications.id
		JOIN _profile_id_map ON _profile_id_map.old_id = profile_qualifications.profile_id`,
		`CREATE TABLE skill_categories_new (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			icon TEXT NOT NULL,
			sort_order INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`INSERT INTO skill_categories_new
		SELECT _skill_category_id_map.new_id, skill_categories.name, skill_categories.icon, skill_categories.sort_order, skill_categories.created_at, skill_categories.updated_at
		FROM skill_categories
		JOIN _skill_category_id_map ON _skill_category_id_map.old_id = skill_categories.id`,
		`CREATE TABLE skill_proficiency_levels_new (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			sort_order INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`INSERT INTO skill_proficiency_levels_new
		SELECT _skill_proficiency_id_map.new_id, skill_proficiency_levels.name, skill_proficiency_levels.sort_order, skill_proficiency_levels.created_at, skill_proficiency_levels.updated_at
		FROM skill_proficiency_levels
		JOIN _skill_proficiency_id_map ON _skill_proficiency_id_map.old_id = skill_proficiency_levels.id`,
		`CREATE TABLE skill_masters_new (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			category_id INTEGER NOT NULL,
			sort_order INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (category_id) REFERENCES skill_categories(id)
		)`,
		`INSERT INTO skill_masters_new
		SELECT
			_skill_master_id_map.new_id,
			skill_masters.name,
			_skill_category_id_map.new_id,
			skill_masters.sort_order,
			skill_masters.created_at,
			skill_masters.updated_at
		FROM skill_masters
		JOIN _skill_master_id_map ON _skill_master_id_map.old_id = skill_masters.id
		JOIN _skill_category_id_map ON _skill_category_id_map.old_id = skill_masters.category_id`,
		`CREATE TABLE skills_new (
			id INTEGER PRIMARY KEY,
			skill_master_id INTEGER NOT NULL UNIQUE,
			experience INTEGER NOT NULL,
			proficiency_level_id INTEGER NOT NULL,
			sort_order INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (skill_master_id) REFERENCES skill_masters(id),
			FOREIGN KEY (proficiency_level_id) REFERENCES skill_proficiency_levels(id)
		)`,
		`INSERT OR IGNORE INTO skills_new
		SELECT
			_skill_id_map.new_id,
			_skill_master_id_map.new_id,
			skills.experience,
			_skill_proficiency_id_map.new_id,
			skills.sort_order,
			skills.created_at,
			skills.updated_at
		FROM skills
		JOIN _skill_id_map ON _skill_id_map.old_id = skills.id
		JOIN _skill_master_id_map ON _skill_master_id_map.old_id = skills.skill_master_id
		JOIN _skill_proficiency_id_map ON _skill_proficiency_id_map.old_id = skills.proficiency_level_id`,
		`CREATE TABLE job_employment_types_new (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			sort_order INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`INSERT INTO job_employment_types_new
		SELECT _job_employment_type_id_map.new_id, job_employment_types.name, job_employment_types.sort_order, job_employment_types.created_at, job_employment_types.updated_at
		FROM job_employment_types
		JOIN _job_employment_type_id_map ON _job_employment_type_id_map.old_id = job_employment_types.id`,
		`CREATE TABLE job_histories_new (
			id INTEGER PRIMARY KEY,
			company_id INTEGER,
			display_name TEXT,
			start_year INTEGER NOT NULL,
			start_month INTEGER NOT NULL,
			end_year INTEGER,
			end_month INTEGER,
			employment_type_id INTEGER NOT NULL,
			project_count INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (company_id) REFERENCES job_companies(id),
			FOREIGN KEY (employment_type_id) REFERENCES job_employment_types(id)
		)`,
		fmt.Sprintf(`INSERT INTO job_histories_new
		SELECT
			_job_history_id_map.new_id,
			%s,
			job_histories.display_name,
			job_histories.start_year,
			job_histories.start_month,
			job_histories.end_year,
			job_histories.end_month,
			_job_employment_type_id_map.new_id,
			job_histories.project_count,
			job_histories.created_at,
			job_histories.updated_at
		FROM job_histories
		JOIN _job_history_id_map ON _job_history_id_map.old_id = job_histories.id
		JOIN _job_employment_type_id_map ON _job_employment_type_id_map.old_id = job_histories.employment_type_id`, jobHistoryCompanyIDExpression),
		`CREATE TABLE projects_new (
			id INTEGER PRIMARY KEY,
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
		`INSERT INTO projects_new
		SELECT
			_project_id_map.new_id,
			projects.name,
			projects.company,
			projects.start_year,
			projects.start_month,
			projects.end_year,
			projects.end_month,
			projects.description,
			projects.role,
			projects.team_size,
			projects.technologies,
			projects.phases,
			projects.achievements,
			projects.is_draft,
			projects.sort_order,
			projects.created_at,
			projects.updated_at
		FROM projects
		JOIN _project_id_map ON _project_id_map.old_id = projects.id`,
		`DROP TABLE profile_qualifications`,
		`DROP TABLE profiles`,
		`DROP TABLE skills`,
		`DROP TABLE skill_masters`,
		`DROP TABLE skill_categories`,
		`DROP TABLE skill_proficiency_levels`,
		`DROP TABLE job_histories`,
		`DROP TABLE job_employment_types`,
		`DROP TABLE projects`,
		`ALTER TABLE profiles_new RENAME TO profiles`,
		`ALTER TABLE profile_qualifications_new RENAME TO profile_qualifications`,
		`ALTER TABLE skill_categories_new RENAME TO skill_categories`,
		`ALTER TABLE skill_proficiency_levels_new RENAME TO skill_proficiency_levels`,
		`ALTER TABLE skill_masters_new RENAME TO skill_masters`,
		`ALTER TABLE skills_new RENAME TO skills`,
		`ALTER TABLE job_employment_types_new RENAME TO job_employment_types`,
		`ALTER TABLE job_histories_new RENAME TO job_histories`,
		`ALTER TABLE projects_new RENAME TO projects`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_skills_skill_master_id_unique ON skills(skill_master_id)`,
		`PRAGMA foreign_keys = ON`,
	}
	for _, statement := range statements {
		if _, err := r.db.ExecContext(ctx, statement); err != nil {
			_, _ = r.db.ExecContext(ctx, "PRAGMA foreign_keys = ON")
			return fmt.Errorf("migrate integer primary keys: %w", err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) seedSkillOptions(ctx context.Context) error {
	categories := []domain.SkillOption{
		{ID: "1", Name: "言語", Icon: "code", SortOrder: 1},
		{ID: "2", Name: "フレームワーク", Icon: "code", SortOrder: 2},
		{ID: "3", Name: "データベース", Icon: "database", SortOrder: 3},
		{ID: "4", Name: "インフラ", Icon: "cloud", SortOrder: 4},
		{ID: "5", Name: "ツール", Icon: "wrench", SortOrder: 5},
		{ID: "6", Name: "その他", Icon: "wrench", SortOrder: 6},
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
		{ID: "1", Name: "初級", SortOrder: 1},
		{ID: "2", Name: "中級", SortOrder: 2},
		{ID: "3", Name: "上級", SortOrder: 3},
		{ID: "4", Name: "エキスパート", SortOrder: 4},
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

func (r *SQLiteProfileRepository) seedProfileLinkMasters(ctx context.Context) error {
	masters := []domain.ProfileLinkMaster{
		{ID: "1", Key: "github", Name: "GitHub", Icon: "github", Placeholder: "https://github.com/username", SortOrder: 1},
		{ID: "2", Key: "zenn", Name: "Zenn", Icon: "book-open", Placeholder: "https://zenn.dev/username", SortOrder: 2},
		{ID: "3", Key: "qiita", Name: "Qiita", Icon: "book-open", Placeholder: "https://qiita.com/username", SortOrder: 3},
		{ID: "4", Key: "website", Name: "個人サイト", Icon: "globe", Placeholder: "https://example.com", SortOrder: 4},
	}
	for _, master := range masters {
		if _, err := r.db.ExecContext(
			ctx,
			"INSERT OR IGNORE INTO profile_link_masters (id, key, name, icon, placeholder, sort_order) VALUES (?, ?, ?, ?, ?, ?)",
			master.ID,
			master.Key,
			master.Name,
			master.Icon,
			master.Placeholder,
			master.SortOrder,
		); err != nil {
			return fmt.Errorf("seed profile link master %s: %w", master.Key, err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) migrateProfileLinks(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT OR IGNORE INTO profile_links (profile_id, link_master_id, url, sort_order)
		SELECT profiles.id, profile_link_masters.id,
			CASE profile_link_masters.key
				WHEN 'github' THEN profiles.github_url
				WHEN 'zenn' THEN profiles.zenn_url
				WHEN 'qiita' THEN profiles.qiita_url
				WHEN 'website' THEN profiles.website_url
				ELSE ''
			END,
			profile_link_masters.sort_order
		FROM profiles
		JOIN profile_link_masters ON profile_link_masters.key IN ('github', 'zenn', 'qiita', 'website')
		WHERE CASE profile_link_masters.key
				WHEN 'github' THEN profiles.github_url
				WHEN 'zenn' THEN profiles.zenn_url
				WHEN 'qiita' THEN profiles.qiita_url
				WHEN 'website' THEN profiles.website_url
				ELSE ''
			END <> ''
	`)
	if err != nil {
		return fmt.Errorf("migrate profile links: %w", err)
	}

	return nil
}

func (r *SQLiteProfileRepository) listProfileLinks(ctx context.Context, profileID string) ([]domain.ProfileLink, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			profile_links.id,
			profile_link_masters.id,
			profile_link_masters.key,
			profile_link_masters.name,
			profile_link_masters.icon,
			profile_link_masters.placeholder,
			profile_links.url,
			profile_links.sort_order
		FROM profile_links
		JOIN profile_link_masters ON profile_link_masters.id = profile_links.link_master_id
		WHERE profile_links.profile_id = ?
		ORDER BY profile_links.sort_order ASC, profile_link_masters.sort_order ASC, profile_link_masters.name ASC
	`, profileID)
	if err != nil {
		return nil, fmt.Errorf("query profile links: %w", err)
	}
	defer rows.Close()

	links := []domain.ProfileLink{}
	for rows.Next() {
		var link domain.ProfileLink
		if err := rows.Scan(&link.ID, &link.LinkMasterID, &link.Key, &link.Name, &link.Icon, &link.Placeholder, &link.URL, &link.SortOrder); err != nil {
			return nil, fmt.Errorf("scan profile link: %w", err)
		}
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate profile links: %w", err)
	}

	return links, nil
}

func (r *SQLiteProfileRepository) saveProfileLinks(ctx context.Context, tx *sql.Tx, profileID string, links []domain.ProfileLink) error {
	if _, err := tx.ExecContext(ctx, "DELETE FROM profile_links WHERE profile_id = ?", profileID); err != nil {
		return fmt.Errorf("delete profile links: %w", err)
	}

	for index, link := range links {
		if strings.TrimSpace(link.LinkMasterID) == "" {
			continue
		}
		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO profile_links (profile_id, link_master_id, url, sort_order)
			VALUES (?, ?, ?, ?)`,
			profileID,
			link.LinkMasterID,
			link.URL,
			int64(index+1),
		); err != nil {
			return fmt.Errorf("insert profile link: %w", err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) seedSkills(ctx context.Context) error {
	skills := []struct {
		id                 string
		skillMasterID      string
		experience         int64
		proficiencyLevelID string
		sortOrder          int64
	}{
		{id: "1", skillMasterID: "1", experience: 3, proficiencyLevelID: "3", sortOrder: 1},
		{id: "2", skillMasterID: "4", experience: 3, proficiencyLevelID: "3", sortOrder: 2},
		{id: "3", skillMasterID: "6", experience: 2, proficiencyLevelID: "2", sortOrder: 3},
		{id: "4", skillMasterID: "7", experience: 2, proficiencyLevelID: "2", sortOrder: 4},
		{id: "5", skillMasterID: "9", experience: 2, proficiencyLevelID: "2", sortOrder: 5},
		{id: "6", skillMasterID: "10", experience: 1, proficiencyLevelID: "1", sortOrder: 6},
		{id: "7", skillMasterID: "11", experience: 4, proficiencyLevelID: "3", sortOrder: 7},
	}
	for _, skill := range skills {
		if _, err := r.db.ExecContext(
			ctx,
			"INSERT OR IGNORE INTO skills (id, skill_master_id, experience, proficiency_level_id, sort_order) VALUES (?, ?, ?, ?, ?)",
			skill.id,
			skill.skillMasterID,
			skill.experience,
			skill.proficiencyLevelID,
			skill.sortOrder,
		); err != nil {
			return fmt.Errorf("seed skill %s: %w", skill.id, err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) seedSkillMasters(ctx context.Context) error {
	skillMasters := []domain.SkillMaster{
		{ID: "1", Name: "TypeScript", CategoryID: "1", SortOrder: 1},
		{ID: "2", Name: "JavaScript", CategoryID: "1", SortOrder: 2},
		{ID: "3", Name: "Go", CategoryID: "1", SortOrder: 3},
		{ID: "4", Name: "React", CategoryID: "2", SortOrder: 4},
		{ID: "5", Name: "Next.js", CategoryID: "2", SortOrder: 5},
		{ID: "6", Name: "Node.js", CategoryID: "2", SortOrder: 6},
		{ID: "7", Name: "PostgreSQL", CategoryID: "3", SortOrder: 7},
		{ID: "8", Name: "SQLite", CategoryID: "3", SortOrder: 8},
		{ID: "9", Name: "Docker", CategoryID: "4", SortOrder: 9},
		{ID: "10", Name: "AWS", CategoryID: "4", SortOrder: 10},
		{ID: "11", Name: "Git", CategoryID: "5", SortOrder: 11},
	}
	for _, skillMaster := range skillMasters {
		if _, err := r.db.ExecContext(
			ctx,
			"INSERT OR IGNORE INTO skill_masters (id, name, category_id, sort_order) VALUES (?, ?, ?, ?)",
			skillMaster.ID,
			skillMaster.Name,
			skillMaster.CategoryID,
			skillMaster.SortOrder,
		); err != nil {
			return fmt.Errorf("seed skill master %s: %w", skillMaster.ID, err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) seedJobEmploymentTypes(ctx context.Context) error {
	employmentTypes := []domain.JobEmploymentType{
		{ID: "1", Name: "正社員", SortOrder: 1},
		{ID: "2", Name: "契約社員", SortOrder: 2},
		{ID: "3", Name: "業務委託", SortOrder: 3},
		{ID: "4", Name: "派遣", SortOrder: 4},
		{ID: "5", Name: "アルバイト", SortOrder: 5},
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

func (r *SQLiteProfileRepository) seedJobCompanies(ctx context.Context) error {
	companyNames := []string{"株式会社A", "株式会社B", "フリーランス"}
	for _, companyName := range companyNames {
		if _, err := r.db.ExecContext(
			ctx,
			"INSERT OR IGNORE INTO job_companies (name, url) VALUES (?, '')",
			companyName,
		); err != nil {
			return fmt.Errorf("seed job company %s: %w", companyName, err)
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
	}{
		{id: "1", company: "株式会社A", displayName: "株式会社A", startYear: 2023, startMonth: 1, endYear: nil, endMonth: nil, employmentTypeID: "1", projectCount: 5},
		{id: "2", company: "株式会社B", displayName: "株式会社B", startYear: 2021, startMonth: 4, endYear: int64Ptr(2022), endMonth: int64Ptr(12), employmentTypeID: "1", projectCount: 4},
		{id: "3", company: "フリーランス", displayName: "フリーランス", startYear: 2020, startMonth: 1, endYear: int64Ptr(2021), endMonth: int64Ptr(3), employmentTypeID: "3", projectCount: 3},
	}
	for _, jobHistory := range jobHistories {
		if _, err := r.db.ExecContext(
			ctx,
			`INSERT OR IGNORE INTO job_histories (
				id,
				company_id,
				display_name,
				start_year,
				start_month,
				end_year,
				end_month,
				employment_type_id,
				project_count
			)
			VALUES (
				?,
				(SELECT id FROM job_companies WHERE name = ? LIMIT 1),
				?,
				?,
				?,
				?,
				?,
				?,
				?
			)`,
			jobHistory.id,
			jobHistory.company,
			jobHistory.displayName,
			jobHistory.startYear,
			jobHistory.startMonth,
			toNullInt64(jobHistory.endYear),
			toNullInt64(jobHistory.endMonth),
			jobHistory.employmentTypeID,
			jobHistory.projectCount,
		); err != nil {
			return fmt.Errorf("seed job history %s: %w", jobHistory.id, err)
		}
	}

	return nil
}

func (r *SQLiteProfileRepository) backfillJobCompanies(ctx context.Context) error {
	hasCompanyColumn, err := r.tableHasColumn(ctx, "job_histories", "company")
	if err != nil {
		return err
	}
	if !hasCompanyColumn {
		return nil
	}

	if _, err := r.db.ExecContext(ctx, `
		INSERT OR IGNORE INTO job_companies (name, url)
		SELECT DISTINCT company, ''
		FROM job_histories
		WHERE company IS NOT NULL AND company <> ''
		ORDER BY company
	`); err != nil {
		return fmt.Errorf("backfill job companies: %w", err)
	}

	return nil
}

func (r *SQLiteProfileRepository) seedProjects(ctx context.Context) error {
	projects := []domain.Project{
		{
			ID:           "1",
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
			ID:           "2",
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

func (r *SQLiteProfileRepository) tableColumnNotNull(ctx context.Context, tableName string, columnName string) (bool, error) {
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
			return notNull != 0, nil
		}
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("iterate %s columns: %w", tableName, err)
	}

	return false, fmt.Errorf("%s.%s column not found", tableName, columnName)
}

func (r *SQLiteProfileRepository) tableColumnType(ctx context.Context, tableName string, columnName string) (string, error) {
	rows, err := r.db.QueryContext(ctx, "PRAGMA table_info("+tableName+")")
	if err != nil {
		return "", fmt.Errorf("inspect %s columns: %w", tableName, err)
	}
	defer func() {
		_ = rows.Close()
	}()

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
			return "", fmt.Errorf("scan %s column: %w", tableName, err)
		}
		if name == columnName {
			return columnType, nil
		}
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("iterate %s columns: %w", tableName, err)
	}

	return "", fmt.Errorf("%s.%s column not found", tableName, columnName)
}

func toDomainProfile(row dbgen.Profile, qualificationRows []dbgen.ProfileQualification, links []domain.ProfileLink) (*domain.Profile, error) {
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
	profile.Links = links

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
		options = append(options, *toDomainSkillProficiencyLevelOption(row))
	}

	return options
}

func toDomainSkillProficiencyLevelOption(row dbgen.SkillProficiencyLevel) *domain.SkillOption {
	return &domain.SkillOption{
		ID:        row.ID,
		Name:      row.Name,
		SortOrder: row.SortOrder,
	}
}

func toDomainSkillMasters(rows []dbgen.ListSkillMastersRow) []domain.SkillMaster {
	skillMasters := make([]domain.SkillMaster, 0, len(rows))
	for _, row := range rows {
		skillMasters = append(skillMasters, domain.SkillMaster{
			ID:           row.ID,
			Name:         row.Name,
			CategoryID:   row.CategoryID,
			CategoryName: row.CategoryName,
			SortOrder:    row.SortOrder,
		})
	}

	return skillMasters
}

func toDomainSkillMaster(row dbgen.GetSkillMasterRow) *domain.SkillMaster {
	return &domain.SkillMaster{
		ID:           row.ID,
		Name:         row.Name,
		CategoryID:   row.CategoryID,
		CategoryName: row.CategoryName,
		SortOrder:    row.SortOrder,
	}
}

func toDomainSkills(rows []dbgen.ListSkillsRow) []domain.Skill {
	skills := make([]domain.Skill, 0, len(rows))
	for _, row := range rows {
		skills = append(skills, domain.Skill{
			ID:                 row.ID,
			SkillMasterID:      row.SkillMasterID,
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
		SkillMasterID:      row.SkillMasterID,
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
			CompanyID:        row.CompanyID,
			Company:          row.Company,
			DisplayName:      row.DisplayName,
			StartYear:        row.StartYear,
			StartMonth:       row.StartMonth,
			EndYear:          fromNullInt64(row.EndYear),
			EndMonth:         fromNullInt64(row.EndMonth),
			EmploymentTypeID: row.EmploymentTypeID,
			EmploymentType:   row.EmploymentType,
			ProjectCount:     row.ProjectCount,
		})
	}

	return jobHistories
}

func toDomainJobHistory(row dbgen.GetJobHistoryRow) *domain.JobHistory {
	return &domain.JobHistory{
		ID:               row.ID,
		CompanyID:        row.CompanyID,
		Company:          row.Company,
		DisplayName:      row.DisplayName,
		StartYear:        row.StartYear,
		StartMonth:       row.StartMonth,
		EndYear:          fromNullInt64(row.EndYear),
		EndMonth:         fromNullInt64(row.EndMonth),
		EmploymentTypeID: row.EmploymentTypeID,
		EmploymentType:   row.EmploymentType,
		ProjectCount:     row.ProjectCount,
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

func syncLegacyProfileLinkColumns(profile *domain.Profile) {
	profile.GitHubURL = ""
	profile.ZennURL = ""
	profile.QiitaURL = ""
	profile.WebsiteURL = ""

	for _, link := range profile.Links {
		switch link.Key {
		case "github":
			profile.GitHubURL = link.URL
		case "zenn":
			profile.ZennURL = link.URL
		case "qiita":
			profile.QiitaURL = link.URL
		case "website":
			profile.WebsiteURL = link.URL
		}
	}
}

func qualificationID(qualification domain.Qualification, index int, now time.Time) string {
	if _, err := strconv.ParseInt(qualification.ID, 10, 64); err == nil && qualification.ID != "" {
		return qualification.ID
	}

	return fmt.Sprintf("%d", now.UnixNano()+int64(index)+1)
}

func skillID(now time.Time) string {
	return fmt.Sprintf("%d", now.UnixNano())
}

func projectID(now time.Time) string {
	return fmt.Sprintf("%d", now.UnixNano())
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
