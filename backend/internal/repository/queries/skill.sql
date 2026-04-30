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

-- name: ListSkillMasters :many
SELECT
  skill_masters.id,
  skill_masters.name,
  skill_masters.category_id,
  skill_categories.name AS category_name,
  skill_masters.sort_order,
  skill_masters.created_at,
  skill_masters.updated_at
FROM skill_masters
JOIN skill_categories ON skill_categories.id = skill_masters.category_id
ORDER BY skill_masters.sort_order ASC, skill_masters.name ASC;

-- name: GetSkillMaster :one
SELECT
  skill_masters.id,
  skill_masters.name,
  skill_masters.category_id,
  skill_categories.name AS category_name,
  skill_masters.sort_order,
  skill_masters.created_at,
  skill_masters.updated_at
FROM skill_masters
JOIN skill_categories ON skill_categories.id = skill_masters.category_id
WHERE skill_masters.id = ?;

-- name: InsertSkillMaster :one
INSERT INTO skill_masters (
  id,
  name,
  category_id,
  sort_order
)
SELECT
  ?,
  ?,
  ?,
  COALESCE(MAX(sort_order), 0) + 1
FROM skill_masters
RETURNING
  id,
  name,
  category_id,
  sort_order,
  created_at,
  updated_at;

-- name: UpdateSkillMaster :one
UPDATE skill_masters
SET
  name = ?,
  category_id = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING
  id,
  name,
  category_id,
  sort_order,
  created_at,
  updated_at;

-- name: ListSkills :many
SELECT
  skills.id,
  skills.name,
  skills.category_id,
  skill_categories.name AS category_name,
  skills.experience,
  skills.proficiency_level_id,
  skill_proficiency_levels.name AS proficiency_name,
  skills.sort_order,
  skills.created_at,
  skills.updated_at
FROM skills
JOIN skill_categories ON skill_categories.id = skills.category_id
JOIN skill_proficiency_levels ON skill_proficiency_levels.id = skills.proficiency_level_id
ORDER BY skills.sort_order ASC, skills.name ASC;

-- name: GetSkill :one
SELECT
  skills.id,
  skills.name,
  skills.category_id,
  skill_categories.name AS category_name,
  skills.experience,
  skills.proficiency_level_id,
  skill_proficiency_levels.name AS proficiency_name,
  skills.sort_order,
  skills.created_at,
  skills.updated_at
FROM skills
JOIN skill_categories ON skill_categories.id = skills.category_id
JOIN skill_proficiency_levels ON skill_proficiency_levels.id = skills.proficiency_level_id
WHERE skills.id = ?;

-- name: InsertSkill :one
INSERT INTO skills (
  id,
  name,
  category_id,
  experience,
  proficiency_level_id,
  sort_order
)
SELECT
  ?,
  ?,
  ?,
  ?,
  ?,
  COALESCE(MAX(sort_order), 0) + 1
FROM skills
RETURNING
  id,
  name,
  category_id,
  experience,
  proficiency_level_id,
  sort_order,
  created_at,
  updated_at;

-- name: UpdateSkill :one
UPDATE skills
SET
  name = ?,
  category_id = ?,
  experience = ?,
  proficiency_level_id = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING
  id,
  name,
  category_id,
  experience,
  proficiency_level_id,
  sort_order,
  created_at,
  updated_at;

-- name: DeleteSkill :execrows
DELETE FROM skills
WHERE id = ?;
