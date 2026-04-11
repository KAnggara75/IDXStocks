DROP TABLE IF EXISTS idxstock.stocks CASCADE;

CREATE TABLE IF NOT EXISTS idxstock.stocks
(
    code           VARCHAR(10)    NOT NULL,
    name           VARCHAR(200)   NOT NULL,
    listing_date   DATE           NOT NULL,
    delisting_date DATE,
    shares         BIGINT         NOT NULL,
    board          idxstock.board NOT NULL DEFAULT 'Main',
    last_modified  TIMESTAMPTZ    NOT NULL DEFAULT now(),
    PRIMARY KEY (code),
    CONSTRAINT unique_code_name UNIQUE (code, name)
);

CREATE INDEX IF NOT EXISTS idx_stocks_code_like ON idxstock.stocks (code text_pattern_ops);

CREATE INDEX IF NOT EXISTS idx_stocks_name_like ON idxstock.stocks (name text_pattern_ops);

CREATE INDEX IF NOT EXISTS idx_stocks_board ON idxstock.stocks (board);

alter table idxstock.stocks
    owner to kanggara;