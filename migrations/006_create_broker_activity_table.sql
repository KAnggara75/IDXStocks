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
    PRIMARY KEY (broker_code, stock_code, date, side)
);

-- Indeks kombinasi untuk mempercepat pencarian (sering dipakai):
CREATE INDEX IF NOT EXISTS idx_broker_act_bd ON idxstock.broker_activity (broker_code, date);
CREATE INDEX IF NOT EXISTS idx_broker_act_bsd ON idxstock.broker_activity (broker_code, stock_code, date);
