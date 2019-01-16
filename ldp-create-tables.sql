START TRANSACTION;


CREATE TABLE users (
    user_id     UUID NOT NULL,
    group_name  TEXT NOT NULL
);


CREATE TABLE locations (
    location_id    TEXT NOT NULL, -- Temporary
    -- location_id    UUID NOT NULL,
    location_name  TEXT NOT NULL
);


CREATE TABLE loans (
    loan_id      UUID NOT NULL,
    user_id      UUID NOT NULL,
    location_id  TEXT NOT NULL, -- Temporary
    -- location_id  UUID NOT NULL,
    loan_date    TIMESTAMP NOT NULL
);


COMMIT;

