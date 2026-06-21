-- name: GetProfile :one
SELECT
  id,
  family_name,
  given_name,
  nickname,
  avatar_url,
  birthday_year,
  birthday_month,
  birthday_day,
  location,
  email,
  pr,
  occupation,
  employment_type,
  preferred_work_style,
  visibility_settings,
  created_at,
  updated_at
FROM profiles
ORDER BY id ASC
LIMIT 1;

-- name: UpsertProfile :exec
INSERT INTO profiles (
  id,
  family_name,
  given_name,
  nickname,
  avatar_url,
  birthday_year,
  birthday_month,
  birthday_day,
  location,
  email,
  pr,
  occupation,
  employment_type,
  preferred_work_style,
  visibility_settings,
  created_at,
  updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?,
  ?, ?, ?, ?, ?, ?, ?, ?
)
ON CONFLICT(id) DO UPDATE SET
  family_name = excluded.family_name,
  given_name = excluded.given_name,
  nickname = excluded.nickname,
  avatar_url = excluded.avatar_url,
  birthday_year = excluded.birthday_year,
  birthday_month = excluded.birthday_month,
  birthday_day = excluded.birthday_day,
  location = excluded.location,
  email = excluded.email,
  pr = excluded.pr,
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
