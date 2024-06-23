-- Verify bearwise:insert_indicators_families_005 on pg

BEGIN;

DO $$
    DECLARE
        bollinger_bands_overlap_exists BOOLEAN;
        candle_price_exists BOOLEAN;
        cci_momentum_exists BOOLEAN;
        cmf_volume_exists BOOLEAN;
        ema_overlap_exists BOOLEAN;
        hma_overlap_exists BOOLEAN;
        ma_overlap_exists BOOLEAN;
        mfi_momentum_exists BOOLEAN;
        mfi_oscillators_exists BOOLEAN;
        mom_momentum_exists BOOLEAN;
        psar_overlap_exists BOOLEAN;
        psar_trend_exists BOOLEAN;
        rsi_momentum_exists BOOLEAN;
        rsi_oscillators_exists BOOLEAN;
        stddev_statistic_exists BOOLEAN;
        stddev_volatility_exists BOOLEAN;
        tr_price_exists BOOLEAN;
        trix_momentum_exists BOOLEAN;
        trix_oscillators_exists BOOLEAN;
        typprice_price_exists BOOLEAN;
        ultosc_momentum_exists BOOLEAN;
        ultosc_oscillators_exists BOOLEAN;
        vwap_overlap_exists BOOLEAN;
        vwap_volume_exists BOOLEAN;
    BEGIN
        -- Check existence for each indicator-family relationship
        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'bollinger bands' AND f.label_en = 'overlap studies'
        ) INTO bollinger_bands_overlap_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'candle' AND f.label_en = 'price'
        ) INTO candle_price_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'commodity channel index' AND f.label_en = 'momentum'
        ) INTO cci_momentum_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'chaikin money flow' AND f.label_en = 'volume'
        ) INTO cmf_volume_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'exponential moving average' AND f.label_en = 'overlap studies'
        ) INTO ema_overlap_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'hull moving average' AND f.label_en = 'overlap studies'
        ) INTO hma_overlap_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'moving average' AND f.label_en = 'overlap studies'
        ) INTO ma_overlap_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'money flow index' AND f.label_en = 'momentum'
        ) INTO mfi_momentum_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'money flow index' AND f.label_en = 'oscillators'
        ) INTO mfi_oscillators_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'momentum' AND f.label_en = 'momentum'
        ) INTO mom_momentum_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'parabolic SAR' AND f.label_en = 'overlap studies'
        ) INTO psar_overlap_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'parabolic SAR' AND f.label_en = 'trend'
        ) INTO psar_trend_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'relative strength index' AND f.label_en = 'momentum'
        ) INTO rsi_momentum_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'relative strength index' AND f.label_en = 'oscillators'
        ) INTO rsi_oscillators_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'standard deviation' AND f.label_en = 'statistic functions'
        ) INTO stddev_statistic_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'standard deviation' AND f.label_en = 'volatility'
        ) INTO stddev_volatility_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'true range' AND f.label_en = 'price'
        ) INTO tr_price_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'TRIX' AND f.label_en = 'momentum'
        ) INTO trix_momentum_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'TRIX' AND f.label_en = 'oscillators'
        ) INTO trix_oscillators_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'typical price' AND f.label_en = 'price'
        ) INTO typprice_price_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'ultimate oscillator' AND f.label_en = 'momentum'
        ) INTO ultosc_momentum_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'ultimate oscillator' AND f.label_en = 'oscillators'
        ) INTO ultosc_oscillators_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'volume weighted average price' AND f.label_en = 'overlap studies'
        ) INTO vwap_overlap_exists;

        SELECT EXISTS (
            SELECT 1
            FROM indicators_families ifa
                     JOIN indicators i ON i.id = ifa.indicator_id
                     JOIN families f ON f.id = ifa.family_id
            WHERE i.label = 'volume weighted average price' AND f.label_en = 'volume'
        ) INTO vwap_volume_exists;

        -- Assert that all required records exist
        ASSERT bollinger_bands_overlap_exists = true, 'Record with label = ''bollinger bands'' and family = ''overlap studies'' does not exist';
        ASSERT candle_price_exists = true, 'Record with label = ''candle'' and family = ''price'' does not exist';
        ASSERT cci_momentum_exists = true, 'Record with label = ''commodity channel index'' and family = ''momentum'' does not exist';
        ASSERT cmf_volume_exists = true, 'Record with label = ''chaikin money flow'' and family = ''volume'' does not exist';
        ASSERT ema_overlap_exists = true, 'Record with label = ''exponential moving average'' and family = ''overlap studies'' does not exist';
        ASSERT hma_overlap_exists = true, 'Record with label = ''hull moving average'' and family = ''overlap studies'' does not exist';
        ASSERT ma_overlap_exists = true, 'Record with label = ''moving average'' and family = ''overlap studies'' does not exist';
        ASSERT mfi_momentum_exists = true, 'Record with label = ''money flow index'' and family = ''momentum'' does not exist';
        ASSERT mfi_oscillators_exists = true, 'Record with label = ''money flow index'' and family = ''oscillators'' does not exist';
        ASSERT mom_momentum_exists = true, 'Record with label = ''momentum'' and family = ''momentum'' does not exist';
        ASSERT psar_overlap_exists = true, 'Record with label = ''parabolic SAR'' and family = ''overlap studies'' does not exist';
        ASSERT psar_trend_exists = true, 'Record with label = ''parabolic SAR'' and family = ''trend'' does not exist';
        ASSERT rsi_momentum_exists = true, 'Record with label = ''relative strength index'' and family = ''momentum'' does not exist';
        ASSERT rsi_oscillators_exists = true, 'Record with label = ''relative strength index'' and family = ''oscillators'' does not exist';
        ASSERT stddev_statistic_exists = true, 'Record with label = ''standard deviation'' and family = ''statistic functions'' does not exist';
        ASSERT stddev_volatility_exists = true, 'Record with label = ''standard deviation'' and family = ''volatility'' does not exist';
        ASSERT tr_price_exists = true, 'Record with label = ''true range'' and family = ''price'' does not exist';
        ASSERT trix_momentum_exists = true, 'Record with label = ''TRIX'' and family = ''momentum'' does not exist';
        ASSERT trix_oscillators_exists = true, 'Record with label = ''TRIX'' and family = ''oscillators'' does not exist';
        ASSERT typprice_price_exists = true, 'Record with label = ''typical price'' and family = ''price'' does not exist';
        ASSERT ultosc_momentum_exists = true, 'Record with label = ''ultimate oscillator'' and family = ''momentum'' does not exist';
        ASSERT ultosc_oscillators_exists = true, 'Record with label = ''ultimate oscillator'' and family = ''oscillators'' does not exist';
        ASSERT vwap_overlap_exists = true, 'Record with label = ''volume weighted average price'' and family = ''overlap studies'' does not exist';
        ASSERT vwap_volume_exists = true, 'Record with label = ''volume weighted average price'' and family = ''volume'' does not exist';

        -- If all assertions pass, print success
        RAISE NOTICE 'Verification successful: All required records exist';

        -- If any assertion fails, an exception will be raised automatically
    END $$;

COMMIT;