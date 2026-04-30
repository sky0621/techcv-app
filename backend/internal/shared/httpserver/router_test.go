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
	router := NewRouter(repository, repository, repository, repository)

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
	router := NewRouter(repository, repository, repository, repository)

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
		SkillMasters []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			CategoryID string `json:"categoryId"`
			Category   string `json:"category"`
		} `json:"skillMasters"`
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
	if len(resp.SkillMasters) != 1 {
		t.Fatalf("expected one skill master, got %#v", resp.SkillMasters)
	}
	if resp.SkillMasters[0].ID != "skill_master_typescript" ||
		resp.SkillMasters[0].Name != "TypeScript" ||
		resp.SkillMasters[0].CategoryID != "skill_category_language" ||
		resp.SkillMasters[0].Category != "言語" {
		t.Fatalf("unexpected first skill master: %#v", resp.SkillMasters[0])
	}
}

func TestSkillCategoryMutationRoutes(t *testing.T) {
	repository := newTestProfileRepository()
	router := NewRouter(repository, repository, repository, repository)

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

func TestSkillMasterMutationRoutes(t *testing.T) {
	repository := newTestProfileRepository()
	router := NewRouter(repository, repository, repository, repository)

	createBody := []byte(`{
		"id":"skill_master_kotlin",
		"name":"Kotlin",
		"categoryId":"skill_category_language"
	}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/skills/masters", bytes.NewReader(createBody))
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status 201, got %d", createRec.Code)
	}

	var createResp struct {
		SkillMaster struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			CategoryID string `json:"categoryId"`
			Category   string `json:"category"`
		} `json:"skillMaster"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}
	if createResp.SkillMaster.ID != "skill_master_kotlin" ||
		createResp.SkillMaster.Name != "Kotlin" ||
		createResp.SkillMaster.Category != "言語" {
		t.Fatalf("unexpected created skill master: %#v", createResp.SkillMaster)
	}

	updateBody := []byte(`{
		"name":"Kotlin/JVM",
		"categoryId":"skill_category_framework"
	}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/skills/masters/skill_master_kotlin", bytes.NewReader(updateBody))
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected update status 200, got %d", updateRec.Code)
	}

	var updateResp struct {
		SkillMaster struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			CategoryID string `json:"categoryId"`
			Category   string `json:"category"`
		} `json:"skillMaster"`
	}
	if err := json.Unmarshal(updateRec.Body.Bytes(), &updateResp); err != nil {
		t.Fatalf("failed to decode update response: %v", err)
	}
	if updateResp.SkillMaster.ID != "skill_master_kotlin" ||
		updateResp.SkillMaster.Name != "Kotlin/JVM" ||
		updateResp.SkillMaster.CategoryID != "skill_category_framework" ||
		updateResp.SkillMaster.Category != "フレームワーク" {
		t.Fatalf("unexpected updated skill master: %#v", updateResp.SkillMaster)
	}
}

func TestSkillRoutes(t *testing.T) {
	repository := newTestProfileRepository()
	router := NewRouter(repository, repository, repository, repository)

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
			Experience         int64  `json:"experience"`
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
		"experience":3,
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
			Experience         int64  `json:"experience"`
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
		"experience":4,
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
			Experience  int64  `json:"experience"`
			Proficiency string `json:"proficiency"`
		} `json:"skill"`
	}
	if err := json.Unmarshal(updateRec.Body.Bytes(), &updateResp); err != nil {
		t.Fatalf("failed to decode update response: %v", err)
	}
	if updateResp.Skill.Experience != 4 || updateResp.Skill.Proficiency != "エキスパート" {
		t.Fatalf("unexpected updated skill: %#v", updateResp.Skill)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/skills/"+createResp.Skill.ID, nil)
	deleteRec := httptest.NewRecorder()
	router.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("expected delete status 204, got %d", deleteRec.Code)
	}
}

func TestJobHistoryRoutes(t *testing.T) {
	repository := newTestProfileRepository()
	router := NewRouter(repository, repository, repository, repository)

	listReq := httptest.NewRequest(http.MethodGet, "/api/job-histories", nil)
	listRec := httptest.NewRecorder()
	router.ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("expected list status 200, got %d", listRec.Code)
	}

	var listResp struct {
		JobHistories []struct {
			ID               string `json:"id"`
			Company          string `json:"company"`
			DisplayName      string `json:"displayName"`
			StartYear        int64  `json:"startYear"`
			StartMonth       int64  `json:"startMonth"`
			EndYear          *int64 `json:"endYear"`
			EndMonth         *int64 `json:"endMonth"`
			EmploymentTypeID string `json:"employmentTypeId"`
			EmploymentType   string `json:"employmentType"`
			ProjectCount     int64  `json:"projectCount"`
		} `json:"jobHistories"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}
	if len(listResp.JobHistories) != 1 {
		t.Fatalf("expected one seeded job history, got %#v", listResp.JobHistories)
	}
	if listResp.JobHistories[0].Company != "株式会社A" ||
		listResp.JobHistories[0].DisplayName != "株式会社A" ||
		listResp.JobHistories[0].StartYear != 2023 ||
		listResp.JobHistories[0].StartMonth != 1 ||
		listResp.JobHistories[0].EndYear != nil ||
		listResp.JobHistories[0].EndMonth != nil ||
		listResp.JobHistories[0].EmploymentTypeID != "job_employment_type_full_time" ||
		listResp.JobHistories[0].EmploymentType != "正社員" ||
		listResp.JobHistories[0].ProjectCount != 5 {
		t.Fatalf("unexpected seeded job history: %#v", listResp.JobHistories[0])
	}

	createBody := []byte(`{
		"company":"株式会社C",
		"displayName":"表示用C",
		"startYear":2024,
		"startMonth":1,
		"endYear":null,
		"endMonth":null,
		"employmentTypeId":"job_employment_type_freelance"
	}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/job-histories", bytes.NewReader(createBody))
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status 201, got %d", createRec.Code)
	}

	var createResp struct {
		JobHistory struct {
			ID               string `json:"id"`
			Company          string `json:"company"`
			DisplayName      string `json:"displayName"`
			StartYear        int64  `json:"startYear"`
			StartMonth       int64  `json:"startMonth"`
			EndYear          *int64 `json:"endYear"`
			EndMonth         *int64 `json:"endMonth"`
			EmploymentTypeID string `json:"employmentTypeId"`
			EmploymentType   string `json:"employmentType"`
			ProjectCount     int64  `json:"projectCount"`
		} `json:"jobHistory"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}
	if createResp.JobHistory.ID == "" ||
		createResp.JobHistory.DisplayName != "表示用C" ||
		createResp.JobHistory.EndYear != nil ||
		createResp.JobHistory.EndMonth != nil ||
		createResp.JobHistory.EmploymentTypeID != "job_employment_type_freelance" ||
		createResp.JobHistory.EmploymentType != "業務委託" ||
		createResp.JobHistory.ProjectCount != 0 {
		t.Fatalf("unexpected created job history: %#v", createResp.JobHistory)
	}

	updateBody := []byte(`{
		"company":"株式会社C Updated",
		"displayName":"表示用C Updated",
		"startYear":2024,
		"startMonth":2,
		"endYear":2025,
		"endMonth":3,
		"employmentTypeId":"job_employment_type_contract"
	}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/job-histories/"+createResp.JobHistory.ID, bytes.NewReader(updateBody))
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected update status 200, got %d", updateRec.Code)
	}

	var updateResp struct {
		JobHistory struct {
			Company          string `json:"company"`
			DisplayName      string `json:"displayName"`
			StartYear        int64  `json:"startYear"`
			StartMonth       int64  `json:"startMonth"`
			EndYear          *int64 `json:"endYear"`
			EndMonth         *int64 `json:"endMonth"`
			EmploymentTypeID string `json:"employmentTypeId"`
			EmploymentType   string `json:"employmentType"`
		} `json:"jobHistory"`
	}
	if err := json.Unmarshal(updateRec.Body.Bytes(), &updateResp); err != nil {
		t.Fatalf("failed to decode update response: %v", err)
	}
	if updateResp.JobHistory.Company != "株式会社C Updated" ||
		updateResp.JobHistory.DisplayName != "表示用C Updated" ||
		updateResp.JobHistory.StartYear != 2024 ||
		updateResp.JobHistory.StartMonth != 2 ||
		updateResp.JobHistory.EndYear == nil ||
		*updateResp.JobHistory.EndYear != 2025 ||
		updateResp.JobHistory.EndMonth == nil ||
		*updateResp.JobHistory.EndMonth != 3 ||
		updateResp.JobHistory.EmploymentTypeID != "job_employment_type_contract" ||
		updateResp.JobHistory.EmploymentType != "契約社員" {
		t.Fatalf("unexpected updated job history: %#v", updateResp.JobHistory)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/job-histories/"+createResp.JobHistory.ID, nil)
	deleteRec := httptest.NewRecorder()
	router.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("expected delete status 204, got %d", deleteRec.Code)
	}
}

func TestProjectRoutes(t *testing.T) {
	repository := newTestProfileRepository()
	router := NewRouter(repository, repository, repository, repository)

	listReq := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	listRec := httptest.NewRecorder()
	router.ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("expected list status 200, got %d", listRec.Code)
	}

	var listResp struct {
		Projects []struct {
			ID           string   `json:"id"`
			Name         string   `json:"name"`
			Company      string   `json:"company"`
			StartYear    int64    `json:"startYear"`
			StartMonth   int64    `json:"startMonth"`
			EndYear      *int64   `json:"endYear"`
			EndMonth     *int64   `json:"endMonth"`
			Technologies []string `json:"technologies"`
			Phases       []string `json:"phases"`
			IsDraft      bool     `json:"isDraft"`
		} `json:"projects"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}
	if len(listResp.Projects) != 1 {
		t.Fatalf("expected one seeded project, got %#v", listResp.Projects)
	}
	if listResp.Projects[0].Name != "ECサイトリニューアル" ||
		listResp.Projects[0].Company != "株式会社A" ||
		listResp.Projects[0].StartYear != 2024 ||
		listResp.Projects[0].StartMonth != 1 ||
		listResp.Projects[0].EndYear != nil ||
		listResp.Projects[0].EndMonth != nil ||
		len(listResp.Projects[0].Technologies) != 2 ||
		listResp.Projects[0].Technologies[0] != "React" ||
		len(listResp.Projects[0].Phases) != 2 {
		t.Fatalf("unexpected seeded project: %#v", listResp.Projects[0])
	}

	createBody := []byte(`{
		"name":"新規案件",
		"company":"株式会社C",
		"startYear":2025,
		"startMonth":1,
		"endYear":null,
		"endMonth":null,
		"description":"新規開発",
		"role":"バックエンドエンジニア",
		"teamSize":"4名",
		"technologies":["Go","SQLite"],
		"phases":["設計","実装"],
		"achievements":"APIを構築",
		"isDraft":true
	}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/projects", bytes.NewReader(createBody))
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status 201, got %d", createRec.Code)
	}

	var createResp struct {
		Project struct {
			ID           string   `json:"id"`
			Name         string   `json:"name"`
			Technologies []string `json:"technologies"`
			IsDraft      bool     `json:"isDraft"`
		} `json:"project"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}
	if createResp.Project.ID == "" ||
		createResp.Project.Name != "新規案件" ||
		len(createResp.Project.Technologies) != 2 ||
		!createResp.Project.IsDraft {
		t.Fatalf("unexpected created project: %#v", createResp.Project)
	}

	updateBody := []byte(`{
		"name":"新規案件 Updated",
		"company":"株式会社C",
		"startYear":2025,
		"startMonth":1,
		"endYear":2025,
		"endMonth":3,
		"description":"新規開発",
		"role":"テックリード",
		"teamSize":"4名",
		"technologies":["Go","SQLite","React"],
		"phases":["設計","実装","テスト"],
		"achievements":"APIと画面を構築",
		"isDraft":false
	}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/projects/"+createResp.Project.ID, bytes.NewReader(updateBody))
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected update status 200, got %d", updateRec.Code)
	}

	var updateResp struct {
		Project struct {
			Name         string   `json:"name"`
			EndYear      *int64   `json:"endYear"`
			EndMonth     *int64   `json:"endMonth"`
			Technologies []string `json:"technologies"`
			Phases       []string `json:"phases"`
			IsDraft      bool     `json:"isDraft"`
		} `json:"project"`
	}
	if err := json.Unmarshal(updateRec.Body.Bytes(), &updateResp); err != nil {
		t.Fatalf("failed to decode update response: %v", err)
	}
	if updateResp.Project.Name != "新規案件 Updated" ||
		updateResp.Project.EndYear == nil ||
		*updateResp.Project.EndYear != 2025 ||
		updateResp.Project.EndMonth == nil ||
		*updateResp.Project.EndMonth != 3 ||
		len(updateResp.Project.Technologies) != 3 ||
		len(updateResp.Project.Phases) != 3 ||
		updateResp.Project.IsDraft {
		t.Fatalf("unexpected updated project: %#v", updateResp.Project)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/projects/"+createResp.Project.ID, nil)
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
	skillMasters      []domain.SkillMaster
	skills            []domain.Skill
	jobHistories      []domain.JobHistory
	employmentTypes   []domain.JobEmploymentType
	projects          []domain.Project
	nextSkillID       int
	nextJobHistoryID  int
	nextProjectID     int
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
		skillMasters: []domain.SkillMaster{
			{
				ID:           "skill_master_typescript",
				Name:         "TypeScript",
				CategoryID:   "skill_category_language",
				CategoryName: "言語",
				SortOrder:    1,
			},
		},
		skills: []domain.Skill{
			{
				ID:                 "skill_typescript",
				Name:               "TypeScript",
				CategoryID:         "skill_category_language",
				CategoryName:       "言語",
				Experience:         3,
				ProficiencyLevelID: "skill_proficiency_advanced",
				ProficiencyName:    "上級",
				SortOrder:          1,
			},
		},
		jobHistories: []domain.JobHistory{
			{
				ID:               "job_history_company_a",
				Company:          "株式会社A",
				DisplayName:      "株式会社A",
				StartYear:        2023,
				StartMonth:       1,
				EndYear:          nil,
				EndMonth:         nil,
				EmploymentTypeID: "job_employment_type_full_time",
				EmploymentType:   "正社員",
				ProjectCount:     5,
				SortOrder:        1,
			},
		},
		employmentTypes: []domain.JobEmploymentType{
			{ID: "job_employment_type_full_time", Name: "正社員", SortOrder: 1},
			{ID: "job_employment_type_contract", Name: "契約社員", SortOrder: 2},
			{ID: "job_employment_type_freelance", Name: "業務委託", SortOrder: 3},
		},
		projects: []domain.Project{
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
				Technologies: []string{"React", "TypeScript"},
				Phases:       []string{"設計", "実装"},
				Achievements: "ページ表示速度を改善",
				IsDraft:      false,
				SortOrder:    1,
			},
		},
		nextSkillID:      1,
		nextJobHistoryID: 1,
		nextProjectID:    1,
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
	skillMasters := make([]domain.SkillMaster, len(r.skillMasters))
	copy(skillMasters, r.skillMasters)

	return &domain.SkillOptions{
		Categories:        categories,
		ProficiencyLevels: proficiencyLevels,
		SkillMasters:      skillMasters,
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

func (r *testProfileRepository) CreateSkillMaster(_ context.Context, input domain.SkillMasterInput) (*domain.SkillMaster, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	category, ok := r.findCategory(input.CategoryID)
	if !ok {
		return nil, sql.ErrNoRows
	}

	skillMaster := domain.SkillMaster{
		ID:           input.ID,
		Name:         input.Name,
		CategoryID:   input.CategoryID,
		CategoryName: category.Name,
		SortOrder:    int64(len(r.skillMasters) + 1),
	}
	r.skillMasters = append(r.skillMasters, skillMaster)

	return &skillMaster, nil
}

func (r *testProfileRepository) UpdateSkillMaster(_ context.Context, id string, input domain.SkillMasterInput) (*domain.SkillMaster, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	category, ok := r.findCategory(input.CategoryID)
	if !ok {
		return nil, sql.ErrNoRows
	}

	for index, skillMaster := range r.skillMasters {
		if skillMaster.ID == id {
			r.skillMasters[index].Name = input.Name
			r.skillMasters[index].CategoryID = input.CategoryID
			r.skillMasters[index].CategoryName = category.Name
			result := r.skillMasters[index]
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

func (r *testProfileRepository) ListJobHistories(context.Context) ([]domain.JobHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobHistories := make([]domain.JobHistory, len(r.jobHistories))
	copy(jobHistories, r.jobHistories)

	return jobHistories, nil
}

func (r *testProfileRepository) CreateJobHistory(_ context.Context, input domain.JobHistoryInput) (*domain.JobHistory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	employmentType, ok := r.findEmploymentType(input.EmploymentTypeID)
	if !ok {
		return nil, sql.ErrNoRows
	}

	r.nextJobHistoryID++
	jobHistory := domain.JobHistory{
		ID:               "job_history_test_new",
		Company:          input.Company,
		DisplayName:      input.DisplayName,
		StartYear:        input.StartYear,
		StartMonth:       input.StartMonth,
		EndYear:          input.EndYear,
		EndMonth:         input.EndMonth,
		EmploymentTypeID: input.EmploymentTypeID,
		EmploymentType:   employmentType.Name,
		ProjectCount:     0,
		SortOrder:        int64(len(r.jobHistories) + 1),
	}
	r.jobHistories = append(r.jobHistories, jobHistory)

	return &jobHistory, nil
}

func (r *testProfileRepository) UpdateJobHistory(_ context.Context, id string, input domain.JobHistoryInput) (*domain.JobHistory, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	employmentType, ok := r.findEmploymentType(input.EmploymentTypeID)
	if !ok {
		return nil, sql.ErrNoRows
	}

	for index, jobHistory := range r.jobHistories {
		if jobHistory.ID == id {
			r.jobHistories[index].Company = input.Company
			r.jobHistories[index].DisplayName = input.DisplayName
			r.jobHistories[index].StartYear = input.StartYear
			r.jobHistories[index].StartMonth = input.StartMonth
			r.jobHistories[index].EndYear = input.EndYear
			r.jobHistories[index].EndMonth = input.EndMonth
			r.jobHistories[index].EmploymentTypeID = input.EmploymentTypeID
			r.jobHistories[index].EmploymentType = employmentType.Name
			result := r.jobHistories[index]
			return &result, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *testProfileRepository) DeleteJobHistory(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, jobHistory := range r.jobHistories {
		if jobHistory.ID == id {
			r.jobHistories = append(r.jobHistories[:index], r.jobHistories[index+1:]...)
			return nil
		}
	}

	return sql.ErrNoRows
}

func (r *testProfileRepository) ListJobHistoryOptions(context.Context) (*domain.JobHistoryOptions, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	employmentTypes := make([]domain.JobEmploymentType, len(r.employmentTypes))
	copy(employmentTypes, r.employmentTypes)

	return &domain.JobHistoryOptions{EmploymentTypes: employmentTypes}, nil
}

func (r *testProfileRepository) CreateJobEmploymentType(_ context.Context, input domain.JobEmploymentTypeInput) (*domain.JobEmploymentType, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	employmentType := domain.JobEmploymentType{
		ID:        input.ID,
		Name:      input.Name,
		SortOrder: int64(len(r.employmentTypes) + 1),
	}
	r.employmentTypes = append(r.employmentTypes, employmentType)

	return &employmentType, nil
}

func (r *testProfileRepository) UpdateJobEmploymentType(_ context.Context, id string, input domain.JobEmploymentTypeInput) (*domain.JobEmploymentType, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, employmentType := range r.employmentTypes {
		if employmentType.ID == id {
			r.employmentTypes[index].Name = input.Name
			result := r.employmentTypes[index]
			return &result, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *testProfileRepository) ListProjects(context.Context) ([]domain.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	projects := make([]domain.Project, len(r.projects))
	copy(projects, r.projects)

	return projects, nil
}

func (r *testProfileRepository) CreateProject(_ context.Context, input domain.ProjectInput) (*domain.Project, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextProjectID++
	project := domain.Project{
		ID:           "project_test_new",
		Name:         input.Name,
		Company:      input.Company,
		StartYear:    input.StartYear,
		StartMonth:   input.StartMonth,
		EndYear:      input.EndYear,
		EndMonth:     input.EndMonth,
		Description:  input.Description,
		Role:         input.Role,
		TeamSize:     input.TeamSize,
		Technologies: input.Technologies,
		Phases:       input.Phases,
		Achievements: input.Achievements,
		IsDraft:      input.IsDraft,
		SortOrder:    int64(len(r.projects) + 1),
	}
	r.projects = append(r.projects, project)

	return &project, nil
}

func (r *testProfileRepository) UpdateProject(_ context.Context, id string, input domain.ProjectInput) (*domain.Project, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, project := range r.projects {
		if project.ID == id {
			r.projects[index].Name = input.Name
			r.projects[index].Company = input.Company
			r.projects[index].StartYear = input.StartYear
			r.projects[index].StartMonth = input.StartMonth
			r.projects[index].EndYear = input.EndYear
			r.projects[index].EndMonth = input.EndMonth
			r.projects[index].Description = input.Description
			r.projects[index].Role = input.Role
			r.projects[index].TeamSize = input.TeamSize
			r.projects[index].Technologies = input.Technologies
			r.projects[index].Phases = input.Phases
			r.projects[index].Achievements = input.Achievements
			r.projects[index].IsDraft = input.IsDraft
			result := r.projects[index]
			return &result, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *testProfileRepository) DeleteProject(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, project := range r.projects {
		if project.ID == id {
			r.projects = append(r.projects[:index], r.projects[index+1:]...)
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

func (r *testProfileRepository) findEmploymentType(id string) (domain.JobEmploymentType, bool) {
	for _, employmentType := range r.employmentTypes {
		if employmentType.ID == id {
			return employmentType, true
		}
	}

	return domain.JobEmploymentType{}, false
}
