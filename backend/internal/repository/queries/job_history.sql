-- name: ListJobHistories :many
SELECT
  job_histories.id,
  job_histories.company,
  job_histories.display_name,
  job_histories.start_year,
  job_histories.start_month,
  job_histories.end_year,
  job_histories.end_month,
  job_histories.employment_type_id,
  job_employment_types.name AS employment_type,
  job_histories.project_count,
  job_histories.sort_order,
  job_histories.created_at,
  job_histories.updated_at
FROM job_histories
JOIN job_employment_types ON job_employment_types.id = job_histories.employment_type_id
ORDER BY job_histories.sort_order ASC, job_histories.start_year DESC, job_histories.start_month DESC, job_histories.company ASC;

-- name: GetJobHistory :one
SELECT
  job_histories.id,
  job_histories.company,
  job_histories.display_name,
  job_histories.start_year,
  job_histories.start_month,
  job_histories.end_year,
  job_histories.end_month,
  job_histories.employment_type_id,
  job_employment_types.name AS employment_type,
  job_histories.project_count,
  job_histories.sort_order,
  job_histories.created_at,
  job_histories.updated_at
FROM job_histories
JOIN job_employment_types ON job_employment_types.id = job_histories.employment_type_id
WHERE job_histories.id = ?;

-- name: InsertJobHistory :one
INSERT INTO job_histories (
  id,
  company,
  display_name,
  start_year,
  start_month,
  end_year,
  end_month,
  employment_type_id,
  project_count,
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
  0,
  COALESCE(MIN(sort_order), 0) - 1
FROM job_histories
RETURNING
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
  updated_at;

-- name: UpdateJobHistory :one
UPDATE job_histories
SET
  company = ?,
  display_name = ?,
  start_year = ?,
  start_month = ?,
  end_year = ?,
  end_month = ?,
  employment_type_id = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING
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
  updated_at;

-- name: DeleteJobHistory :execrows
DELETE FROM job_histories
WHERE id = ?;

-- name: ListJobEmploymentTypes :many
SELECT
  id,
  name,
  sort_order,
  created_at,
  updated_at
FROM job_employment_types
ORDER BY sort_order ASC, name ASC;

-- name: InsertJobEmploymentType :one
INSERT INTO job_employment_types (
  id,
  name,
  sort_order
)
SELECT
  ?,
  ?,
  COALESCE(MAX(sort_order), 0) + 1
FROM job_employment_types
RETURNING
  id,
  name,
  sort_order,
  created_at,
  updated_at;

-- name: UpdateJobEmploymentType :one
UPDATE job_employment_types
SET
  name = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING
  id,
  name,
  sort_order,
  created_at,
  updated_at;
