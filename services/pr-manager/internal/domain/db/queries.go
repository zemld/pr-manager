package db

const (
	CreateUsersTable = `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT NOT NULL,
		username TEXT NOT NULL,
		team_name TEXT NOT NULL,
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		PRIMARY KEY (id)
	)
	`
	CreatePullRequestsStatusesTable = `
	CREATE TABLE IF NOT EXISTS pull_requests_statuses (
		id SERIAL,
		status TEXT NOT NULL UNIQUE,
		PRIMARY KEY (id)
	)
	`
	CreatePullRequestsTable = `
	CREATE TABLE IF NOT EXISTS pull_requests (
		id TEXT NOT NULL,
		name TEXT NOT NULL,
		author_id TEXT NOT NULL,
		status_id INT NOT NULL,
		assigned_reviewers TEXT NOT NULL,
		created_at TIMESTAMP,
		merged_at TIMESTAMP,
		PRIMARY KEY (id),
		FOREIGN KEY (author_id) REFERENCES users (id),
		FOREIGN KEY (status_id) REFERENCES pull_requests_statuses (id)
	)
	`

	FillPullRequestsStatusesTable = `
	INSERT INTO pull_requests_statuses (status) VALUES ('open'), ('merged')
	ON CONFLICT (status) DO NOTHING
	`

	SelectTeam = `
	SELECT 
		team_name,
		json_agg(
			json_build_object(
				'user_id', id,
				'username', username,
				'is_active', is_active
			)
		) as members
	FROM users 
	WHERE team_name = $1
	GROUP BY team_name
	`

	UpdateUserStatus = `
	UPDATE users SET is_active = $1 WHERE id = $2
	`
	InsertUser = `
	INSERT INTO users (id, username, team_name, is_active) VALUES ($1, $2, $3, $4)
	ON CONFLICT (id) DO NOTHING
	`

	CreatePullRequest = `
		INSERT INTO
			pull_requests
			(id, name, author_id, status_id, assigned_reviewers, created_at, merged_at)
		VALUES
			($1, $2, $3, (SELECT id FROM pull_requests_statuses WHERE status = 'open' LIMIT 1), $4, NOW(), NULL)
		ON CONFLICT
			(id) DO NOTHING
	`
	MergePullRequest = `
		UPDATE
			pull_requests
		SET
			status_id = (SELECT id FROM pull_requests_statuses WHERE status = 'merged' LIMIT 1),
			merged_at = NOW()
		WHERE
			id = $1
			AND status_id != (SELECT id FROM pull_requests_statuses WHERE status = 'merged' LIMIT 1)
	`
	ReassignPullRequest = `
	UPDATE
		pull_requests
	SET
		assigned_reviewers = $2
	WHERE
		id = $1
	`
	UserPullRequestsReviews = `
	SELECT 
		id, 
		name, 
		author_id, 
		(
			SELECT
				status
			FROM
				pull_requests_statuses
			WHERE id = pull_requests.status_id
			LIMIT 1
		) as status,
		assigned_reviewers,
		created_at,
		merged_at
	FROM
		pull_requests
	WHERE
		assigned_reviewers LIKE '%' || $1 || '%'
	`
	SelectPullRequest = `
	SELECT
		id,
		name,
		author_id,
		(
			SELECT
				status
			FROM
				pull_requests_statuses
			WHERE id = pull_requests.status_id
			LIMIT 1
		) as status,
		assigned_reviewers,
		created_at,
		merged_at
	FROM
		pull_requests
	WHERE id = $1
	LIMIT 1
	`
	SelectUser = `
	SELECT
		id as user_id,
		username,
		team_name,
		is_active
	FROM
		users
	WHERE id = $1
	LIMIT 1
	`
)
