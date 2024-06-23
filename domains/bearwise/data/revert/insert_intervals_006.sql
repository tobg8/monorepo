-- Revert bearwise:insert_intervals_006 from pg

BEGIN;

TRUNCATE TABLE intervals RESTART IDENTITY CASCADE;


COMMIT;
