-- name: ListProjects :many
SELECT
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
FROM projects
ORDER BY sort_order ASC, start_year DESC, start_month DESC, name ASC;

-- name: GetProject :one
SELECT
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
FROM projects
WHERE id = ?;

-- name: InsertProject :one
INSERT INTO projects (
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
  updated_at;

-- name: UpdateProject :one
UPDATE projects
SET
  name = ?,
  company = ?,
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
