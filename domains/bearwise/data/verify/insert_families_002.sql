-- Verify bearwise:insert_families_002 on pg

BEGIN;

DO $$
    DECLARE
        band_exists BOOLEAN;
        breakout_exists BOOLEAN;
        math_operators_exists BOOLEAN;
        math_transform_exists BOOLEAN;
        momentum_exists BOOLEAN;
        oscillators_exists BOOLEAN;
        overlap_studies_exists BOOLEAN;
        pattern_recognition_exists BOOLEAN;
        price_exists BOOLEAN;
        sentiment_exists BOOLEAN;
        statistic_functions_exists BOOLEAN;
        support_resistance_exists BOOLEAN;
        trend_exists BOOLEAN;
        volatility_exists BOOLEAN;
        volume_exists BOOLEAN;
    BEGIN
        -- Check existence for each label_en
        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'bands'
        ) INTO band_exists;

        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'breakouts'
        ) INTO breakout_exists;

        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'math operators'
        ) INTO math_operators_exists;

        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'math transform'
        ) INTO math_transform_exists;

        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'momentum'
        ) INTO momentum_exists;

        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'oscillators'
        ) INTO oscillators_exists;

        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'overlap studies'
        ) INTO overlap_studies_exists;

        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'pattern recognition'
        ) INTO pattern_recognition_exists;

        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'price'
        ) INTO price_exists;

        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'sentiment'
        ) INTO sentiment_exists;

        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'statistic functions'
        ) INTO statistic_functions_exists;

        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'support & resistance'
        ) INTO support_resistance_exists;

        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'trend'
        ) INTO trend_exists;

        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'volatility'
        ) INTO volatility_exists;

        SELECT EXISTS (
            SELECT 1
            FROM families
            WHERE label_en = 'volume'
        ) INTO volume_exists;

        -- Assert that all required records exist
        ASSERT band_exists = true, 'Record with label_en = ''wowowowo'' does not exist';
        ASSERT breakout_exists = true, 'Record with label_en = ''breakouts'' does not exist';
        ASSERT math_operators_exists = true, 'Record with label_en = ''math operators'' does not exist';
        ASSERT math_transform_exists = true, 'Record with label_en = ''math transform'' does not exist';
        ASSERT momentum_exists = true, 'Record with label_en = ''momentum'' does not exist';
        ASSERT oscillators_exists = true, 'Record with label_en = ''oscillators'' does not exist';
        ASSERT overlap_studies_exists = true, 'Record with label_en = ''overlap studies'' does not exist';
        ASSERT pattern_recognition_exists = true, 'Record with label_en = ''pattern recognition'' does not exist';
        ASSERT price_exists = true, 'Record with label_en = ''price'' does not exist';
        ASSERT sentiment_exists = true, 'Record with label_en = ''sentiment'' does not exist';
        ASSERT statistic_functions_exists = true, 'Record with label_en = ''statistic functions'' does not exist';
        ASSERT support_resistance_exists = true, 'Record with label_en = ''support & resistance'' does not exist';
        ASSERT trend_exists = true, 'Record with label_en = ''trend'' does not exist';
        ASSERT volatility_exists = true, 'Record with label_en = ''volatility'' does not exist';
        ASSERT volume_exists = true, 'Record with label_en = ''volume'' does not exist';

        -- If all assertions pass, print success
        RAISE NOTICE 'Verification successful: All required records exist';

        -- If any assertion fails, an exception will be raised automatically
    END $$;

COMMIT;