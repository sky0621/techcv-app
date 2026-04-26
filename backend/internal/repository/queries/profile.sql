-- name: GetProfileByUserID :one
SELECT
  id,
  user_id,
  full_name,
  nickname,
  location,
  email,
  summary,
  github_url,
  zenn_url,
  qiita_url,
  website_url,
  occupation,
  employment_type,
  preferred_work_style,
  visibility_settings,
  created_at,
  updated_at
FROM profiles
WHERE user_id = ?
LIMIT 1;

-- name: UpsertProfile :exec
INSERT INTO profiles (
  id,
  user_id,
  full_name,
  nickname,
  location,
  email,
  summary,
  github_url,
  zenn_url,
  qiita_url,
  website_url,
  occupation,
  employment_type,
  preferred_work_style,
  visibility_settings,
  created_at,
  updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?,
  ?, ?, ?, ?, ?, ?, ?, ?, ?
)
ON CONFLICT(user_id) DO UPDATE SET
  full_name = excluded.full_name,
  nickname = excluded.nickname,
  location = excluded.location,
  email = excluded.email,
  summary = excluded.summary,
  github_url = excluded.github_url,
  zenn_url = excluded.zenn_url,
  qiita_url = excluded.qiita_url,
  website_url = excluded.website_url,
  occupation = excluded.occupation,
  employment_type = excluded.employment_type,
  preferred_work_style = excluded.preferred_work_style,
  visibility_settings = excluded.visibility_settings,
  updated_at = excluded.updated_at;

-- name: ListQualificationsByProfileID :many
SELECT
  id,
  profile_id,
  name,
  acquired_date,
  organization,
  url,
  memo,
  sort_order,
  created_at,
  updated_at
FROM profile_qualifications
WHERE profile_id = ?
ORDER BY sort_order ASC, created_at ASC, id ASC;

-- name: DeleteQualificationsByProfileID :exec
DELETE FROM profile_qualifications
WHERE profile_id = ?;

-- name: InsertQualification :exec
INSERT INTO profile_qualifications (
  id,
  profile_id,
  name,
  acquired_date,
  organization,
  url,
  memo,
  sort_order,
  created_at,
  updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);
