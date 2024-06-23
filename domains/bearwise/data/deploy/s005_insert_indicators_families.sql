-- Deploy bearwise:insert_indicators_families_005 to pg

BEGIN;

-- Bollinger Bands
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'bollinger bands'),
    (SELECT id FROM families WHERE label_en = 'overlap studies')
);

-- Candle
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'candle'),
    (SELECT id FROM families WHERE label_en = 'price')
);

-- Commodity Channel Index (CCI)
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'commodity channel index'),
    (SELECT id FROM families WHERE label_en = 'momentum')
);

-- Chaikin Money Flow (CMF)
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'chaikin money flow'),
    (SELECT id FROM families WHERE label_en = 'volume')
);

-- Exponential Moving Average (EMA)
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'exponential moving average'),
    (SELECT id FROM families WHERE label_en = 'overlap studies')
);

-- Hull Moving Average (HMA)
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
(SELECT id FROM indicators WHERE label = 'hull moving average'),
(SELECT id FROM families WHERE label_en = 'overlap studies')
);

-- Moving Average (MA)
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'moving average'),
    (SELECT id FROM families WHERE label_en = 'overlap studies')
);

-- Money Flow Index (MFI) + momentum
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'money flow index'),
    (SELECT id FROM families WHERE label_en = 'momentum')
);
-- Money Flow Index (MFI) + oscillators
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'money flow index'),
    (SELECT id FROM families WHERE label_en = 'oscillators')
);

-- Momentum (MOM)
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'momentum'),
    (SELECT id FROM families WHERE label_en = 'momentum')
);

-- Parabolic SAR (PSAR) + overlap studies
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'parabolic SAR'),
    (SELECT id FROM families WHERE label_en = 'overlap studies')
);
-- Parabolic SAR (PSAR) + trend
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'parabolic SAR'),
    (SELECT id FROM families WHERE label_en = 'trend')
);

-- Relative Strength Index (RSI) + momentum
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'relative strength index'),
    (SELECT id FROM families WHERE label_en = 'momentum')
);
-- Relative Strength Index (RSI) + oscillators
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'relative strength index'),
    (SELECT id FROM families WHERE label_en = 'oscillators')
);

-- Standard Deviation (STDDEV) + statistic functions
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'standard deviation'),
    (SELECT id FROM families WHERE label_en = 'statistic functions')
);
-- Standard Deviation (STDDEV) + volatility
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'standard deviation'),
    (SELECT id FROM families WHERE label_en = 'volatility')
);

-- True Range (TR)
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'true range'),
    (SELECT id FROM families WHERE label_en = 'price')
);

-- TRIX + momentum
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'TRIX'),
    (SELECT id FROM families WHERE label_en = 'momentum')
);
-- TRIX + oscillators
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'TRIX'),
    (SELECT id FROM families WHERE label_en = 'oscillators')
);

-- Typical Price (TYPPRICE)
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'typical price'),
    (SELECT id FROM families WHERE label_en = 'price')
);

-- Ultimate Oscillator (ULTOSC) + momentum
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'ultimate oscillator'),
    (SELECT id FROM families WHERE label_en = 'momentum')
);
-- Ultimate Oscillator (ULTOSC) + oscillators
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'ultimate oscillator'),
    (SELECT id FROM families WHERE label_en = 'oscillators')
);


-- Volume Weighted Average Price (VWAP) + overlap studies
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'volume weighted average price'),
    (SELECT id FROM families WHERE label_en = 'overlap studies')
);
-- Volume Weighted Average Price (VWAP) + volume
INSERT INTO indicators_families (indicator_id, family_id) VALUES (
    (SELECT id FROM indicators WHERE label = 'volume weighted average price'),
    (SELECT id FROM families WHERE label_en = 'volume')
);

COMMIT;
