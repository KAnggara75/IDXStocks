CREATE TABLE IF NOT EXISTS idxstock.stocks
(
    id                  INT            NULL,
    code                VARCHAR(10)    NOT NULL,
    name                VARCHAR(200)   NOT NULL,
    listing_date        DATE           NOT NULL,
    delisting_date      DATE,
    shares              BIGINT         NOT NULL DEFAULT 0,
    board               idxstock.board NOT NULL DEFAULT 'Main',
    total_employees     VARCHAR(50),
    annual_dividend     NUMERIC,
    general_information TEXT,
    founding_date       VARCHAR(50),
    sector_id           INT,
    sub_sector_id       INT,
    industry_id         INT,
    sub_industry_id     INT,
    last_modified       TIMESTAMPTZ    NOT NULL DEFAULT now(),
    PRIMARY KEY (code),
    CONSTRAINT unique_code_name UNIQUE (code, name)
);

CREATE INDEX IF NOT EXISTS idx_stocks_code_like ON idxstock.stocks (code text_pattern_ops);

CREATE INDEX IF NOT EXISTS idx_stocks_name_like ON idxstock.stocks (name text_pattern_ops);

CREATE INDEX IF NOT EXISTS idx_stocks_board ON idxstock.stocks (board);
