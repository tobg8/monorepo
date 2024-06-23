-- Verify bearwise:insert_types_003 on pg

BEGIN;

DO $$
    DECLARE
        bollinger_exists BOOLEAN;
        candle_exists BOOLEAN;
        single_exists BOOLEAN;
    BEGIN
        -- Check existence for each label
        SELECT EXISTS (
            SELECT 1
            FROM types
            WHERE label = 'bollinger'
        ) INTO bollinger_exists;

        SELECT EXISTS (
            SELECT 1
            FROM types
            WHERE label = 'candle'
        ) INTO candle_exists;

        SELECT EXISTS (
            SELECT 1
            FROM types
            WHERE label = 'single'
        ) INTO single_exists;

        -- Assert that all required records exist
        ASSERT bollinger_exists = true, 'Record with label = ''bollinger'' does not exist';
        ASSERT candle_exists = true, 'Record with label = ''candle'' does not exist';
        ASSERT single_exists = true, 'Record with label = ''single'' does not exist';

        -- If all assertions pass, print success
        RAISE NOTICE 'Verification successful: All required records exist';

        -- If any assertion fails, an exception will be raised automatically
    END $$;

COMMIT;