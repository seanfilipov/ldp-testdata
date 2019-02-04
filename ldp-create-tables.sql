START TRANSACTION ISOLATION LEVEL SERIALIZABLE;

-- Extension for crosstab() requires superuser to install:
-- CREATE EXTENSION IF NOT EXISTS tablefunc;

CREATE SCHEMA normal;
COMMENT ON SCHEMA normal IS 'Extra tables used for denormalization';

CREATE SCHEMA loading;
COMMENT ON SCHEMA loading IS 'Internal area used for data loading';

CREATE TABLE loading.exlock ();
COMMENT ON TABLE loading.exlock IS 'Exclusive lock to prevent concurrent data loads';

-------------------------------------------------------------------------------
-- NORMALIZED SCHEMA ----------------------------------------------------------
-------------------------------------------------------------------------------

-- normal.groups

CREATE TABLE normal.groups (
    group_id     UUID NOT NULL PRIMARY KEY,
    group_name   TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
	CHECK (group_name <> ''),
    description  TEXT NOT NULL DEFAULT 'NOT AVAILABLE'
);

-- INSERT INTO normal.groups (id) VALUES ('00000000-0000-0000-0000-000000000000');

-- normal.users

/*
CREATE TABLE normal.users (
    id               UUID NOT NULL PRIMARY KEY,
    username         TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
        CHECK (username <> ''),
    barcode          TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    user_type        TEXT NOT NULL DEFAULT 'NOT AVAILABLE',
    active           BOOLEAN NOT NULL DEFAULT FALSE,
    patron_group_id  UUID NOT NULL REFERENCES normal.groups (id)
            DEFAULT '00000000-0000-0000-0000-000000000000'
);

INSERT INTO normal.users (id) VALUES ('00000000-0000-0000-0000-000000000000');
*/

-- normal.loans

/*
CREATE TABLE normal.loans (
    id           UUID NOT NULL PRIMARY KEY,
    user_id      UUID NOT NULL REFERENCES normal.users (id)
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

-- normal.tmp_loans_locations

CREATE TABLE normal.tmp_loans_locations (
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
    record_effective   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE users IS 'User records';

COMMENT ON COLUMN users.user_key IS 'Primary key of user record';
COMMENT ON COLUMN users.user_id IS 'FOLIO ID of the user';
COMMENT ON COLUMN users.username IS 'The user''s unique name, typically used for login';
COMMENT ON COLUMN users.barcode IS 'The library barcode for the user';
COMMENT ON COLUMN users.user_type IS 'The user class';
COMMENT ON COLUMN users.active IS
        'A flag to determine if a user can log in, take out loans, etc.';
COMMENT ON COLUMN users.group_name IS 'The group that the user is associated with';
COMMENT ON COLUMN users.group_description IS 'A description of the associated group';
COMMENT ON COLUMN users.record_effective IS
        'Date and time when the user record becomes effective';

CREATE INDEX ON users (user_id);

-- INSERT INTO users_dim (id) VALUES ('00000000-0000-0000-0000-000000000000');

CREATE TABLE loading.users (
    user_id          UUID NOT NULL PRIMARY KEY,
    username         TEXT NOT NULL,
    barcode          TEXT NOT NULL,
    user_type        TEXT NOT NULL,
    active           BOOLEAN NOT NULL,
    patron_group_id  UUID NOT NULL
);

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

COMMENT ON TABLE loans IS 'Loan transactions';

COMMENT ON COLUMN loans.loan_key IS 'Primary key of loan transaction';
COMMENT ON COLUMN loans.loan_id IS 'FOLIO ID of the loan transaction';
COMMENT ON COLUMN loans.user_key IS 'Foreign key of the user record for this loan';
-- COMMENT ON COLUMN loans.location_key IS '';
-- COMMENT ON COLUMN loans.item_key IS '';
COMMENT ON COLUMN loans.action IS 'Last action performed on the loan';
COMMENT ON COLUMN loans.status_name IS 'Overall status of the loan';
COMMENT ON COLUMN loans.loan_date IS 'Date and time when the loan began';
COMMENT ON COLUMN loans.due_date IS 'Date and time when the item is due to be returned';

CREATE INDEX ON loans (loan_date);

-- INSERT INTO loans_fact (id) VALUES ('00000000-0000-0000-0000-000000000000');

CREATE TABLE loading.loans (
    loan_id      UUID NOT NULL PRIMARY KEY,
    user_id      UUID NOT NULL,
    location_id  TEXT NOT NULL,
    item_id      UUID NOT NULL,
    action       TEXT NOT NULL,
    status_name  TEXT NOT NULL,
    loan_date    TIMESTAMP NOT NULL,
    due_date     TIMESTAMP NOT NULL
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

