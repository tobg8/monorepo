-- Verify bearwise:insert_indicators_004 on pg

BEGIN;

DO $$
    DECLARE
        bollinger_bands_exists BOOLEAN;
        candle_exists BOOLEAN;
        commodity_channel_index_exists BOOLEAN;
        chaikin_money_flow_exists BOOLEAN;
        exponential_moving_average_exists BOOLEAN;
        hull_moving_average_exists BOOLEAN;
        moving_average_exists BOOLEAN;
        money_flow_index_exists BOOLEAN;
        momentum_exists BOOLEAN;
        parabolic_SAR_exists BOOLEAN;
        relative_strength_index_exists BOOLEAN;
        standard_deviation_exists BOOLEAN;
        true_range_exists BOOLEAN;
        TRIX_exists BOOLEAN;
        typical_price_exists BOOLEAN;
        ultimate_oscillator_exists BOOLEAN;
        volume_weighted_average_price_exists BOOLEAN;
    BEGIN
        -- Check existence for each label
        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'bollinger bands'
        ) INTO bollinger_bands_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'candle'
        ) INTO candle_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'commodity channel index'
        ) INTO commodity_channel_index_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'chaikin money flow'
        ) INTO chaikin_money_flow_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'exponential moving average'
        ) INTO exponential_moving_average_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'hull moving average'
        ) INTO hull_moving_average_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'moving average'
        ) INTO moving_average_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'money flow index'
        ) INTO money_flow_index_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'momentum'
        ) INTO momentum_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'parabolic SAR'
        ) INTO parabolic_SAR_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'relative strength index'
        ) INTO relative_strength_index_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'standard deviation'
        ) INTO standard_deviation_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'true range'
        ) INTO true_range_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'TRIX'
        ) INTO TRIX_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'typical price'
        ) INTO typical_price_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'ultimate oscillator'
        ) INTO ultimate_oscillator_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators
            WHERE label = 'volume weighted average price'
        ) INTO volume_weighted_average_price_exists;

        -- Assert that all required records exist
        ASSERT bollinger_bands_exists = true, 'Record with label = ''bollinger bands'' does not exist';
        ASSERT candle_exists = true, 'Record with label = ''candle'' does not exist';
        ASSERT commodity_channel_index_exists = true, 'Record with label = ''commodity channel index'' does not exist';
        ASSERT chaikin_money_flow_exists = true, 'Record with label = ''chaikin money flow'' does not exist';
        ASSERT exponential_moving_average_exists = true, 'Record with label = ''exponential moving average'' does not exist';
        ASSERT hull_moving_average_exists = true, 'Record with label = ''hull moving average'' does not exist';
        ASSERT moving_average_exists = true, 'Record with label = ''moving average'' does not exist';
        ASSERT money_flow_index_exists = true, 'Record with label = ''money flow index'' does not exist';
        ASSERT momentum_exists = true, 'Record with label = ''momentum'' does not exist';
        ASSERT parabolic_SAR_exists = true, 'Record with label = ''parabolic SAR'' does not exist';
        ASSERT relative_strength_index_exists = true, 'Record with label = ''relative strength index'' does not exist';
        ASSERT standard_deviation_exists = true, 'Record with label = ''standard deviation'' does not exist';
        ASSERT true_range_exists = true, 'Record with label = ''true range'' does not exist';
        ASSERT TRIX_exists = true, 'Record with label = ''TRIX'' does not exist';
        ASSERT typical_price_exists = true, 'Record with label = ''typical price'' does not exist';
        ASSERT ultimate_oscillator_exists = true, 'Record with label = ''ultimate oscillator'' does not exist';
        ASSERT volume_weighted_average_price_exists = true, 'Record with label = ''volume weighted average price'' does not exist';

        -- If all assertions pass, print success
        RAISE NOTICE 'Verification successful: All required records exist';

        -- If any assertion fails, an exception will be raised automatically
    END $$;

COMMIT;