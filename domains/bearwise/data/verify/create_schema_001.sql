-- Verify bearwise:create_schema_001 on pg

BEGIN;

SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'users'
) AS users_exist;

SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'providers'
) AS providers_exist;

SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'currencies'
) AS currencies_exist;

SELECT EXISTS (
    SELECT 1
    FROM pg_type
    WHERE typname = 'bot_status'
) AS bot_status_exist;

SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'bots'
) AS bots_exist;

SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'intervals'
) AS intervals_exist;

SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'types'
) AS types_exist;

SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'indicators'
) AS indicators_exist;

SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'parameters'
) AS parameters_exist;

SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'parameters_indicators'
) AS parameters_indicators_exist;

SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'families'
) AS families_exist;

SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'indicators_families'
) AS indicators_families_exist;

SELECT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_name = 'alerts'
) AS alerts_exist;

ROLLBACK;
