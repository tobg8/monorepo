-- Deploy bearwise:insert_types_003 to pg

BEGIN;

-- Insert first JSON object with label "bollinger"
INSERT INTO types (label, config) VALUES (
    'bollinger',
    '{
     "valueUpperBand": 23491.10808422158,
     "valueMiddleBand": 21935.281499999997,
     "valueLowerBand": 20379.454915778413
    }'::json
);

-- Insert second JSON object with label "candle"
INSERT INTO types (label, config) VALUES (
    'candle',
    '{
     "timestampHuman": "2021-01-14 15:00:00 (Thursday) UTC",
     "timestamp": 1610636400,
     "open": 39577.53,
     "high": 39666,
     "low": 39294.7,
     "close": 39607.09,
     "volume": 1211.2841909999893
    }'::json
);

-- Insert third JSON object with label "single"
INSERT INTO types (label, config) VALUES (
    'single',
    '{
     "value": 177.65955750772062
    }'::json
);

COMMIT;