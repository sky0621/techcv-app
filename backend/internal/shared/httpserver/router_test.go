package httpserver

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/sky0621/techcv-app/backend/internal/domain"
)

func TestProfileRoutes(t *testing.T) {
	repository := newTestProfileRepository()
	router := NewRouter(repository, repository)

	getReq := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", getRec.Code)
	}

	var getResp struct {
		Profile struct {
			DisplayName string `json:"displayName"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(getRec.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("failed to decode get response: %v", err)
	}

	if getResp.Profile.DisplayName != "Initial User" {
		t.Fatalf("expected initial displayName, got %v", getResp.Profile.DisplayName)
	}

	body := []byte(`{
		"displayName":"Sky Sample",
		"location":"Tokyo",
		"email":"me@example.com",
		"bio":"Backend engineer",
		"githubUrl":"https://github.com/sky0621",
		"zennUrl":"https://zenn.dev/sky0621",
		"qiitaUrl":"https://qiita.com/sky0621",
		"websiteUrl":"https://example.com",
		"occupation":"Software Engineer",
		"employmentType":"Freelance",
		"workStyle":"Full remote",
		"qualifications":[
			{
				"name":"AWS Certified Solutions Architect",
				"acquiredDate":"2026-04-26",
				"organization":"Amazon Web Services",
				"url":"https://aws.amazon.com/certification/certified-solutions-architect-associate/",
				"memo":"Associate"
			},
			{
				"name":"基本情報技術者",
				"acquiredDate":"2020-10-01",
				"organization":"IPA",
				"url":"https://www.ipa.go.jp/shiken/kubun/fe.html",
				"memo":""
			}
		],
		"visibilitySettings":{"email":false,"location":false}
	}`)

	putReq := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewReader(body))
	putRec := httptest.NewRecorder()
	router.ServeHTTP(putRec, putReq)

	if putRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", putRec.Code)
	}

	var putResp struct {
		Profile struct {
			DisplayName    string `json:"displayName"`
			Location       string `json:"location"`
			Email          string `json:"email"`
			Bio            string `json:"bio"`
			GithubURL      string `json:"githubUrl"`
			ZennURL        string `json:"zennUrl"`
			QiitaURL       string `json:"qiitaUrl"`
			WebsiteURL     string `json:"websiteUrl"`
			Occupation     string `json:"occupation"`
			EmploymentType string `json:"employmentType"`
			WorkStyle      string `json:"workStyle"`
			Qualifications []struct {
				ID           string `json:"id"`
				Name         string `json:"name"`
				AcquiredDate string `json:"acquiredDate"`
				Organization string `json:"organization"`
				URL          string `json:"url"`
				Memo         string `json:"memo"`
			} `json:"qualifications"`
			Visibility struct {
				Email    bool `json:"email"`
				Location bool `json:"location"`
			} `json:"visibilitySettings"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(putRec.Body.Bytes(), &putResp); err != nil {
		t.Fatalf("failed to decode put response: %v", err)
	}

	if putResp.Profile.DisplayName != "Sky Sample" {
		t.Fatalf("expected updated displayName, got %v", putResp.Profile.DisplayName)
	}
	if putResp.Profile.Location != "Tokyo" {
		t.Fatalf("expected updated location, got %v", putResp.Profile.Location)
	}
	if putResp.Profile.Email != "me@example.com" {
		t.Fatalf("expected updated email, got %v", putResp.Profile.Email)
	}
	if putResp.Profile.Bio != "Backend engineer" {
		t.Fatalf("expected updated bio, got %v", putResp.Profile.Bio)
	}
	if putResp.Profile.GithubURL != "https://github.com/sky0621" {
		t.Fatalf("expected updated githubUrl, got %v", putResp.Profile.GithubURL)
	}
	if putResp.Profile.ZennURL != "https://zenn.dev/sky0621" {
		t.Fatalf("expected updated zennUrl, got %v", putResp.Profile.ZennURL)
	}
	if putResp.Profile.QiitaURL != "https://qiita.com/sky0621" {
		t.Fatalf("expected updated qiitaUrl, got %v", putResp.Profile.QiitaURL)
	}
	if putResp.Profile.WebsiteURL != "https://example.com" {
		t.Fatalf("expected updated websiteUrl, got %v", putResp.Profile.WebsiteURL)
	}
	if putResp.Profile.Occupation != "Software Engineer" {
		t.Fatalf("expected updated occupation, got %v", putResp.Profile.Occupation)
	}
	if putResp.Profile.EmploymentType != "Freelance" {
		t.Fatalf("expected updated employmentType, got %v", putResp.Profile.EmploymentType)
	}
	if putResp.Profile.WorkStyle != "Full remote" {
		t.Fatalf("expected updated workStyle, got %v", putResp.Profile.WorkStyle)
	}
	if len(putResp.Profile.Qualifications) != 2 {
		t.Fatalf("expected two qualifications, got %#v", putResp.Profile.Qualifications)
	}
	if putResp.Profile.Qualifications[0].Name != "AWS Certified Solutions Architect" {
		t.Fatalf("expected first qualification to preserve order, got %#v", putResp.Profile.Qualifications[0])
	}
	if putResp.Profile.Qualifications[0].URL != "https://aws.amazon.com/certification/certified-solutions-architect-associate/" {
		t.Fatalf("expected first qualification URL to persist, got %#v", putResp.Profile.Qualifications[0])
	}
	if putResp.Profile.Qualifications[1].Name != "基本情報技術者" {
		t.Fatalf("expected second qualification to preserve order, got %#v", putResp.Profile.Qualifications[1])
	}
	if putResp.Profile.Visibility.Email != false || putResp.Profile.Visibility.Location != false {
		t.Fatalf("expected updated visibility settings, got %+v", putResp.Profile.Visibility)
	}

	getUpdatedReq := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	getUpdatedRec := httptest.NewRecorder()
	router.ServeHTTP(getUpdatedRec, getUpdatedReq)

	var getUpdatedResp struct {
		Profile struct {
			DisplayName    string `json:"displayName"`
			Location       string `json:"location"`
			Email          string `json:"email"`
			Bio            string `json:"bio"`
			GithubURL      string `json:"githubUrl"`
			ZennURL        string `json:"zennUrl"`
			QiitaURL       string `json:"qiitaUrl"`
			WebsiteURL     string `json:"websiteUrl"`
			Occupation     string `json:"occupation"`
			EmploymentType string `json:"employmentType"`
			WorkStyle      string `json:"workStyle"`
			Qualifications []struct {
				ID           string `json:"id"`
				Name         string `json:"name"`
				AcquiredDate string `json:"acquiredDate"`
				Organization string `json:"organization"`
				URL          string `json:"url"`
				Memo         string `json:"memo"`
			} `json:"qualifications"`
			Visibility struct {
				Email    bool `json:"email"`
				Location bool `json:"location"`
			} `json:"visibilitySettings"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(getUpdatedRec.Body.Bytes(), &getUpdatedResp); err != nil {
		t.Fatalf("failed to decode updated get response: %v", err)
	}

	if getUpdatedResp.Profile.DisplayName != "Sky Sample" {
		t.Fatalf("expected persisted displayName, got %v", getUpdatedResp.Profile.DisplayName)
	}
	if getUpdatedResp.Profile.Location != "Tokyo" {
		t.Fatalf("expected persisted location, got %v", getUpdatedResp.Profile.Location)
	}
	if getUpdatedResp.Profile.Email != "me@example.com" {
		t.Fatalf("expected persisted email, got %v", getUpdatedResp.Profile.Email)
	}
	if getUpdatedResp.Profile.Bio != "Backend engineer" {
		t.Fatalf("expected persisted bio, got %v", getUpdatedResp.Profile.Bio)
	}
	if getUpdatedResp.Profile.GithubURL != "https://github.com/sky0621" {
		t.Fatalf("expected persisted githubUrl, got %v", getUpdatedResp.Profile.GithubURL)
	}
	if getUpdatedResp.Profile.ZennURL != "https://zenn.dev/sky0621" {
		t.Fatalf("expected persisted zennUrl, got %v", getUpdatedResp.Profile.ZennURL)
	}
	if getUpdatedResp.Profile.QiitaURL != "https://qiita.com/sky0621" {
		t.Fatalf("expected persisted qiitaUrl, got %v", getUpdatedResp.Profile.QiitaURL)
	}
	if getUpdatedResp.Profile.WebsiteURL != "https://example.com" {
		t.Fatalf("expected persisted websiteUrl, got %v", getUpdatedResp.Profile.WebsiteURL)
	}
	if getUpdatedResp.Profile.Occupation != "Software Engineer" {
		t.Fatalf("expected persisted occupation, got %v", getUpdatedResp.Profile.Occupation)
	}
	if getUpdatedResp.Profile.EmploymentType != "Freelance" {
		t.Fatalf("expected persisted employmentType, got %v", getUpdatedResp.Profile.EmploymentType)
	}
	if getUpdatedResp.Profile.WorkStyle != "Full remote" {
		t.Fatalf("expected persisted workStyle, got %v", getUpdatedResp.Profile.WorkStyle)
	}
	if len(getUpdatedResp.Profile.Qualifications) != 2 {
		t.Fatalf("expected persisted qualifications, got %#v", getUpdatedResp.Profile.Qualifications)
	}
	if getUpdatedResp.Profile.Qualifications[0].Name != "AWS Certified Solutions Architect" {
		t.Fatalf("expected persisted first qualification, got %#v", getUpdatedResp.Profile.Qualifications[0])
	}
	if getUpdatedResp.Profile.Qualifications[1].Organization != "IPA" {
		t.Fatalf("expected persisted second qualification, got %#v", getUpdatedResp.Profile.Qualifications[1])
	}
	if getUpdatedResp.Profile.Qualifications[1].URL != "https://www.ipa.go.jp/shiken/kubun/fe.html" {
		t.Fatalf("expected persisted second qualification URL, got %#v", getUpdatedResp.Profile.Qualifications[1])
	}
	if getUpdatedResp.Profile.Visibility.Email != false || getUpdatedResp.Profile.Visibility.Location != false {
		t.Fatalf("expected persisted visibility settings, got %+v", getUpdatedResp.Profile.Visibility)
	}

	replacementBody := []byte(`{
		"displayName":"Sky Sample",
		"location":"Tokyo",
		"email":"me@example.com",
		"bio":"Backend engineer",
		"githubUrl":"https://github.com/sky0621",
		"zennUrl":"https://zenn.dev/sky0621",
		"qiitaUrl":"https://qiita.com/sky0621",
		"websiteUrl":"https://example.com",
		"occupation":"Software Engineer",
		"employmentType":"Freelance",
		"workStyle":"Full remote",
		"qualifications":[
			{
				"id":"replacement_qualification",
				"name":"Google Cloud Professional Cloud Architect",
				"acquiredDate":"2025-01-01",
				"organization":"Google Cloud",
				"url":"https://cloud.google.com/learn/certification/cloud-architect",
				"memo":"更新済み"
			}
		],
		"visibilitySettings":{"email":false,"location":false}
	}`)

	replacementReq := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewReader(replacementBody))
	replacementRec := httptest.NewRecorder()
	router.ServeHTTP(replacementRec, replacementReq)

	if replacementRec.Code != http.StatusOK {
		t.Fatalf("expected replacement status 200, got %d", replacementRec.Code)
	}

	getReplacedReq := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	getReplacedRec := httptest.NewRecorder()
	router.ServeHTTP(getReplacedRec, getReplacedReq)

	var getReplacedResp struct {
		Profile struct {
			Qualifications []struct {
				ID           string `json:"id"`
				Name         string `json:"name"`
				AcquiredDate string `json:"acquiredDate"`
				Organization string `json:"organization"`
				URL          string `json:"url"`
				Memo         string `json:"memo"`
			} `json:"qualifications"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(getReplacedRec.Body.Bytes(), &getReplacedResp); err != nil {
		t.Fatalf("failed to decode replaced get response: %v", err)
	}
	if len(getReplacedResp.Profile.Qualifications) != 1 {
		t.Fatalf("expected replaced qualifications, got %#v", getReplacedResp.Profile.Qualifications)
	}
	if getReplacedResp.Profile.Qualifications[0].ID != "replacement_qualification" {
		t.Fatalf("expected replacement qualification id to persist, got %#v", getReplacedResp.Profile.Qualifications[0])
	}
	if getReplacedResp.Profile.Qualifications[0].URL != "https://cloud.google.com/learn/certification/cloud-architect" {
		t.Fatalf("expected replacement qualification URL to persist, got %#v", getReplacedResp.Profile.Qualifications[0])
	}
}

func TestSkillOptionsRoute(t *testing.T) {
	repository := newTestProfileRepository()
	router := NewRouter(repository, repository)

	req := httptest.NewRequest(http.MethodGet, "/api/skills/options", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp struct {
		Categories []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Icon      string `json:"icon"`
			SortOrder int64  `json:"sortOrder"`
		} `json:"categories"`
		ProficiencyLevels []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			SortOrder int64  `json:"sortOrder"`
		} `json:"proficiencyLevels"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode skill options response: %v", err)
	}

	if len(resp.Categories) != 6 {
		t.Fatalf("expected six categories, got %#v", resp.Categories)
	}
	if resp.Categories[0].ID != "skill_category_language" || resp.Categories[0].Name != "言語" {
		t.Fatalf("unexpected first category: %#v", resp.Categories[0])
	}
	if resp.Categories[0].Icon != "code" {
		t.Fatalf("unexpected first category icon: %#v", resp.Categories[0])
	}
	if len(resp.ProficiencyLevels) != 4 {
		t.Fatalf("expected four proficiency levels, got %#v", resp.ProficiencyLevels)
	}
	if resp.ProficiencyLevels[0].ID != "skill_proficiency_beginner" || resp.ProficiencyLevels[0].Name != "初級" {
		t.Fatalf("unexpected first proficiency level: %#v", resp.ProficiencyLevels[0])
	}
}

func TestSkillCategoryMutationRoutes(t *testing.T) {
	repository := newTestProfileRepository()
	router := NewRouter(repository, repository)

	createBody := []byte(`{
		"id":"skill_category_backend",
		"name":"バックエンド",
		"icon":"code"
	}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/skills/categories", bytes.NewReader(createBody))
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status 201, got %d", createRec.Code)
	}

	var createResp struct {
		Category struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Icon      string `json:"icon"`
			SortOrder int64  `json:"sortOrder"`
		} `json:"category"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}
	if createResp.Category.ID != "skill_category_backend" || createResp.Category.Name != "バックエンド" {
		t.Fatalf("unexpected created category: %#v", createResp.Category)
	}

	updateBody := []byte(`{
		"name":"バックエンド・API",
		"icon":"database"
	}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/skills/categories/skill_category_backend", bytes.NewReader(updateBody))
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected update status 200, got %d", updateRec.Code)
	}

	var updateResp struct {
		Category struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Icon string `json:"icon"`
		} `json:"category"`
	}
	if err := json.Unmarshal(updateRec.Body.Bytes(), &updateResp); err != nil {
		t.Fatalf("failed to decode update response: %v", err)
	}
	if updateResp.Category.ID != "skill_category_backend" ||
		updateResp.Category.Name != "バックエンド・API" ||
		updateResp.Category.Icon != "database" {
		t.Fatalf("unexpected updated category: %#v", updateResp.Category)
	}
}

func TestSkillRoutes(t *testing.T) {
	repository := newTestProfileRepository()
	router := NewRouter(repository, repository)

	listReq := httptest.NewRequest(http.MethodGet, "/api/skills", nil)
	listRec := httptest.NewRecorder()
	router.ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("expected list status 200, got %d", listRec.Code)
	}

	var listResp struct {
		Skills []struct {
			ID                 string `json:"id"`
			Name               string `json:"name"`
			CategoryID         string `json:"categoryId"`
			Category           string `json:"category"`
			Experience         string `json:"experience"`
			ProficiencyLevelID string `json:"proficiencyLevelId"`
			Proficiency        string `json:"proficiency"`
		} `json:"skills"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}
	if len(listResp.Skills) != 1 {
		t.Fatalf("expected one seeded skill, got %#v", listResp.Skills)
	}
	if listResp.Skills[0].Category != "言語" || listResp.Skills[0].Proficiency != "上級" {
		t.Fatalf("expected skill master names, got %#v", listResp.Skills[0])
	}

	createBody := []byte(`{
		"name":"React",
		"categoryId":"skill_category_framework",
		"experience":"3年",
		"proficiencyLevelId":"skill_proficiency_advanced"
	}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/skills", bytes.NewReader(createBody))
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status 201, got %d", createRec.Code)
	}

	var createResp struct {
		Skill struct {
			ID                 string `json:"id"`
			Name               string `json:"name"`
			CategoryID         string `json:"categoryId"`
			Category           string `json:"category"`
			Experience         string `json:"experience"`
			ProficiencyLevelID string `json:"proficiencyLevelId"`
			Proficiency        string `json:"proficiency"`
		} `json:"skill"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}
	if createResp.Skill.ID == "" || createResp.Skill.Category != "フレームワーク" {
		t.Fatalf("unexpected created skill: %#v", createResp.Skill)
	}

	updateBody := []byte(`{
		"name":"React",
		"categoryId":"skill_category_framework",
		"experience":"4年",
		"proficiencyLevelId":"skill_proficiency_expert"
	}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/skills/"+createResp.Skill.ID, bytes.NewReader(updateBody))
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected update status 200, got %d", updateRec.Code)
	}

	var updateResp struct {
		Skill struct {
			ID          string `json:"id"`
			Experience  string `json:"experience"`
			Proficiency string `json:"proficiency"`
		} `json:"skill"`
	}
	if err := json.Unmarshal(updateRec.Body.Bytes(), &updateResp); err != nil {
		t.Fatalf("failed to decode update response: %v", err)
	}
	if updateResp.Skill.Experience != "4年" || updateResp.Skill.Proficiency != "エキスパート" {
		t.Fatalf("unexpected updated skill: %#v", updateResp.Skill)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/skills/"+createResp.Skill.ID, nil)
	deleteRec := httptest.NewRecorder()
	router.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("expected delete status 204, got %d", deleteRec.Code)
	}
}

type testProfileRepository struct {
	mu                sync.RWMutex
	profile           *domain.Profile
	categories        []domain.SkillOption
	proficiencyLevels []domain.SkillOption
	skills            []domain.Skill
	nextSkillID       int
}

func newTestProfileRepository() *testProfileRepository {
	now := time.Now().UTC()

	return &testProfileRepository{
		profile: &domain.Profile{
			ID:                 "profile_01",
			UserID:             "user_01",
			FullName:           "Initial User",
			VisibilitySettings: map[string]any{"email": false},
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		categories: []domain.SkillOption{
			{ID: "skill_category_language", Name: "言語", Icon: "code", SortOrder: 1},
			{ID: "skill_category_framework", Name: "フレームワーク", Icon: "code", SortOrder: 2},
			{ID: "skill_category_database", Name: "データベース", Icon: "database", SortOrder: 3},
			{ID: "skill_category_infrastructure", Name: "インフラ", Icon: "cloud", SortOrder: 4},
			{ID: "skill_category_tool", Name: "ツール", Icon: "wrench", SortOrder: 5},
			{ID: "skill_category_other", Name: "その他", Icon: "wrench", SortOrder: 6},
		},
		proficiencyLevels: []domain.SkillOption{
			{ID: "skill_proficiency_beginner", Name: "初級", SortOrder: 1},
			{ID: "skill_proficiency_intermediate", Name: "中級", SortOrder: 2},
			{ID: "skill_proficiency_advanced", Name: "上級", SortOrder: 3},
			{ID: "skill_proficiency_expert", Name: "エキスパート", SortOrder: 4},
		},
		skills: []domain.Skill{
			{
				ID:                 "skill_typescript",
				Name:               "TypeScript",
				CategoryID:         "skill_category_language",
				CategoryName:       "言語",
				Experience:         "3年",
				ProficiencyLevelID: "skill_proficiency_advanced",
				ProficiencyName:    "上級",
				SortOrder:          1,
			},
		},
		nextSkillID: 1,
	}
}

func (r *testProfileRepository) Get(_ context.Context) (*domain.Profile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	profileCopy := *r.profile
	return &profileCopy, nil
}

func (r *testProfileRepository) Save(_ context.Context, profile *domain.Profile) (*domain.Profile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	profileCopy := *profile
	profileCopy.UpdatedAt = time.Now().UTC()
	if profileCopy.CreatedAt.IsZero() {
		profileCopy.CreatedAt = profileCopy.UpdatedAt
	}

	r.profile = &profileCopy

	savedCopy := *r.profile
	return &savedCopy, nil
}

func (r *testProfileRepository) ListSkillOptions(context.Context) (*domain.SkillOptions, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make([]domain.SkillOption, len(r.categories))
	copy(categories, r.categories)
	proficiencyLevels := make([]domain.SkillOption, len(r.proficiencyLevels))
	copy(proficiencyLevels, r.proficiencyLevels)

	return &domain.SkillOptions{
		Categories:        categories,
		ProficiencyLevels: proficiencyLevels,
	}, nil
}

func (r *testProfileRepository) CreateSkillCategory(_ context.Context, input domain.SkillCategoryInput) (*domain.SkillOption, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	category := domain.SkillOption{
		ID:        input.ID,
		Name:      input.Name,
		Icon:      input.Icon,
		SortOrder: int64(len(r.categories) + 1),
	}
	r.categories = append(r.categories, category)

	return &category, nil
}

func (r *testProfileRepository) UpdateSkillCategory(_ context.Context, id string, input domain.SkillCategoryInput) (*domain.SkillOption, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, category := range r.categories {
		if category.ID == id {
			r.categories[index].Name = input.Name
			r.categories[index].Icon = input.Icon
			result := r.categories[index]
			return &result, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *testProfileRepository) ListSkills(context.Context) ([]domain.Skill, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	skills := make([]domain.Skill, len(r.skills))
	copy(skills, r.skills)

	return skills, nil
}

func (r *testProfileRepository) CreateSkill(_ context.Context, input domain.SkillInput) (*domain.Skill, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	category, ok := r.findCategory(input.CategoryID)
	if !ok {
		return nil, sql.ErrNoRows
	}
	proficiencyLevel, ok := r.findProficiencyLevel(input.ProficiencyLevelID)
	if !ok {
		return nil, sql.ErrNoRows
	}

	r.nextSkillID++
	skill := domain.Skill{
		ID:                 "skill_test_new",
		Name:               input.Name,
		CategoryID:         input.CategoryID,
		CategoryName:       category.Name,
		Experience:         input.Experience,
		ProficiencyLevelID: input.ProficiencyLevelID,
		ProficiencyName:    proficiencyLevel.Name,
		SortOrder:          int64(len(r.skills) + 1),
	}
	r.skills = append(r.skills, skill)

	return &skill, nil
}

func (r *testProfileRepository) UpdateSkill(_ context.Context, id string, input domain.SkillInput) (*domain.Skill, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	category, ok := r.findCategory(input.CategoryID)
	if !ok {
		return nil, sql.ErrNoRows
	}
	proficiencyLevel, ok := r.findProficiencyLevel(input.ProficiencyLevelID)
	if !ok {
		return nil, sql.ErrNoRows
	}

	for index, skill := range r.skills {
		if skill.ID == id {
			r.skills[index].Name = input.Name
			r.skills[index].CategoryID = input.CategoryID
			r.skills[index].CategoryName = category.Name
			r.skills[index].Experience = input.Experience
			r.skills[index].ProficiencyLevelID = input.ProficiencyLevelID
			r.skills[index].ProficiencyName = proficiencyLevel.Name
			result := r.skills[index]
			return &result, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *testProfileRepository) DeleteSkill(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, skill := range r.skills {
		if skill.ID == id {
			r.skills = append(r.skills[:index], r.skills[index+1:]...)
			return nil
		}
	}

	return sql.ErrNoRows
}

func (r *testProfileRepository) findCategory(id string) (domain.SkillOption, bool) {
	for _, category := range r.categories {
		if category.ID == id {
			return category, true
		}
	}

	return domain.SkillOption{}, false
}

func (r *testProfileRepository) findProficiencyLevel(id string) (domain.SkillOption, bool) {
	for _, level := range r.proficiencyLevels {
		if level.ID == id {
			return level, true
		}
	}

	return domain.SkillOption{}, false
}
