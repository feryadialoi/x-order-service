-- migrate:up
CREATE TABLE order_states
(
    correlation_id VARCHAR(36) NOT NULL PRIMARY KEY,
    current_state  VARCHAR(64) NOT NULL,
    data           JSONB       NOT NULL,
    created_at     TIMESTAMPTZ NULL,
    updated_at     TIMESTAMPTZ NULL
);

-- migrate:down
DROP TABLE order_states;