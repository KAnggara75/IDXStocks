CREATE TABLE IF NOT EXISTS brokers (
    code VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    investor_type VARCHAR,
    total_value BIGINT,
    net_value BIGINT,
    buy_value BIGINT,
    sell_value BIGINT,
    total_volume BIGINT,
    total_frequency BIGINT,
    broker_group VARCHAR, -- Menghindari tipe reserved word 'group'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
