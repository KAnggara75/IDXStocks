CREATE TABLE IF NOT EXISTS idxstock.sector
(
    id             INT            PRIMARY KEY,
    code           VARCHAR(50),
    name           VARCHAR(200)   NOT NULL,
    name_en        VARCHAR(200),
    description    TEXT,
    last_modified  TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS idxstock.sub_sector
(
    id             INT            PRIMARY KEY,
    code           VARCHAR(50),
    name           VARCHAR(200)   NOT NULL,
    name_en        VARCHAR(200),
    description    TEXT,
    sector_id      INT            NOT NULL REFERENCES idxstock.sector(id),
    last_modified  TIMESTAMPTZ    NOT NULL DEFAULT now()
);
