CREATE DATABASE IF NOT EXISTS benchmark_repository;
USE benchmark_repository;


CREATE TABLE IF NOT EXISTS repository (
    id INT AUTO_INCREMENT PRIMARY KEY,
    owner VARCHAR(255) NOT NULL,
    repo VARCHAR(255) NOT NULL,
    star_number INT NOT NULL,
    fork_number INT NOT NULL,
    collected_on DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_entry
        UNIQUE (owner, repo, collected_on)
);

CREATE TABLE IF NOT EXISTS pull_requests(
    repository_id INT UNIQUE NOT NULL,
    open_pull_requests_number INT NOT NULL,
    closed_pull_requests_number INT NOT NULL,
    average_comments_per_pull_request FLOAT8 NOT NULL,
    average_commits_per_pull_request FLOAT8 NOT NULL,
    CONSTRAINT pull_requests_fk
        FOREIGN KEY (repository_id)
        REFERENCES repository (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS issues(
    repository_id INT UNIQUE NOT NULL,
    open_issues_number INT NOT NULL,
    closed_issues_number INT NOT NULL,
    average_comments_per_issue FLOAT8 NOT NULL,
    CONSTRAINT issues_fk
        FOREIGN KEY (repository_id)
        REFERENCES repository (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS community_documents(
    repository_id INT UNIQUE NOT NULL,
    health_percentage int NOT NULL,
    has_code_of_conduct BOOLEAN NOT NULL,
    has_contributing BOOLEAN NOT NULL,
    has_issue_template BOOLEAN NOT NULL,
    has_pull_request_template BOOLEAN NOT NULL,
    has_license BOOLEAN NOT NULL,
    has_readme BOOLEAN NOT NULL,
    has_content_reports_enabled BOOLEAN NOT NULL,
    has_wiki BOOLEAN NOT NULL,
    CONSTRAINT community_documents_fk
        FOREIGN KEY (repository_id)
        REFERENCES repository (id) ON DELETE CASCADE
);