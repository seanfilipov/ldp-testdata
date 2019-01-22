-- Requires superuser:
-- CREATE EXTENSION IF NOT EXISTS tablefunc;

CREATE TABLE groups (
    id          UUID NOT NULL PRIMARY KEY,
    group_name  TEXT NOT NULL UNIQUE,
        CHECK (group_name <> '')
);

CREATE TABLE users (
    id               UUID NOT NULL PRIMARY KEY,
    -- username         TEXT NOT NULL UNIQUE,
    username         TEXT NOT NULL, -- TODO fix test data
        CHECK (username <> ''),
    active           BOOLEAN NOT NULL,
    patron_group_id  UUID NOT NULL
);

CREATE TABLE loans (
    id         UUID NOT NULL PRIMARY KEY,
    user_id    UUID NOT NULL,
    loan_date  TIMESTAMP NOT NULL
);

CREATE TABLE tmp_loans_locations (
    loan_id        UUID NOT NULL PRIMARY KEY,
    location_name  TEXT NOT NULL,
        CHECK (location_name <> '')
);

CREATE VIEW users_dim AS
SELECT u.id, u.username, u.active, g.group_name
    FROM users u
        LEFT JOIN groups g ON u.patron_group_id = g.id;

CREATE VIEW locations_dim AS
SELECT 'id-' || replace(lower(tll.location_name), ' ', '-') location_id,
       tll.location_name
    FROM (
        SELECT DISTINCT location_name FROM tmp_loans_locations
    ) tll;

CREATE VIEW loans_fact AS
SELECT l.id,
       l.user_id,
       'id-' || replace(lower(tll.location_name), ' ', '-') location_id,
       l.loan_date
    FROM loans l
        JOIN tmp_loans_locations tll ON l.id = tll.loan_id;


