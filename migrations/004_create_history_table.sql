CREATE TABLE IF NOT EXISTS idxstock.history (
    code                  VARCHAR(10) NOT NULL,
    date                  DATE        NOT NULL,
    previous              NUMERIC,
    open_price            NUMERIC,
    first_trade           NUMERIC,
    high                  NUMERIC,
    low                   NUMERIC,
    close                 NUMERIC,
    change                NUMERIC,
    volume                NUMERIC,
    value                 NUMERIC,
    frequency             NUMERIC,
    index_individual      NUMERIC,
    offer                 NUMERIC,
    offer_volume          NUMERIC,
    bid                   NUMERIC,
    bid_volume            NUMERIC,
    listed_shares         NUMERIC,
    tradeble_shares       NUMERIC,
    weight_for_index      NUMERIC,
    foreign_sell          NUMERIC,
    foreign_buy           NUMERIC,
    delisting_date        DATE,
    non_regular_volume    NUMERIC,
    non_regular_value     NUMERIC,
    non_regular_frequency NUMERIC,
    last_modified         TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (code, date),
    CONSTRAINT fk_history_stock_code FOREIGN KEY (code) REFERENCES idxstock.stocks (code) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_history_code_like ON idxstock.history (code text_pattern_ops);
CREATE INDEX IF NOT EXISTS idx_history_code_fk ON idxstock.history (code);
CREATE INDEX IF NOT EXISTS idx_history_date ON idxstock.history (date);
