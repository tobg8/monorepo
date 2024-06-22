-- Verify bearwise:insert_families_002 on pg

BEGIN;

SELECT label_en,
       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'bands'
       ) AS bands_exist,

       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'breakouts'
       ) AS breakouts_exist,

       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'math operators'
       ) AS math_operators_exist,

       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'math transform'
       ) AS math_transform_exist,

       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'momentum'
       ) AS momentum_exist,

       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'oscillators'
       ) AS oscillators_exist,

       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'overlap studies'
       ) AS overlap_studies_exist,

       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'pattern recognition'
       ) AS pattern_recognition_exist,

       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'price'
       ) AS price_exist,

       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'sentiment'
       ) AS sentiment_exist,

       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'statistic functions'
       ) AS statistic_functions_exist,

       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'support & resistance'
       ) AS support_resistance_exist,

       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'trend'
       ) AS trend_exist,

       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'volatility'
       ) AS volatility_exist,

       EXISTS (
           SELECT 1
           FROM families
           WHERE label_en = 'volume'
       ) AS volume_exist

FROM families
         LIMIT 1;

COMMIT;