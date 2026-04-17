CREATE TABLE IF NOT EXISTS idxstock.broker_activity (
    broker_code VARCHAR NOT NULL REFERENCES idxstock.brokers(code),
    stock_code VARCHAR NOT NULL REFERENCES idxstock.stocks(code),
    date DATE NOT NULL,
    side VARCHAR(10) NOT NULL CHECK (side IN ('buy', 'sell')),
    lot BIGINT,
    value BIGINT,
    avg_price DECIMAL,
    freq BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (broker_code, date, stock_code, side)
) PARTITION BY RANGE (date);

-- Indeks kombinasi untuk mempercepat pencarian:
CREATE INDEX IF NOT EXISTS idx_broker_act_bd ON idxstock.broker_activity (broker_code, date);
CREATE INDEX IF NOT EXISTS idx_broker_act_bsd ON idxstock.broker_activity (broker_code, date, stock_code);

-- Default partition untuk menampung data yang belum memiliki partisi spesifik
CREATE TABLE IF NOT EXISTS idxstock.broker_activity_default
PARTITION OF idxstock.broker_activity DEFAULT;

-- Contoh pembuatan partisi bulanan manual (Opsional: Bisa ditambahkan sesuai kebutuhan range data)
-- CREATE TABLE IF NOT EXISTS idxstock.broker_activity_y2026m04
-- PARTITION OF idxstock.broker_activity
-- FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');
