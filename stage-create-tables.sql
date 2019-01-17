START TRANSACTION;


CREATE TABLE stage (
    id      BIGSERIAL NOT NULL PRIMARY KEY,
    jtype   TEXT NOT NULL,
    jid     UUID NOT NULL,
    j       JSONB NOT NULL,
);

CREATE INDEX ON stage (jtype, jid);


COMMIT;

