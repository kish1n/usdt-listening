-- +migrate Up

CREATE TABLE transfers (
   id UUID PRIMARY KEY,
   from_address VARCHAR(42) NOT NULL,
   to_address VARCHAR(42) NOT NULL,
   value NUMERIC NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE transfers;