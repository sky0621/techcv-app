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
    occupation TEXT NOT NULL,
    employment_type TEXT NOT NULL,
    preferred_work_style TEXT NOT NULL,
    visibility_settings TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS profile_qualifications (
    id TEXT NOT NULL PRIMARY KEY,
    profile_id TEXT NOT NULL,
    name TEXT NOT NULL,
    acquired_date TEXT NOT NULL,
    organization TEXT NOT NULL,
    url TEXT NOT NULL,
    memo TEXT NOT NULL,
    sort_order INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (profile_id) REFERENCES profiles(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS skill_categories (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    icon TEXT NOT NULL,
    sort_order INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS skill_proficiency_levels (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    sort_order INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS skill_masters (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    category_id TEXT NOT NULL,
    sort_order INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES skill_categories(id)
);

CREATE TABLE IF NOT EXISTS skills (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    category_id TEXT NOT NULL,
    experience INTEGER NOT NULL,
    proficiency_level_id TEXT NOT NULL,
    sort_order INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES skill_categories(id),
    FOREIGN KEY (proficiency_level_id) REFERENCES skill_proficiency_levels(id)
);

CREATE TABLE IF NOT EXISTS job_employment_types (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    sort_order INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS job_histories (
    id TEXT NOT NULL PRIMARY KEY,
    company TEXT NOT NULL,
    display_name TEXT NOT NULL,
    start_year INTEGER NOT NULL,
    start_month INTEGER NOT NULL,
    end_year INTEGER,
    end_month INTEGER,
    employment_type_id TEXT NOT NULL,
    project_count INTEGER NOT NULL DEFAULT 0,
    sort_order INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employment_type_id) REFERENCES job_employment_types(id)
);

CREATE TABLE IF NOT EXISTS projects (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    company TEXT NOT NULL,
    start_year INTEGER NOT NULL,
    start_month INTEGER NOT NULL,
    end_year INTEGER,
    end_month INTEGER,
    description TEXT NOT NULL,
    role TEXT NOT NULL,
    team_size TEXT NOT NULL,
    technologies TEXT NOT NULL,
    phases TEXT NOT NULL,
    achievements TEXT NOT NULL,
    is_draft INTEGER NOT NULL DEFAULT 0,
    sort_order INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
