-- Revert bearwise:insert_types_003 from pg

BEGIN;

DELETE FROM types WHERE label = 'bollinger';
DELETE FROM types WHERE label = 'candle';
DELETE FROM types WHERE label = 'single';

COMMIT;
