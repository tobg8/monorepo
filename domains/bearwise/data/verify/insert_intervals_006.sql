-- Verify bearwise:insert_intervals_006 on pg

BEGIN;

DO $$
DECLARE
    interval_1W_exists BOOLEAN;
    interval_1D_exists BOOLEAN;
    interval_4H_exists BOOLEAN;
    interval_1H_exists BOOLEAN;
BEGIN
    -- Check existence for each interval value
SELECT EXISTS (SELECT 1 FROM intervals WHERE value = '1W') INTO interval_1W_exists;
SELECT EXISTS (SELECT 1 FROM intervals WHERE value = '1D') INTO interval_1D_exists;
SELECT EXISTS (SELECT 1 FROM intervals WHERE value = '4H') INTO interval_4H_exists;
SELECT EXISTS (SELECT 1 FROM intervals WHERE value = '1H') INTO interval_1H_exists;

-- Assert that all required records exist
ASSERT interval_1W_exists = true, 'Record with value = ''1W'' does not exist';
    ASSERT interval_1D_exists = true, 'Record with value = ''1D'' does not exist';
    ASSERT interval_4H_exists = true, 'Record with value = ''4H'' does not exist';
    ASSERT interval_1H_exists = true, 'Record with value = ''1H'' does not exist';

    -- If all assertions pass, print success
    RAISE NOTICE 'Verification successful: All required records exist';

    -- If any assertion fails, an exception will be raised automatically
END $$;

COMMIT;