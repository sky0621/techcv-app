-- name: GetProfileByUserID :one
SELECT
  id,
  user_id,
  full_name,
  nickname,
  location,
  email,
  phone,
  summary,
  github_url,
  zenn_url,
  qiita_url,
  website_url,
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
  phone,
  summary,
  github_url,
  zenn_url,
  qiita_url,
  website_url,
  preferred_work_style,
  visibility_settings,
  created_at,
  updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?,
  ?, ?, ?, ?, ?, ?, ?, ?
)
ON DUPLICATE KEY UPDATE
  full_name = VALUES(full_name),
  nickname = VALUES(nickname),
  location = VALUES(location),
  email = VALUES(email),
  phone = VALUES(phone),
  summary = VALUES(summary),
  github_url = VALUES(github_url),
  zenn_url = VALUES(zenn_url),
  qiita_url = VALUES(qiita_url),
  website_url = VALUES(website_url),
  preferred_work_style = VALUES(preferred_work_style),
  visibility_settings = VALUES(visibility_settings),
  updated_at = VALUES(updated_at);
