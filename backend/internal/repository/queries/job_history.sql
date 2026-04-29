-- name: ListJobHistories :many
SELECT
  id,
  company,
  start_date,
  end_date,
  employment_type,
  project_count,
  sort_order,
  created_at,
  updated_at
FROM job_histories
ORDER BY sort_order ASC, start_date DESC, company ASC;

-- name: GetJobHistory :one
SELECT
  id,
  company,
  start_date,
  end_date,
  employment_type,
  project_count,
  sort_order,
  created_at,
  updated_at
FROM job_histories
WHERE id = ?;

-- name: InsertJobHistory :one
INSERT INTO job_histories (
  id,
  company,
  start_date,
  end_date,
  employment_type,
  project_count,
  sort_order
)
SELECT
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
  start_date,
  end_date,
  employment_type,
  project_count,
  sort_order,
  created_at,
  updated_at;

-- name: UpdateJobHistory :one
UPDATE job_histories
SET
  company = ?,
  start_date = ?,
  end_date = ?,
  employment_type = ?,
  updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING
  id,
  company,
  start_date,
  end_date,
  employment_type,
  project_count,
  sort_order,
  created_at,
  updated_at;

-- name: DeleteJobHistory :execrows
DELETE FROM job_histories
WHERE id = ?;
