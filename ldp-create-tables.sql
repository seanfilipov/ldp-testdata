START TRANSACTION ISOLATION LEVEL SERIALIZABLE;

-- Extension for crosstab() requires superuser to install:
-- CREATE EXTENSION IF NOT EXISTS tablefunc;

CREATE SCHEMA norm;
CREATE SCHEMA internal;

CREATE TABLE internal.lock ();

-------------------------------------------------------------------------------
-- NORMALIZED SCHEMA ----------------------------------------------------------
-------------------------------------------------------------------------------

-- norm.groups

CREATE TABLE norm.groups (
    id           UUID NOT NULL PRIMARY KEY,
    group_name   TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
	CHECK (group_name <> ''),
    description  TEXT NOT NULL DEFAULT 'NOT AVAILABLE'
);

INSERT INTO norm.groups (id) VALUES ('00000000-0000-0000-0000-000000000000');

-- norm.users

/*
CREATE TABLE norm.users (
    id               UUID NOT NULL PRIMARY KEY,
    username         TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
        CHECK (username <> ''),
    barcode          TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    user_type        TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    active           BOOLEAN NOT NULL DEFAULT FALSE,
    patron_group_id  UUID NOT NULL REFERENCES norm.groups (id)
            DEFAULT '00000000-0000-0000-0000-000000000000'
);

INSERT INTO norm.users (id) VALUES ('00000000-0000-0000-0000-000000000000');
*/

-- norm.loans

/*
CREATE TABLE norm.loans (
    id           UUID NOT NULL PRIMARY KEY,
    user_id      UUID NOT NULL REFERENCES norm.users (id)
            DEFAULT '00000000-0000-0000-0000-000000000000',
    item_id      UUID NOT NULL
            DEFAULT '00000000-0000-0000-0000-000000000000',
    action       TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    status_name  TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    loan_date    TIMESTAMP NOT NULL DEFAULT 'epoch',
    due_date     TIMESTAMP NOT NULL DEFAULT 'epoch'
);

-- CREATE INDEX ON loans (loan_date);

-- INSERT INTO loans (id) VALUES ('00000000-0000-0000-0000-000000000000');
*/

-- norm.tmp_loans_locations

CREATE TABLE norm.tmp_loans_locations (
    loan_id        UUID NOT NULL PRIMARY KEY,
    location_name  TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
        CHECK (location_name <> '')
);

-- INSERT INTO tmp_loans_locations (loan_id)
    -- VALUES ('00000000-0000-0000-0000-000000000000');

-------------------------------------------------------------------------------
-- STAR SCHEMAS ---------------------------------------------------------------
-------------------------------------------------------------------------------

CREATE TABLE users (
    user_key           BIGSERIAL NOT NULL PRIMARY KEY,
        CHECK (user_key > 0),
    user_id            UUID NOT NULL,
    username           TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
        CHECK (username <> ''),
    barcode            TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    user_type          TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    active             BOOLEAN NOT NULL DEFAULT FALSE,
    group_name         TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
	CHECK (group_name <> ''),
    group_description  TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    record_time        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX ON users (user_id);

-- INSERT INTO users_dim (id) VALUES ('00000000-0000-0000-0000-000000000000');

/*
CREATE VIEW users_dim AS
SELECT u.id,
       u.username,
       u.barcode,
       u.user_type,
       u.active,
       g.group_name,
       g.description group_description
    FROM users u
        LEFT JOIN groups g ON u.patron_group_id = g.id;
*/

CREATE TABLE locations (
    id             TEXT NOT NULL PRIMARY KEY,
    location_name  TEXT NOT NULL DEFAULT 'NOT AVAILABLE'
);

-- INSERT INTO locations_dim (id) VALUES ('00000000-0000-0000-0000-000000000000');

/*
CREATE VIEW locations_dim AS
SELECT 'id-' || replace(lower(tll.location_name), ' ', '-') id,
       tll.location_name
    FROM (
        SELECT DISTINCT location_name FROM tmp_loans_locations
    ) tll;
*/

CREATE TABLE loans (
    loan_key     BIGSERIAL NOT NULL PRIMARY KEY,
        CHECK (loan_key > 0),
    loan_id      UUID NOT NULL UNIQUE,
    -- user_id      UUID NOT NULL --REFERENCES users_dim (id)
            -- DEFAULT '00000000-0000-0000-0000-000000000000',
    user_key     BIGINT NOT NULL REFERENCES users (user_key),
    location_id  TEXT NOT NULL REFERENCES locations (id)
            DEFAULT '00000000-0000-0000-0000-000000000000',
    item_id      UUID NOT NULL
            DEFAULT '00000000-0000-0000-0000-000000000000',
    action       TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    status_name  TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    loan_date    TIMESTAMP NOT NULL DEFAULT 'epoch',
    due_date     TIMESTAMP NOT NULL DEFAULT 'epoch'
);

CREATE INDEX ON loans (loan_date);

-- INSERT INTO loans_fact (id) VALUES ('00000000-0000-0000-0000-000000000000');

CREATE TABLE internal.loans (
    loan_id      UUID,
    user_id      UUID,
    location_id  TEXT,
    item_id      UUID,
    action       TEXT,
    status_name  TEXT,
    loan_date    TIMESTAMP,
    due_date     TIMESTAMP
);

/*
CREATE VIEW loans_fact AS
SELECT l.id,
       l.user_id,
       'id-' || replace(lower(tll.location_name), ' ', '-') location_id,
       l.item_id,
       l.action,
       l.status_name,
       l.loan_date,
       l.due_date
    FROM loans l
        LEFT JOIN tmp_loans_locations tll ON l.id = tll.loan_id;

*/

GRANT SELECT ON ALL TABLES IN SCHEMA public TO ldp;

COMMIT;

