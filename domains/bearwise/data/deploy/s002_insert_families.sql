-- Deploy bearwise:insert_families_002 to pg

BEGIN;

INSERT INTO families (label_fr, label_en, label_es, label_it, label_nl, label_de) VALUES
    ('bandes', 'bands', 'bandas', 'bande', 'banden', 'bänder'),
    ('évasions', 'breakouts', 'rupturas', 'rotture', 'uitbraken', 'ausbrüche'),
    ('opérateurs mathématiques', 'math operators', 'operadores matemáticos', 'operatori matematici', 'wiskundige operators', 'mathematische operatoren'),
    ('transformation mathématique', 'math transform', 'transformación matemática', 'trasformazione matematica', 'wiskundige transformatie', 'mathematische transformation'),
    ('momentum', 'momentum', 'momentum', 'momento', 'momentum', 'momentum'),
    ('oscillateurs', 'oscillators', 'osciladores', 'oscillatori', 'oscillatoren', 'oszillatoren'),
    ('études de chevauchement', 'overlap studies', 'estudios de solapamiento', 'studi di sovrapposizione', 'overlap studies', 'überlappungsstudien'),
    ('reconnaissance de formes', 'pattern recognition', 'reconocimiento de patrones', 'riconoscimento dei modelli', 'patroonherkenning', 'mustererkennung'),
    ('prix', 'price', 'precio', 'prezzo', 'prijs', 'preis'),
    ('sentiment', 'sentiment', 'sentimiento', 'sentimento', 'sentiment', 'sentiment'),
    ('fonctions statistiques', 'statistic functions', 'funciones estadísticas', 'funzioni statistiche', 'statistische functies', 'statistische funktionen'),
    ('support et résistance', 'support & resistance', 'soporte y resistencia', 'supporto e resistenza', 'ondersteuning & weerstand', 'unterstützung & widerstand'),
    ('tendance', 'trend', 'tendencia', 'tendenza', 'trend', 'trend'),
    ('volatilité', 'volatility', 'volatilidad', 'volatilità', 'volatiliteit', 'volatilität'),
    ('volume', 'volume', 'volumen', 'volume', 'volume', 'volumen');

COMMIT;
