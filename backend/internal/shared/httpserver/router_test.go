package httpserver

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
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
			FamilyName string `json:"familyName"`
			GivenName  string `json:"givenName"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(getRec.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("failed to decode get response: %v", err)
	}

	if getResp.Profile.FamilyName != "Initial" || getResp.Profile.GivenName != "User" {
		t.Fatalf("expected initial name fields, got %+v", getResp.Profile)
	}

	body := []byte(`{
		"familyName":"Sky",
		"givenName":"Sample",
		"location":"Tokyo",
		"email":"me@example.com",
		"bio":"Backend engineer",
		"occupation":"Software Engineer",
		"employmentType":"Freelance",
		"workStyle":"Full remote",
		"links":[
			{"linkMasterId":"1","key":"github","url":"https://github.com/sky0621"},
			{"linkMasterId":"2","key":"zenn","url":"https://zenn.dev/sky0621"},
			{"linkMasterId":"3","key":"qiita","url":"https://qiita.com/sky0621"},
			{"linkMasterId":"4","key":"website","url":"https://example.com"}
		],
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
			FamilyName     string `json:"familyName"`
			GivenName      string `json:"givenName"`
			Location       string `json:"location"`
			Email          string `json:"email"`
			Bio            string `json:"bio"`
			Occupation     string `json:"occupation"`
			EmploymentType string `json:"employmentType"`
			WorkStyle      string `json:"workStyle"`
			Links          []struct {
				LinkMasterID string `json:"linkMasterId"`
				Key          string `json:"key"`
				URL          string `json:"url"`
			} `json:"links"`
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

	if putResp.Profile.FamilyName != "Sky" || putResp.Profile.GivenName != "Sample" {
		t.Fatalf("expected updated name fields, got %+v", putResp.Profile)
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
	if len(putResp.Profile.Links) != 4 || putResp.Profile.Links[0].URL != "https://github.com/sky0621" {
		t.Fatalf("expected updated profile links, got %#v", putResp.Profile.Links)
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
			FamilyName     string `json:"familyName"`
			GivenName      string `json:"givenName"`
			Location       string `json:"location"`
			Email          string `json:"email"`
			Bio            string `json:"bio"`
			Occupation     string `json:"occupation"`
			EmploymentType string `json:"employmentType"`
			WorkStyle      string `json:"workStyle"`
			Links          []struct {
				LinkMasterID string `json:"linkMasterId"`
				Key          string `json:"key"`
				URL          string `json:"url"`
			} `json:"links"`
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

	if getUpdatedResp.Profile.FamilyName != "Sky" || getUpdatedResp.Profile.GivenName != "Sample" {
		t.Fatalf("expected persisted name fields, got %+v", getUpdatedResp.Profile)
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
	if len(getUpdatedResp.Profile.Links) != 4 || getUpdatedResp.Profile.Links[0].URL != "https://github.com/sky0621" {
		t.Fatalf("expected persisted profile links, got %#v", getUpdatedResp.Profile.Links)
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
		"familyName":"Sky",
		"givenName":"Sample",
		"location":"Tokyo",
		"email":"me@example.com",
		"bio":"Backend engineer",
		"occupation":"Software Engineer",
		"employmentType":"Freelance",
		"workStyle":"Full remote",
		"links":[
			{"linkMasterId":"1","key":"github","url":"https://github.com/sky0621"}
		],
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
	if resp.Categories[0].ID != "1" || resp.Categories[0].Name != "言語" {
		t.Fatalf("unexpected first category: %#v", resp.Categories[0])
	}
	if resp.Categories[0].Icon != "code" {
		t.Fatalf("unexpected first category icon: %#v", resp.Categories[0])
	}
	if len(resp.SkillMasters) != 2 {
		t.Fatalf("expected two skill masters, got %#v", resp.SkillMasters)
	}
	if resp.SkillMasters[0].ID != "4" ||
		resp.SkillMasters[0].Name != "React" ||
		resp.SkillMasters[0].CategoryID != "2" ||
		resp.SkillMasters[0].Category != "フレームワーク" {
		t.Fatalf("unexpected first skill master: %#v", resp.SkillMasters[0])
	}
}

func TestSkillCategoryMutationRoutes(t *testing.T) {
	repository := newTestProfileRepository()
	router := NewRouter(repository, repository, repository, repository)

	createBody := []byte(`{
		"id":"7",
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
	if createResp.Category.ID != "7" || createResp.Category.Name != "バックエンド" {
		t.Fatalf("unexpected created category: %#v", createResp.Category)
	}

	updateBody := []byte(`{
		"name":"バックエンド・API",
		"icon":"database"
	}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/skills/categories/7", bytes.NewReader(updateBody))
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
	if updateResp.Category.ID != "7" ||
		updateResp.Category.Name != "バックエンド・API" ||
		updateResp.Category.Icon != "database" {
		t.Fatalf("unexpected updated category: %#v", updateResp.Category)
	}
}

func TestSkillMasterMutationRoutes(t *testing.T) {
	repository := newTestProfileRepository()
	router := NewRouter(repository, repository, repository, repository)

	createBody := []byte(`{
		"name":"Kotlin",
		"categoryId":"1"
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
	if createResp.SkillMaster.ID != "5" ||
		createResp.SkillMaster.Name != "Kotlin" ||
		createResp.SkillMaster.Category != "言語" {
		t.Fatalf("unexpected created skill master: %#v", createResp.SkillMaster)
	}

	updateBody := []byte(`{
		"name":"Kotlin/JVM",
		"categoryId":"2"
	}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/skills/masters/5", bytes.NewReader(updateBody))
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
	if updateResp.SkillMaster.ID != "5" ||
		updateResp.SkillMaster.Name != "Kotlin/JVM" ||
		updateResp.SkillMaster.CategoryID != "2" ||
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
			ID            string `json:"id"`
			SkillMasterID string `json:"skillMasterId"`
			Name          string `json:"name"`
			CategoryID    string `json:"categoryId"`
			Category      string `json:"category"`
			Experience    int64  `json:"experience"`
		} `json:"skills"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}
	if len(listResp.Skills) != 1 {
		t.Fatalf("expected one seeded skill, got %#v", listResp.Skills)
	}
	if listResp.Skills[0].SkillMasterID != "1" ||
		listResp.Skills[0].Category != "言語" {
		t.Fatalf("expected skill master names, got %#v", listResp.Skills[0])
	}

	duplicateBody := []byte(`{
		"skillMasterId":"1",
		"experience":5
	}`)
	duplicateReq := httptest.NewRequest(http.MethodPost, "/api/skills", bytes.NewReader(duplicateBody))
	duplicateRec := httptest.NewRecorder()
	router.ServeHTTP(duplicateRec, duplicateReq)

	if duplicateRec.Code != http.StatusConflict {
		t.Fatalf("expected duplicate create status 409, got %d", duplicateRec.Code)
	}

	createBody := []byte(`{
		"skillMasterId":"4",
		"experience":3
	}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/skills", bytes.NewReader(createBody))
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status 201, got %d", createRec.Code)
	}

	var createResp struct {
		Skill struct {
			ID            string `json:"id"`
			SkillMasterID string `json:"skillMasterId"`
			Name          string `json:"name"`
			CategoryID    string `json:"categoryId"`
			Category      string `json:"category"`
			Experience    int64  `json:"experience"`
		} `json:"skill"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}
	if createResp.Skill.ID == "" ||
		createResp.Skill.SkillMasterID != "4" ||
		createResp.Skill.Category != "フレームワーク" {
		t.Fatalf("unexpected created skill: %#v", createResp.Skill)
	}

	updateBody := []byte(`{
		"skillMasterId":"4",
		"experience":4
	}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/skills/"+createResp.Skill.ID, bytes.NewReader(updateBody))
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected update status 200, got %d", updateRec.Code)
	}

	var updateResp struct {
		Skill struct {
			ID         string `json:"id"`
			Experience int64  `json:"experience"`
		} `json:"skill"`
	}
	if err := json.Unmarshal(updateRec.Body.Bytes(), &updateResp); err != nil {
		t.Fatalf("failed to decode update response: %v", err)
	}
	if updateResp.Skill.Experience != 4 {
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
			CompanyID        string `json:"companyId"`
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
	if listResp.JobHistories[0].CompanyID != "1" ||
		listResp.JobHistories[0].Company != "株式会社A" ||
		listResp.JobHistories[0].DisplayName != "株式会社A" ||
		listResp.JobHistories[0].StartYear != 2023 ||
		listResp.JobHistories[0].StartMonth != 1 ||
		listResp.JobHistories[0].EndYear != nil ||
		listResp.JobHistories[0].EndMonth != nil ||
		listResp.JobHistories[0].EmploymentTypeID != "1" ||
		listResp.JobHistories[0].EmploymentType != "正社員" ||
		listResp.JobHistories[0].ProjectCount != 5 {
		t.Fatalf("unexpected seeded job history: %#v", listResp.JobHistories[0])
	}

	optionsReq := httptest.NewRequest(http.MethodGet, "/api/job-histories/options", nil)
	optionsRec := httptest.NewRecorder()
	router.ServeHTTP(optionsRec, optionsReq)

	if optionsRec.Code != http.StatusOK {
		t.Fatalf("expected options status 200, got %d", optionsRec.Code)
	}

	var optionsResp struct {
		Companies []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"companies"`
	}
	if err := json.Unmarshal(optionsRec.Body.Bytes(), &optionsResp); err != nil {
		t.Fatalf("failed to decode options response: %v", err)
	}
	if len(optionsResp.Companies) != 1 ||
		optionsResp.Companies[0].Name != "株式会社A" ||
		optionsResp.Companies[0].URL != "https://example.com/a" {
		t.Fatalf("unexpected company options: %#v", optionsResp.Companies)
	}

	createBody := []byte(`{
		"companyId":"1",
		"displayName":"表示用C",
		"startYear":2024,
		"startMonth":1,
		"endYear":null,
		"endMonth":null,
		"employmentTypeId":"3"
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
			CompanyID        string `json:"companyId"`
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
		createResp.JobHistory.CompanyID != "1" ||
		createResp.JobHistory.Company != "株式会社A" ||
		createResp.JobHistory.DisplayName != "表示用C" ||
		createResp.JobHistory.EndYear != nil ||
		createResp.JobHistory.EndMonth != nil ||
		createResp.JobHistory.EmploymentTypeID != "3" ||
		createResp.JobHistory.EmploymentType != "業務委託" ||
		createResp.JobHistory.ProjectCount != 0 {
		t.Fatalf("unexpected created job history: %#v", createResp.JobHistory)
	}

	updateBody := []byte(`{
		"companyId":"1",
		"displayName":"表示用C Updated",
		"startYear":2024,
		"startMonth":2,
		"endYear":2025,
		"endMonth":3,
		"employmentTypeId":"2"
	}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/job-histories/"+createResp.JobHistory.ID, bytes.NewReader(updateBody))
	updateRec := httptest.NewRecorder()
	router.ServeHTTP(updateRec, updateReq)

	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected update status 200, got %d", updateRec.Code)
	}

	var updateResp struct {
		JobHistory struct {
			CompanyID        string `json:"companyId"`
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
	if updateResp.JobHistory.CompanyID != "1" ||
		updateResp.JobHistory.Company != "株式会社A" ||
		updateResp.JobHistory.DisplayName != "表示用C Updated" ||
		updateResp.JobHistory.StartYear != 2024 ||
		updateResp.JobHistory.StartMonth != 2 ||
		updateResp.JobHistory.EndYear == nil ||
		*updateResp.JobHistory.EndYear != 2025 ||
		updateResp.JobHistory.EndMonth == nil ||
		*updateResp.JobHistory.EndMonth != 3 ||
		updateResp.JobHistory.EmploymentTypeID != "2" ||
		updateResp.JobHistory.EmploymentType != "契約社員" {
		t.Fatalf("unexpected updated job history: %#v", updateResp.JobHistory)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/job-histories/"+createResp.JobHistory.ID, nil)
	deleteRec := httptest.NewRecorder()
	router.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("expected delete status 204, got %d", deleteRec.Code)
	}

	companyCreateBody := []byte(`{"name":"株式会社D","url":"https://example.com/d"}`)
	companyCreateReq := httptest.NewRequest(http.MethodPost, "/api/job-companies", bytes.NewReader(companyCreateBody))
	companyCreateRec := httptest.NewRecorder()
	router.ServeHTTP(companyCreateRec, companyCreateReq)

	if companyCreateRec.Code != http.StatusCreated {
		t.Fatalf("expected company create status 201, got %d", companyCreateRec.Code)
	}

	var companyCreateResp struct {
		Company struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"company"`
	}
	if err := json.Unmarshal(companyCreateRec.Body.Bytes(), &companyCreateResp); err != nil {
		t.Fatalf("failed to decode company create response: %v", err)
	}
	if companyCreateResp.Company.ID == "" ||
		companyCreateResp.Company.Name != "株式会社D" ||
		companyCreateResp.Company.URL != "https://example.com/d" {
		t.Fatalf("unexpected created company: %#v", companyCreateResp.Company)
	}

	companyUpdateBody := []byte(`{"name":"株式会社D Updated","url":"https://example.com/d-updated"}`)
	companyUpdateReq := httptest.NewRequest(http.MethodPut, "/api/job-companies/"+companyCreateResp.Company.ID, bytes.NewReader(companyUpdateBody))
	companyUpdateRec := httptest.NewRecorder()
	router.ServeHTTP(companyUpdateRec, companyUpdateReq)

	if companyUpdateRec.Code != http.StatusOK {
		t.Fatalf("expected company update status 200, got %d", companyUpdateRec.Code)
	}
}

func TestProjectRoutes(t *testing.T) {
	repository := newTestProfileRepository()
	router := NewRouter(repository, repository, repository, repository)

	optionsReq := httptest.NewRequest(http.MethodGet, "/api/projects/options", nil)
	optionsRec := httptest.NewRecorder()
	router.ServeHTTP(optionsRec, optionsReq)

	if optionsRec.Code != http.StatusOK {
		t.Fatalf("expected options status 200, got %d", optionsRec.Code)
	}

	var optionsResp struct {
		Phases []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			SortOrder int64  `json:"sortOrder"`
		} `json:"phases"`
	}
	if err := json.Unmarshal(optionsRec.Body.Bytes(), &optionsResp); err != nil {
		t.Fatalf("failed to decode project options response: %v", err)
	}
	if len(optionsResp.Phases) != 2 || optionsResp.Phases[0].Name != "設計" {
		t.Fatalf("unexpected project phases: %#v", optionsResp.Phases)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	listRec := httptest.NewRecorder()
	router.ServeHTTP(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("expected list status 200, got %d", listRec.Code)
	}

	var listResp struct {
		Projects []struct {
			ID            string   `json:"id"`
			Name          string   `json:"name"`
			CompanyID     string   `json:"companyId"`
			Company       string   `json:"company"`
			TechnologyIDs []string `json:"technologyIds"`
			StartYear     int64    `json:"startYear"`
			StartMonth    int64    `json:"startMonth"`
			EndYear       *int64   `json:"endYear"`
			EndMonth      *int64   `json:"endMonth"`
			Technologies  []string `json:"technologies"`
			Phases        []string `json:"phases"`
			IsDraft       bool     `json:"isDraft"`
		} `json:"projects"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("failed to decode list response: %v", err)
	}
	if len(listResp.Projects) != 1 {
		t.Fatalf("expected one seeded project, got %#v", listResp.Projects)
	}
	if listResp.Projects[0].Name != "ECサイトリニューアル" ||
		listResp.Projects[0].CompanyID != "1" ||
		listResp.Projects[0].Company != "株式会社A" ||
		listResp.Projects[0].StartYear != 2024 ||
		listResp.Projects[0].StartMonth != 1 ||
		listResp.Projects[0].EndYear != nil ||
		listResp.Projects[0].EndMonth != nil ||
		len(listResp.Projects[0].Technologies) != 2 ||
		len(listResp.Projects[0].TechnologyIDs) != 2 ||
		listResp.Projects[0].TechnologyIDs[0] != "4" ||
		listResp.Projects[0].Technologies[0] != "React" ||
		len(listResp.Projects[0].Phases) != 2 {
		t.Fatalf("unexpected seeded project: %#v", listResp.Projects[0])
	}

	createBody := []byte(`{
		"name":"新規案件",
		"companyId":"1",
		"startYear":2025,
		"startMonth":1,
		"endYear":null,
		"endMonth":null,
		"description":"新規開発",
		"role":"バックエンドエンジニア",
		"teamSize":"4名",
		"technologyIds":["1","4"],
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
			ID            string   `json:"id"`
			Name          string   `json:"name"`
			CompanyID     string   `json:"companyId"`
			Company       string   `json:"company"`
			TechnologyIDs []string `json:"technologyIds"`
			Technologies  []string `json:"technologies"`
			IsDraft       bool     `json:"isDraft"`
		} `json:"project"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}
	if createResp.Project.ID == "" ||
		createResp.Project.Name != "新規案件" ||
		createResp.Project.CompanyID != "1" ||
		createResp.Project.Company != "株式会社A" ||
		len(createResp.Project.TechnologyIDs) != 2 ||
		len(createResp.Project.Technologies) != 2 ||
		!createResp.Project.IsDraft {
		t.Fatalf("unexpected created project: %#v", createResp.Project)
	}

	updateBody := []byte(`{
		"name":"新規案件 Updated",
		"companyId":"1",
		"startYear":2025,
		"startMonth":1,
		"endYear":2025,
		"endMonth":3,
		"description":"新規開発",
		"role":"テックリード",
		"teamSize":"4名",
		"technologyIds":["1"],
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
			Name          string   `json:"name"`
			EndYear       *int64   `json:"endYear"`
			EndMonth      *int64   `json:"endMonth"`
			TechnologyIDs []string `json:"technologyIds"`
			Technologies  []string `json:"technologies"`
			Phases        []string `json:"phases"`
			IsDraft       bool     `json:"isDraft"`
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
		len(updateResp.Project.TechnologyIDs) != 1 ||
		len(updateResp.Project.Technologies) != 1 ||
		len(updateResp.Project.Phases) != 3 ||
		updateResp.Project.IsDraft {
		t.Fatalf("unexpected updated project: %#v", updateResp.Project)
	}

	phaseCreateBody := []byte(`{"name":"テスト","sortOrder":3}`)
	phaseCreateReq := httptest.NewRequest(http.MethodPost, "/api/project-phases", bytes.NewReader(phaseCreateBody))
	phaseCreateRec := httptest.NewRecorder()
	router.ServeHTTP(phaseCreateRec, phaseCreateReq)

	if phaseCreateRec.Code != http.StatusCreated {
		t.Fatalf("expected phase create status 201, got %d", phaseCreateRec.Code)
	}

	var phaseCreateResp struct {
		Phase struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			SortOrder int64  `json:"sortOrder"`
		} `json:"phase"`
	}
	if err := json.Unmarshal(phaseCreateRec.Body.Bytes(), &phaseCreateResp); err != nil {
		t.Fatalf("failed to decode phase create response: %v", err)
	}
	if phaseCreateResp.Phase.ID == "" || phaseCreateResp.Phase.Name != "テスト" {
		t.Fatalf("unexpected created phase: %#v", phaseCreateResp.Phase)
	}

	phaseUpdateBody := []byte(`{"name":"テスト・検証","sortOrder":4}`)
	phaseUpdateReq := httptest.NewRequest(http.MethodPut, "/api/project-phases/"+phaseCreateResp.Phase.ID, bytes.NewReader(phaseUpdateBody))
	phaseUpdateRec := httptest.NewRecorder()
	router.ServeHTTP(phaseUpdateRec, phaseUpdateReq)

	if phaseUpdateRec.Code != http.StatusOK {
		t.Fatalf("expected phase update status 200, got %d", phaseUpdateRec.Code)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/projects/"+createResp.Project.ID, nil)
	deleteRec := httptest.NewRecorder()
	router.ServeHTTP(deleteRec, deleteReq)

	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("expected delete status 204, got %d", deleteRec.Code)
	}
}

type testProfileRepository struct {
	mu                 sync.RWMutex
	profile            *domain.Profile
	categories         []domain.SkillOption
	skillMasters       []domain.SkillMaster
	profileLinkMasters []domain.ProfileLinkMaster
	skills             []domain.Skill
	jobHistories       []domain.JobHistory
	employmentTypes    []domain.JobEmploymentType
	jobCompanies       []domain.JobCompany
	projects           []domain.Project
	projectPhases      []domain.ProjectPhase
	nextSkillID        int
	nextJobHistoryID   int
	nextProjectID      int
}

func newTestProfileRepository() *testProfileRepository {
	now := time.Now().UTC()

	return &testProfileRepository{
		profile: &domain.Profile{
			ID:                 "1",
			FamilyName:         "Initial",
			GivenName:          "User",
			VisibilitySettings: map[string]any{"email": false},
			CreatedAt:          now,
			UpdatedAt:          now,
		},
		categories: []domain.SkillOption{
			{ID: "1", Name: "言語", Icon: "code", SortOrder: 1},
			{ID: "2", Name: "フレームワーク", Icon: "code", SortOrder: 2},
			{ID: "3", Name: "データベース", Icon: "database", SortOrder: 3},
			{ID: "4", Name: "インフラ", Icon: "cloud", SortOrder: 4},
			{ID: "5", Name: "ツール", Icon: "wrench", SortOrder: 5},
			{ID: "6", Name: "その他", Icon: "wrench", SortOrder: 6},
		},
		skillMasters: []domain.SkillMaster{
			{
				ID:           "1",
				Name:         "TypeScript",
				CategoryID:   "1",
				CategoryName: "言語",
			},
			{
				ID:           "4",
				Name:         "React",
				CategoryID:   "2",
				CategoryName: "フレームワーク",
			},
		},
		profileLinkMasters: []domain.ProfileLinkMaster{
			{ID: "1", Key: "github", Name: "GitHub", Icon: "github", Placeholder: "https://github.com/username", SortOrder: 1},
			{ID: "2", Key: "zenn", Name: "Zenn", Icon: "book-open", Placeholder: "https://zenn.dev/username", SortOrder: 2},
			{ID: "3", Key: "qiita", Name: "Qiita", Icon: "book-open", Placeholder: "https://qiita.com/username", SortOrder: 3},
			{ID: "4", Key: "website", Name: "個人サイト", Icon: "globe", Placeholder: "https://example.com", SortOrder: 4},
		},
		skills: []domain.Skill{
			{
				ID:            "skill_typescript",
				SkillMasterID: "1",
				Name:          "TypeScript",
				CategoryID:    "1",
				CategoryName:  "言語",
				Experience:    3,
				SortOrder:     1,
			},
		},
		jobHistories: []domain.JobHistory{
			{
				ID:               "1",
				CompanyID:        "1",
				Company:          "株式会社A",
				DisplayName:      "株式会社A",
				StartYear:        2023,
				StartMonth:       1,
				EndYear:          nil,
				EndMonth:         nil,
				EmploymentTypeID: "1",
				EmploymentType:   "正社員",
				ProjectCount:     5,
			},
		},
		employmentTypes: []domain.JobEmploymentType{
			{ID: "1", Name: "正社員", SortOrder: 1},
			{ID: "2", Name: "契約社員", SortOrder: 2},
			{ID: "3", Name: "業務委託", SortOrder: 3},
		},
		jobCompanies: []domain.JobCompany{
			{ID: "1", Name: "株式会社A", URL: "https://example.com/a"},
		},
		projects: []domain.Project{
			{
				ID:            "1",
				Name:          "ECサイトリニューアル",
				CompanyID:     "1",
				Company:       "株式会社A",
				StartYear:     2024,
				StartMonth:    1,
				EndYear:       nil,
				EndMonth:      nil,
				Description:   "大手ECサイトのフロントエンド刷新プロジェクト",
				Role:          "フロントエンドエンジニア",
				TeamSize:      "8名",
				TechnologyIDs: []string{"4", "1"},
				Technologies:  []string{"React", "TypeScript"},
				Phases:        []string{"設計", "実装"},
				Achievements:  "ページ表示速度を改善",
				IsDraft:       false,
				SortOrder:     1,
			},
		},
		projectPhases: []domain.ProjectPhase{
			{ID: "1", Name: "設計", SortOrder: 1},
			{ID: "2", Name: "実装", SortOrder: 2},
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

func (r *testProfileRepository) ListProfileLinkMasters(context.Context) ([]domain.ProfileLinkMaster, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	masters := make([]domain.ProfileLinkMaster, len(r.profileLinkMasters))
	copy(masters, r.profileLinkMasters)

	return masters, nil
}

func (r *testProfileRepository) CreateProfileLinkMaster(_ context.Context, input domain.ProfileLinkMasterInput) (*domain.ProfileLinkMaster, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	master := domain.ProfileLinkMaster{
		ID:          input.ID,
		Key:         input.Key,
		Name:        input.Name,
		Icon:        input.Icon,
		Placeholder: input.Placeholder,
		SortOrder:   input.SortOrder,
	}
	if master.SortOrder == 0 {
		master.SortOrder = int64(len(r.profileLinkMasters) + 1)
	}
	r.profileLinkMasters = append(r.profileLinkMasters, master)

	return &master, nil
}

func (r *testProfileRepository) UpdateProfileLinkMaster(_ context.Context, id string, input domain.ProfileLinkMasterInput) (*domain.ProfileLinkMaster, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, master := range r.profileLinkMasters {
		if master.ID == id {
			r.profileLinkMasters[index].Key = input.Key
			r.profileLinkMasters[index].Name = input.Name
			r.profileLinkMasters[index].Icon = input.Icon
			r.profileLinkMasters[index].Placeholder = input.Placeholder
			if input.SortOrder != 0 {
				r.profileLinkMasters[index].SortOrder = input.SortOrder
			}
			result := r.profileLinkMasters[index]
			return &result, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *testProfileRepository) ListSkillOptions(context.Context) (*domain.SkillOptions, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make([]domain.SkillOption, len(r.categories))
	copy(categories, r.categories)
	skillMasters := make([]domain.SkillMaster, len(r.skillMasters))
	copy(skillMasters, r.skillMasters)
	sort.Slice(skillMasters, func(i, j int) bool {
		return skillMasters[i].Name < skillMasters[j].Name
	})

	return &domain.SkillOptions{
		Categories:   categories,
		SkillMasters: skillMasters,
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
			if input.SortOrder != 0 {
				r.categories[index].SortOrder = input.SortOrder
			}
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
		ID:           strconv.Itoa(nextSkillMasterID(r.skillMasters)),
		Name:         input.Name,
		CategoryID:   input.CategoryID,
		CategoryName: category.Name,
	}
	r.skillMasters = append(r.skillMasters, skillMaster)

	return &skillMaster, nil
}

func nextSkillMasterID(skillMasters []domain.SkillMaster) int {
	maxID := 0
	for _, skillMaster := range skillMasters {
		id, err := strconv.Atoi(skillMaster.ID)
		if err != nil {
			continue
		}
		if id > maxID {
			maxID = id
		}
	}

	return maxID + 1
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

	skillMaster, ok := r.findSkillMaster(input.SkillMasterID)
	if !ok {
		return nil, sql.ErrNoRows
	}
	for _, skill := range r.skills {
		if skill.SkillMasterID == input.SkillMasterID {
			return nil, errors.New("UNIQUE constraint failed: skills.skill_master_id")
		}
	}

	r.nextSkillID++
	skill := domain.Skill{
		ID:            "skill_test_new",
		SkillMasterID: skillMaster.ID,
		Name:          skillMaster.Name,
		CategoryID:    skillMaster.CategoryID,
		CategoryName:  skillMaster.CategoryName,
		Experience:    input.Experience,
		SortOrder:     int64(len(r.skills) + 1),
	}
	r.skills = append(r.skills, skill)

	return &skill, nil
}

func (r *testProfileRepository) UpdateSkill(_ context.Context, id string, input domain.SkillInput) (*domain.Skill, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	skillMaster, ok := r.findSkillMaster(input.SkillMasterID)
	if !ok {
		return nil, sql.ErrNoRows
	}
	for _, skill := range r.skills {
		if skill.ID != id && skill.SkillMasterID == input.SkillMasterID {
			return nil, errors.New("UNIQUE constraint failed: skills.skill_master_id")
		}
	}

	for index, skill := range r.skills {
		if skill.ID == id {
			r.skills[index].SkillMasterID = skillMaster.ID
			r.skills[index].Name = skillMaster.Name
			r.skills[index].CategoryID = skillMaster.CategoryID
			r.skills[index].CategoryName = skillMaster.CategoryName
			r.skills[index].Experience = input.Experience
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
	company, ok := r.findJobCompany(input.CompanyID)
	if input.CompanyID != "" && !ok {
		return nil, sql.ErrNoRows
	}

	r.nextJobHistoryID++
	jobHistory := domain.JobHistory{
		ID:               "job_history_test_new",
		CompanyID:        input.CompanyID,
		Company:          company.Name,
		DisplayName:      input.DisplayName,
		StartYear:        input.StartYear,
		StartMonth:       input.StartMonth,
		EndYear:          input.EndYear,
		EndMonth:         input.EndMonth,
		EmploymentTypeID: input.EmploymentTypeID,
		EmploymentType:   employmentType.Name,
		ProjectCount:     0,
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
	company, ok := r.findJobCompany(input.CompanyID)
	if input.CompanyID != "" && !ok {
		return nil, sql.ErrNoRows
	}

	for index, jobHistory := range r.jobHistories {
		if jobHistory.ID == id {
			r.jobHistories[index].CompanyID = input.CompanyID
			r.jobHistories[index].Company = company.Name
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
	jobCompanies := make([]domain.JobCompany, len(r.jobCompanies))
	copy(jobCompanies, r.jobCompanies)

	return &domain.JobHistoryOptions{
		EmploymentTypes: employmentTypes,
		Companies:       jobCompanies,
	}, nil
}

func (r *testProfileRepository) CreateJobEmploymentType(_ context.Context, input domain.JobEmploymentTypeInput) (*domain.JobEmploymentType, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	employmentType := domain.JobEmploymentType{
		ID:        input.ID,
		Name:      input.Name,
		SortOrder: input.SortOrder,
	}
	if employmentType.SortOrder == 0 {
		employmentType.SortOrder = int64(len(r.employmentTypes) + 1)
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
			if input.SortOrder != 0 {
				r.employmentTypes[index].SortOrder = input.SortOrder
			}
			result := r.employmentTypes[index]
			return &result, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *testProfileRepository) CreateJobCompany(_ context.Context, input domain.JobCompanyInput) (*domain.JobCompany, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := input.ID
	if id == "" {
		id = strconv.Itoa(len(r.jobCompanies) + 1)
	}
	company := domain.JobCompany{
		ID:   id,
		Name: input.Name,
		URL:  input.URL,
	}
	r.jobCompanies = append(r.jobCompanies, company)

	return &company, nil
}

func (r *testProfileRepository) UpdateJobCompany(_ context.Context, id string, input domain.JobCompanyInput) (*domain.JobCompany, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, company := range r.jobCompanies {
		if company.ID == id {
			r.jobCompanies[index].Name = input.Name
			r.jobCompanies[index].URL = input.URL
			result := r.jobCompanies[index]
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
	company, ok := r.findJobCompany(input.CompanyID)
	if !ok {
		return nil, sql.ErrNoRows
	}

	project := domain.Project{
		ID:            "project_test_new",
		Name:          input.Name,
		CompanyID:     input.CompanyID,
		Company:       company.Name,
		StartYear:     input.StartYear,
		StartMonth:    input.StartMonth,
		EndYear:       input.EndYear,
		EndMonth:      input.EndMonth,
		Description:   input.Description,
		Role:          input.Role,
		TeamSize:      input.TeamSize,
		TechnologyIDs: input.TechnologyIDs,
		Technologies:  r.skillMasterNames(input.TechnologyIDs),
		Phases:        input.Phases,
		Achievements:  input.Achievements,
		IsDraft:       input.IsDraft,
		SortOrder:     int64(len(r.projects) + 1),
	}
	r.projects = append(r.projects, project)

	return &project, nil
}

func (r *testProfileRepository) UpdateProject(_ context.Context, id string, input domain.ProjectInput) (*domain.Project, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, project := range r.projects {
		if project.ID == id {
			company, ok := r.findJobCompany(input.CompanyID)
			if !ok {
				return nil, sql.ErrNoRows
			}

			r.projects[index].Name = input.Name
			r.projects[index].CompanyID = input.CompanyID
			r.projects[index].Company = company.Name
			r.projects[index].StartYear = input.StartYear
			r.projects[index].StartMonth = input.StartMonth
			r.projects[index].EndYear = input.EndYear
			r.projects[index].EndMonth = input.EndMonth
			r.projects[index].Description = input.Description
			r.projects[index].Role = input.Role
			r.projects[index].TeamSize = input.TeamSize
			r.projects[index].TechnologyIDs = input.TechnologyIDs
			r.projects[index].Technologies = r.skillMasterNames(input.TechnologyIDs)
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

func (r *testProfileRepository) ListProjectOptions(context.Context) (*domain.ProjectOptions, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	phases := make([]domain.ProjectPhase, len(r.projectPhases))
	copy(phases, r.projectPhases)

	return &domain.ProjectOptions{Phases: phases}, nil
}

func (r *testProfileRepository) CreateProjectPhase(_ context.Context, input domain.ProjectPhaseInput) (*domain.ProjectPhase, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := input.ID
	if id == "" {
		id = strconv.Itoa(len(r.projectPhases) + 1)
	}
	phase := domain.ProjectPhase{
		ID:        id,
		Name:      input.Name,
		SortOrder: input.SortOrder,
	}
	if phase.SortOrder == 0 {
		phase.SortOrder = int64(len(r.projectPhases) + 1)
	}
	r.projectPhases = append(r.projectPhases, phase)

	return &phase, nil
}

func (r *testProfileRepository) UpdateProjectPhase(_ context.Context, id string, input domain.ProjectPhaseInput) (*domain.ProjectPhase, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index, phase := range r.projectPhases {
		if phase.ID == id {
			r.projectPhases[index].Name = input.Name
			if input.SortOrder != 0 {
				r.projectPhases[index].SortOrder = input.SortOrder
			}
			result := r.projectPhases[index]
			return &result, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *testProfileRepository) findCategory(id string) (domain.SkillOption, bool) {
	for _, category := range r.categories {
		if category.ID == id {
			return category, true
		}
	}

	return domain.SkillOption{}, false
}

func (r *testProfileRepository) findSkillMaster(id string) (domain.SkillMaster, bool) {
	for _, skillMaster := range r.skillMasters {
		if skillMaster.ID == id {
			return skillMaster, true
		}
	}

	return domain.SkillMaster{}, false
}

func (r *testProfileRepository) skillMasterNames(ids []string) []string {
	names := make([]string, 0, len(ids))
	for _, id := range ids {
		skillMaster, ok := r.findSkillMaster(id)
		if !ok {
			continue
		}
		names = append(names, skillMaster.Name)
	}

	return names
}

func (r *testProfileRepository) findEmploymentType(id string) (domain.JobEmploymentType, bool) {
	for _, employmentType := range r.employmentTypes {
		if employmentType.ID == id {
			return employmentType, true
		}
	}

	return domain.JobEmploymentType{}, false
}

func (r *testProfileRepository) findJobCompany(id string) (domain.JobCompany, bool) {
	for _, company := range r.jobCompanies {
		if company.ID == id {
			return company, true
		}
	}

	return domain.JobCompany{}, false
}
