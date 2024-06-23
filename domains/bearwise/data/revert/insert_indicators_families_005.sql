-- Revert bearwise:insert_indicators_families_005 from pg

BEGIN;

-- Bollinger Bands
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'bollinger bands')
  AND family_id = (SELECT id FROM families WHERE label_en = 'overlap studies');

-- Candle
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'candle')
  AND family_id = (SELECT id FROM families WHERE label_en = 'price');

-- Commodity Channel Index (CCI)
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'commodity channel index')
  AND family_id = (SELECT id FROM families WHERE label_en = 'momentum');

-- Chaikin Money Flow (CMF)
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'chaikin money flow')
  AND family_id = (SELECT id FROM families WHERE label_en = 'volume');

-- Exponential Moving Average (EMA)
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'exponential moving average')
  AND family_id = (SELECT id FROM families WHERE label_en = 'overlap studies');

-- Hull Moving Average (HMA)
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'hull moving average')
  AND family_id = (SELECT id FROM families WHERE label_en = 'overlap studies');

-- Moving Average (MA)
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'moving average')
  AND family_id = (SELECT id FROM families WHERE label_en = 'overlap studies');

-- Money Flow Index (MFI) + momentum
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'money flow index')
  AND family_id = (SELECT id FROM families WHERE label_en = 'momentum');
-- Money Flow Index (MFI) + oscillators
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'money flow index')
  AND family_id = (SELECT id FROM families WHERE label_en = 'oscillators');

-- Momentum (MOM)
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'momentum')
  AND family_id = (SELECT id FROM families WHERE label_en = 'momentum');

-- Parabolic SAR (PSAR) + overlap studies
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'parabolic SAR')
  AND family_id = (SELECT id FROM families WHERE label_en = 'overlap studies');
-- Parabolic SAR (PSAR) + trend
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'parabolic SAR')
  AND family_id = (SELECT id FROM families WHERE label_en = 'trend');

-- Relative Strength Index (RSI) + momentum
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'relative strength index')
  AND family_id = (SELECT id FROM families WHERE label_en = 'momentum');
-- Relative Strength Index (RSI) + oscillators
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'relative strength index')
  AND family_id = (SELECT id FROM families WHERE label_en = 'oscillators');

-- Standard Deviation (STDDEV) + statistic functions
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'standard deviation')
  AND family_id = (SELECT id FROM families WHERE label_en = 'statistic functions');
-- Standard Deviation (STDDEV) + volatility
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'standard deviation')
  AND family_id = (SELECT id FROM families WHERE label_en = 'volatility');

-- True Range (TR)
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'true range')
  AND family_id = (SELECT id FROM families WHERE label_en = 'price');

-- TRIX + momentum
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'TRIX')
  AND family_id = (SELECT id FROM families WHERE label_en = 'momentum');
-- TRIX + oscillators
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'TRIX')
  AND family_id = (SELECT id FROM families WHERE label_en = 'oscillators');

-- Typical Price (TYPPRICE)
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'typical price')
  AND family_id = (SELECT id FROM families WHERE label_en = 'price');

-- Ultimate Oscillator (ULTOSC) + momentum
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'ultimate oscillator')
  AND family_id = (SELECT id FROM families WHERE label_en = 'momentum');
-- Ultimate Oscillator (ULTOSC) + oscillators
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'ultimate oscillator')
  AND family_id = (SELECT id FROM families WHERE label_en = 'oscillators');

-- Volume Weighted Average Price (VWAP) + overlap studies
DELETE FROM indicators_families
WHERE indicator_id = (SELECT id FROM indicators WHERE label = 'volume weighted average price')
  AND family_id = (SELECT id FROM families WHERE label_en = 'overlap studies');
-- Volume Weighted Average Price (VWAP) + volume
DELETE FROM indicators_families
  AND family_id = (SELECT id FROM families WHERE label_en = 'volume');

COMMIT;
