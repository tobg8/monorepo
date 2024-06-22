-- Deploy bearwise:create_schema_001 to pg

BEGIN;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE providers (
    id SERIAL PRIMARY KEY,
    value VARCHAR(255) NOT NULL
);

CREATE TABLE currencies (
    id SERIAL PRIMARY KEY,
    value VARCHAR(255) NOT NULL,
    provider_id INTEGER NOT NULL REFERENCES providers(id)
);

CREATE TYPE bot_status AS ENUM ('active', 'inactive', 'bug');

CREATE TABLE bots (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    status bot_status NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER NOT NULL REFERENCES users(id),
    currency_id INTEGER NOT NULL REFERENCES currencies(id)
);

CREATE TABLE intervals (
    id SERIAL PRIMARY KEY,
    value VARCHAR(255) NOT NULL
);

CREATE TABLE types (
    id SERIAL PRIMARY KEY,
    config JSON NOT NULL
);

CREATE TABLE indicators (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    symbol VARCHAR(255) NOT NULL,
    information TEXT,
    type_id INTEGER NOT NULL REFERENCES types(id)
);

CREATE TABLE parameters (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    optional BOOLEAN NOT NULL,
    type VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    "default" FLOAT
);

CREATE TABLE parameters_indicators (
    id SERIAL PRIMARY KEY,
    parameter_id INTEGER NOT NULL REFERENCES parameters(id),
    indicator_id INTEGER NOT NULL REFERENCES indicators(id)
);

CREATE TABLE families (
    id SERIAL PRIMARY KEY,
    value VARCHAR(255) NOT NULL
);

CREATE TABLE indicators_families (
    id SERIAL PRIMARY KEY,
    indicator_id INTEGER NOT NULL REFERENCES indicators(id),
    family_id INTEGER NOT NULL REFERENCES families(id)
);

CREATE TABLE alerts (
    id SERIAL PRIMARY KEY,
    config JSON NOT NULL,
    bot_id INTEGER NOT NULL REFERENCES bots(id),
    indicator_id INTEGER NOT NULL REFERENCES indicators(id),
    interval_id INTEGER NOT NULL REFERENCES intervals(id)
);

COMMIT;
