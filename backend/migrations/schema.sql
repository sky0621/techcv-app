CREATE TABLE IF NOT EXISTS profiles (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE,
    full_name TEXT NOT NULL,
    nickname TEXT NOT NULL,
    location TEXT NOT NULL,
    email TEXT NOT NULL,
    summary TEXT NOT NULL,
    github_url TEXT NOT NULL,
    zenn_url TEXT NOT NULL,
    qiita_url TEXT NOT NULL,
    website_url TEXT NOT NULL,
    preferred_work_style TEXT NOT NULL,
    visibility_settings TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
