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

-- name: ListSkillProficiencyLevels :many
SELECT
  id,
  name,
  sort_order,
  created_at,
  updated_at
FROM skill_proficiency_levels
ORDER BY sort_order ASC, name ASC;
