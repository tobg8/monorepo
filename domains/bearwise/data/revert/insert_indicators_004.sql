-- Revert bearwise:insert_indicators_004 from pg

BEGIN;

DELETE FROM indicators WHERE label = 'bollinger bands';
DELETE FROM indicators WHERE label = 'candle';
DELETE FROM indicators WHERE label = 'commodity channel index';
DELETE FROM indicators WHERE label = 'chaikin money flow';
DELETE FROM indicators WHERE label = 'exponential moving average';
DELETE FROM indicators WHERE label = 'hull moving average';
DELETE FROM indicators WHERE label = 'moving average';
DELETE FROM indicators WHERE label = 'money flow index';
DELETE FROM indicators WHERE label = 'momentum';
DELETE FROM indicators WHERE label = 'parabolic SAR';
DELETE FROM indicators WHERE label = 'relative strength index';
DELETE FROM indicators WHERE label = 'standard deviation';
DELETE FROM indicators WHERE label = 'true range';
DELETE FROM indicators WHERE label = 'TRIX';
DELETE FROM indicators WHERE label = 'typical price';
DELETE FROM indicators WHERE label = 'ultimate oscillator';
DELETE FROM indicators WHERE label = 'volume weighted average price';

COMMIT;