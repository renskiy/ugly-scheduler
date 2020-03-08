-- +migrate Up
CREATE TABLE events
(
    id      SERIAL PRIMARY KEY,
    "when"  TIMESTAMPTZ        NOT NULL,
    message VARCHAR(1024)      NOT NULL,
    done    BOOL DEFAULT FALSE NOT NULL
);

CREATE INDEX events_when_idx
    ON events ("when") WHERE NOT done;

-- +migrate Down
DROP TABLE events CASCADE;
