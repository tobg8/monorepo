-- Verify bearwise:create_schema_001 on pg

BEGIN;

DO $$
    DECLARE
        users_exist BOOLEAN;
        providers_exist BOOLEAN;
        currencies_exist BOOLEAN;
        bot_status_exist BOOLEAN;
        bots_exist BOOLEAN;
        parameters_exist BOOLEAN;
        parameters_options_exist BOOLEAN;
        indicators_exist BOOLEAN;
        alerts_exist BOOLEAN;
        intervals_exist BOOLEAN;
        parameters_config_exist BOOLEAN;
        parameters_indicators_exist BOOLEAN;
        indicators_families_exist BOOLEAN;
        families_exist BOOLEAN;
        types_exist BOOLEAN;
    BEGIN
        -- Check existence for each table or type
        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = 'users'
        ) INTO users_exist;

        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = 'providers'
        ) INTO providers_exist;

        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = 'currencies'
        ) INTO currencies_exist;

        SELECT EXISTS (
            SELECT 1
            FROM pg_type
            WHERE typname = 'bot_status'
        ) INTO bot_status_exist;

        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = 'bots'
        ) INTO bots_exist;

        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = 'parameters'
        ) INTO parameters_exist;

        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = 'parameters_options'
        ) INTO parameters_options_exist;

        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = 'indicators'
        ) INTO indicators_exist;

        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = 'alerts'
        ) INTO alerts_exist;

        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = 'intervals'
        ) INTO intervals_exist;

        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = 'parameters_config'
        ) INTO parameters_config_exist;

        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = 'parameters_indicators'
        ) INTO parameters_indicators_exist;

        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = 'indicators_families'
        ) INTO indicators_families_exist;

        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = 'families'
        ) INTO families_exist;

        SELECT EXISTS (
            SELECT 1
            FROM information_schema.tables
            WHERE table_name = 'types'
        ) INTO types_exist;

        -- Assert that all required tables and types exist
        ASSERT users_exist = true, 'Table users does not exist';
        ASSERT providers_exist = true, 'Table providers does not exist';
        ASSERT currencies_exist = true, 'Table currencies does not exist';
        ASSERT bot_status_exist = true, 'Type bot_status does not exist';
        ASSERT bots_exist = true, 'Table bots does not exist';
        ASSERT parameters_exist = true, 'Table parameters does not exist';
        ASSERT parameters_options_exist = true, 'Table parameters_options does not exist';
        ASSERT indicators_exist = true, 'Table indicators does not exist';
        ASSERT alerts_exist = true, 'Table alerts does not exist';
        ASSERT intervals_exist = true, 'Table intervals does not exist';
        ASSERT parameters_config_exist = true, 'Table parameters_config does not exist';
        ASSERT parameters_indicators_exist = true, 'Table parameters_indicators does not exist';
        ASSERT indicators_families_exist = true, 'Table indicators_families does not exist';
        ASSERT families_exist = true, 'Table families does not exist';
        ASSERT types_exist = true, 'Table types does not exist';

        -- If all assertions pass, print success
        RAISE NOTICE 'Verification successful: All required tables and types exist';

        -- If any assertion fails, an exception will be raised automatically
    END $$;

COMMIT;