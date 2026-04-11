CREATE TABLE IF NOT EXISTS idxstock.sectors
(
    id             INT            PRIMARY KEY,
    code           VARCHAR(10)    NOT NULL,
    name           VARCHAR(200)   NOT NULL,
    name_en        VARCHAR(200),
    description    TEXT,
    last_modified  TIMESTAMPTZ    NOT NULL DEFAULT now()
);

-- Index for faster lookup by code
CREATE INDEX IF NOT EXISTS idx_sectors_code ON idxstock.sectors (code);

ALTER TABLE idxstock.sectors OWNER TO pakaiwa_app;
