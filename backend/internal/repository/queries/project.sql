-- name: ListProjects :many
SELECT
  projects.id,
  projects.name,
  CAST(projects.company_id AS TEXT) AS company_id,
  job_companies.name AS company,
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
JOIN job_companies ON job_companies.id = projects.company_id
ORDER BY projects.sort_order ASC, projects.start_year DESC, projects.start_month DESC, projects.name ASC;

-- name: GetProject :one
SELECT
  projects.id,
  projects.name,
  CAST(projects.company_id AS TEXT) AS company_id,
  job_companies.name AS company,
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
JOIN job_companies ON job_companies.id = projects.company_id
WHERE projects.id = ?;

-- name: InsertProject :one
INSERT INTO projects (
  id,
  name,
  company_id,
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
)
SELECT
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  COALESCE(MIN(sort_order), 0) - 1
FROM projects
RETURNING
  id,
  name,
  CAST(company_id AS TEXT) AS company_id,
  '' AS company,
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
  updated_at;

-- name: UpdateProject :one
UPDATE projects
SET
  name = ?,
  company_id = ?,
  start_year = ?,
  start_month = ?,
  end_year = ?,
  end_month = ?,
  description = ?,
  role = ?,
  team_size = ?,
  technologies = ?,
  phases = ?,
  achievements = ?,
  is_draft = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING
  id,
  name,
  CAST(company_id AS TEXT) AS company_id,
  '' AS company,
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
  updated_at;

-- name: DeleteProject :execrows
DELETE FROM projects
WHERE id = ?;

-- name: ListProjectPhases :many
SELECT
  id,
  name,
  sort_order,
  created_at,
  updated_at
FROM project_phases
ORDER BY sort_order ASC, name ASC;

-- name: InsertProjectPhase :one
INSERT INTO project_phases (
  id,
  name,
  sort_order
)
SELECT
  ?,
  ?,
  COALESCE(NULLIF(?, 0), COALESCE(MAX(sort_order), 0) + 1)
FROM project_phases
RETURNING
  id,
  name,
  sort_order,
  created_at,
  updated_at;

-- name: InsertProjectPhaseAuto :one
INSERT INTO project_phases (
  name,
  sort_order
)
SELECT
  ?,
  COALESCE(NULLIF(?, 0), COALESCE(MAX(sort_order), 0) + 1)
FROM project_phases
RETURNING
  id,
  name,
  sort_order,
  created_at,
  updated_at;

-- name: UpdateProjectPhase :one
UPDATE project_phases
SET
  name = ?,
  sort_order = COALESCE(NULLIF(?, 0), sort_order),
  updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING
  id,
  name,
  sort_order,
  created_at,
  updated_at;
