-- name: ListSkillCategories :many
SELECT
  id,
  name,
  icon,
  sort_order,
  created_at,
  updated_at
FROM skill_categories
ORDER BY sort_order ASC, name ASC;

-- name: InsertSkillCategory :one
INSERT INTO skill_categories (
  id,
  name,
  icon,
  sort_order
)
SELECT
  ?,
  ?,
  ?,
  COALESCE(MAX(sort_order), 0) + 1
FROM skill_categories
RETURNING
  id,
  name,
  icon,
  sort_order,
  created_at,
  updated_at;

-- name: UpdateSkillCategory :one
UPDATE skill_categories
SET
  name = ?,
  icon = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING
  id,
  name,
  icon,
  sort_order,
  created_at,
  updated_at;

-- name: ListSkillProficiencyLevels :many
SELECT
  id,
  name,
  sort_order,
  created_at,
  updated_at
FROM skill_proficiency_levels
ORDER BY sort_order ASC, name ASC;
