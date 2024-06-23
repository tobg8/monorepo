-- Revert bearwise:insert_families_002 from pg

BEGIN;

TRUNCATE TABLE families RESTART IDENTITY CASCADE;

COMMIT;
