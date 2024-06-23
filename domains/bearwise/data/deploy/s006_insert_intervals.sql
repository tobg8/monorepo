-- Deploy bearwise:insert_intervals_006 to pg

BEGIN;

INSERT INTO intervals (value) VALUES
    ('1W'),
    ('1D'),
    ('4H'),
    ('1H')
;

COMMIT;
