-- Deploy bearwise:insert_indicators_004 to pg

BEGIN;

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'bbands',

    'The Bollinger Bands is a popular technical indicator used by traders in financial markets to analyze price volatility and potential trend reversals. ' ||
    'It consists of three lines plotted on a price chart: the upper band, the middle band, and the lower band. The middle band represents the moving average of the underlying asset’s price over a specified period, typically 20 candles (adjustable using the optional period parameter).
    The upper and lower bands are placed a certain number of standard deviations away from the middle band, usually two standard deviations. The result returned in the given format provides the numerical values for the upper band, middle band, and lower band.' ||
    'Traders can use Bollinger Bands to identify periods of low volatility (when the bands are narrow) and high volatility (when the bands widen). When the price touches or moves beyond the upper band, it suggests that the asset may be overbought, potentially indicating a selling opportunity. Conversely, when the price reaches or falls below the lower band, it suggests oversold conditions, potentially signaling a buying opportunity. ' ||
    'Traders often look for price reversals or bounces off the bands to make trading decisions and assess potential price targets and stop-loss levels.',

    'bollinger bands',
    (SELECT id FROM types WHERE label = 'bollinger')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'candle',
    'Returns the latest closed candle',
    'candle',
    (SELECT id FROM types WHERE label = 'candle')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'cci',

    'Developed by Donald Lambert, the Commodity Channel Index​ (CCI) is a momentum-based oscillator used to help determine when an asset is reaching overbought or oversold conditions.',

    'commodity channel index',
    (SELECT id FROM types WHERE label = 'single')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'cmf',
    '',
    'chaikin money flow',
    (SELECT id FROM types WHERE label = 'single')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'ema',

    'The Exponential Moving Average (EMA) is a widely used technical indicator in financial markets that helps to smooth price data and identify trends more effectively.' ||
    'The EMA, a type of moving average, places greater weight on the most recent price data, making it more responsive to new information compared to the Simple Moving Average (SMA). ' ||
    'This characteristic allows traders and analysts to better capture short-term price movements and trends. The Exponential Moving Average is calculated by applying a multiplier to the most recent closing price, adding it to the previous period’s EMA, and then dividing by the total number of periods. Commonly used periods for EMA calculations include 12, 20, and 50 days.' ||
    ' Traders utilize this indicator to identify trend direction, potential entry and exit points, and to generate trading signals.',

    'exponential moving average',
    (SELECT id FROM types WHERE label = 'single')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'hma',

    'Unlike the conventional moving averages, the Hull Moving Average incorporates weighted moving averages and a square root factor, resulting in a smoother curve that reacts more swiftly to price changes. This enhanced responsiveness makes the HMA particularly useful for traders seeking to capture trends with minimal delay, providing a valuable alternative to the more typical simple and exponential moving averages in the realm of technical analysis. ' ||
    'We offer you the added benefit of adjusting the period parameter in the Hull Moving Average (HMA) to tailor your technical analysis precisely to your liking. This feature enables you to customize the indicator based on your specific trading preferences and timeframes. ' ||
    'By tweaking the period, you have control over the number of candles or periods considered in the HMA calculation. Opting for a shorter period provides you with a more responsive HMA, ideal for capturing rapid trend changes. Conversely, choosing a longer period results in a smoother HMA that excels at identifying prolonged trends, catering to those with a patient trading approach. ' ||
    'Fine-tune the period parameter and enhance your ability to adapt the Hull Moving Average to different market conditions so that it aligns with your preferred trading style. ',

    'hull moving average',
    (SELECT id FROM types WHERE label = 'single')
    );

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'ma',

    'This indicators returns the moving average (MA) values',

    'moving average',
    (SELECT id FROM types WHERE label = 'single')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'mfi',

    'The Money Flow Index (MFI) is a technical indicator used in financial markets to assess the strength and momentum of a price movement. ' ||
    'Developed by Gene Quong and Avrum Soudack, MFI combines price and volume data to provide insights into the buying and selling pressure behind a security.The MFI is an oscillator and is calculated using a formula that takes into account the ratio of positive to negative money flow and scales it to a value between 0 and 100. A high MFI suggests a potentially overbought condition, indicating that a reversal or correction may be imminent, while a low MFI may signal an oversold condition, suggesting a potential buying opportunity.',

    'money flow index',
    (SELECT id FROM types WHERE label = 'single')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'mom',

    'The Momentum Indicator (MOM) is a leading indicator measuring a security’s rate-of-change. ' ||
    'It compares the current price with the previous price from a number of periods ago.The ongoing plot forms an oscillator that moves above and below 0. It is a fully unbounded oscillator and has no lower or upper limit. Bullish and bearish interpretations are found by looking for divergences, centerline crossovers and extreme readings.' ||
    'The indicator is often used in combination with other signals.',

    'momentum',
    (SELECT id FROM types WHERE label = 'single')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'psar',

    'The Parabolic SAR (Stop and Reverse) is a popular technical indicator used in financial markets to identify potential trend reversals. ' ||
    'Developed by J. Welles Wilder Jr., it utilizes a series of dots plotted on a price chart to highlight potential entry and exit points.' ||
    'The Parabolic SAR dots are positioned above or below the price, depending on the direction of the prevailing trend. When the dots are below the price, it indicates an uptrend, while dots above the price indicate a downtrend. The dots gradually adjust their position as the price evolves, creating a parabolic shape. ' ||
    'Traders often use the Parabolic SAR as a tool for setting stop-loss orders and trailing stops, as it dynamically adapts to market conditions and can help identify potential trend reversals in a timely manner.',

    'parabolic SAR',
    (SELECT id FROM types WHERE label = 'single')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'rsi',

    'The Relative Strength Index (RSI) is a popular technical indicator used in financial markets to assess the momentum and strength of a price trend. We provide the best RSI API available.' ||
    'Developed by J. Welles Wilder, the RSI is a momentum oscillator that measures the speed and change of price movements.' ||
    'It is typically applied to a price chart, oscillating between 0 and 100, with readings above 70 considered overbought and readings below 30 indicating oversold conditions. The Relative Strength Index is calculated based on the average gains and losses over a specified period, commonly 14 days.' ||
    'Traders and analysts use this indicator to identify potential trend reversals, overbought or oversold conditions, and to generate signals for buying or selling assets.',

    'relative strength index',
    (SELECT id FROM types WHERE label = 'single')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'stddev',

    'The standard deviation technical indicator is a statistical measure that quantifies the dispersion of a set of data points from its mean (average).' ||
    'In the context of financial markets, it is often used to assess the volatility or variability of price movements.',

    'standard deviation',
    (SELECT id FROM types WHERE label = 'single')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'tr',

    'True range measures the daily range plus any gap from the closing price of the preceding day.',

    'true range',
    (SELECT id FROM types WHERE label = 'single')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'trix',

    'The Triple Exponential Average (TRIX) is a momentum indicator used by technical traders that shows the percentage change in a triple exponentially smoothed moving average.' ||
    'When it is applied to triple smoothing of moving averages, it is designed to filter out price movements that are considered insignificant or unimportant.',

    'TRIX',
    (SELECT id FROM types WHERE label = 'single')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'typprice',

    '',

    'typical price',
    (SELECT id FROM types WHERE label = 'single')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'ultosc',

    '',

    'ultimate oscillator',
    (SELECT id FROM types WHERE label = 'single')
);

INSERT INTO indicators (symbol, description, label, type_id) VALUES (
    'vwap',

    'The volume weighted average price (VWAP) is a trading benchmark used by traders that gives the average price a security has traded at throughout the day, based on both volume and price.' ||
    'It is important because it provides traders with insight into both the trend and value of a security.',

    'volume weighted average price',
    (SELECT id FROM types WHERE label = 'single')
);


COMMIT;
