package profiledb

import "context"

const getProfileByUserID = `-- name: GetProfileByUserID :one
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
LIMIT 1
`

func (q *Queries) GetProfileByUserID(ctx context.Context, userID string) (Profile, error) {
	row := q.db.QueryRowContext(ctx, getProfileByUserID, userID)
	var profile Profile
	err := row.Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FullName,
		&profile.Nickname,
		&profile.Location,
		&profile.Email,
		&profile.Phone,
		&profile.Summary,
		&profile.GithubUrl,
		&profile.ZennUrl,
		&profile.QiitaUrl,
		&profile.WebsiteUrl,
		&profile.PreferredWorkStyle,
		&profile.VisibilitySettings,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	return profile, err
}

const upsertProfile = `-- name: UpsertProfile :exec
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
  updated_at = VALUES(updated_at)
`

func (q *Queries) UpsertProfile(ctx context.Context, arg UpsertProfileParams) error {
	_, err := q.db.ExecContext(
		ctx,
		upsertProfile,
		arg.ID,
		arg.UserID,
		arg.FullName,
		arg.Nickname,
		arg.Location,
		arg.Email,
		arg.Phone,
		arg.Summary,
		arg.GithubUrl,
		arg.ZennUrl,
		arg.QiitaUrl,
		arg.WebsiteUrl,
		arg.PreferredWorkStyle,
		arg.VisibilitySettings,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}
