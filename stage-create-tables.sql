START TRANSACTION;


CREATE SCHEMA denorm;
CREATE SCHEMA dim;
CREATE SCHEMA tmp;


CREATE TYPE ldutype AS ENUM (
    'groups',
    'loans',
    'users',
    'tmp_loans_locations'
);


CREATE TABLE stage (
    id     BIGSERIAL NOT NULL PRIMARY KEY,
    jtype  ldutype NOT NULL,
    jid    UUID NOT NULL,
    j      JSONB NOT NULL
);


CREATE TABLE denorm.groups (
    id           UUID NOT NULL PRIMARY KEY,
    name         TEXT NOT NULL,
    description  TEXT NOT NULL
);


CREATE TABLE dim.users (
    id        UUID NOT NULL PRIMARY KEY,
    group_id  UUID NOT NULL
);


CREATE TABLE tmp.loans_locations (
    loan_id        UUID NOT NULL PRIMARY KEY,
    location_name  TEXT NOT NULL
);


COMMIT;

