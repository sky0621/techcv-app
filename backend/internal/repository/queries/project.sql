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
