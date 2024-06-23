-- Deploy bearwise:create_schema_001 to pg

BEGIN;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE providers (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE currencies (
    id SERIAL PRIMARY KEY,
    pair VARCHAR NOT NULL,
    provider_id INTEGER NOT NULL,
    FOREIGN KEY (provider_id) REFERENCES providers(id)
);

CREATE TYPE bot_status AS ENUM ('active', 'inactive', 'bug');

CREATE TABLE bots (
    id SERIAL PRIMARY KEY,
    label VARCHAR NOT NULL,
    status bot_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    user_id INTEGER NOT NULL,
    currency_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (currency_id) REFERENCES currencies(id)
);

CREATE TABLE parameters (
    id SERIAL PRIMARY KEY,
    optional BOOLEAN NOT NULL,
    type_id INTEGER NOT NULL,
    label VARCHAR NOT NULL,
    description TEXT NOT NULL,
    default_value VARCHAR
);

CREATE TABLE parameters_options (
    id SERIAL PRIMARY KEY,
    option VARCHAR NOT NULL,
    parameter_id INTEGER NOT NULL,
    FOREIGN KEY (parameter_id) REFERENCES parameters(id)
);

CREATE TABLE types (
    id SERIAL PRIMARY KEY,
    label VARCHAR,
    config JSON NOT NULL
);

CREATE TABLE indicators (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    label VARCHAR NOT NULL,
    type_id INTEGER NOT NULL,
    FOREIGN KEY (type_id) REFERENCES types(id)
);

CREATE TABLE parameters_config (
    id SERIAL PRIMARY KEY,
    config JSON NOT NULL
);

CREATE TABLE parameters_indicators (
    id SERIAL PRIMARY KEY,
    parameter_id INTEGER NOT NULL,
    indicator_id INTEGER NOT NULL,
    FOREIGN KEY (parameter_id) REFERENCES parameters(id),
    FOREIGN KEY (indicator_id) REFERENCES indicators(id)
);

CREATE TABLE families (
    id SERIAL PRIMARY KEY,
    label_fr VARCHAR,
    label_en VARCHAR  NOT NULL,
    label_es VARCHAR,
    label_it VARCHAR,
    label_nl VARCHAR,
    label_de VARCHAR
);

CREATE TABLE indicators_families (
    id SERIAL PRIMARY KEY,
    indicator_id INTEGER NOT NULL,
    family_id INTEGER NOT NULL,
    FOREIGN KEY (indicator_id) REFERENCES indicators(id),
    FOREIGN KEY (family_id) REFERENCES families(id)
);

CREATE TABLE intervals (
    id SERIAL PRIMARY KEY,
    value VARCHAR NOT NULL
);

CREATE TABLE alerts (
    id SERIAL PRIMARY KEY,
    config JSON NOT NULL,
    bot_id INTEGER NOT NULL,
    indicator_id INTEGER NOT NULL,
    parameters_config_id INTEGER NOT NULL,
    interval_id INTEGER NOT NULL,
    FOREIGN KEY (bot_id) REFERENCES bots(id),
    FOREIGN KEY (indicator_id) REFERENCES indicators(id),
    FOREIGN KEY (parameters_config_id) REFERENCES parameters_config(id),
    FOREIGN KEY (interval_id) REFERENCES intervals(id)
);

COMMIT;
