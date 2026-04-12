CREATE TABLE IF NOT EXISTS idxstock.sector
(
    id             INT            PRIMARY KEY,
    name           VARCHAR(200)   NOT NULL,
    last_modified  TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS idxstock.sub_sector
(
    id             INT            PRIMARY KEY,
    name           VARCHAR(200)   NOT NULL,
    sector_id      INT            NOT NULL REFERENCES idxstock.sector(id),
    last_modified  TIMESTAMPTZ    NOT NULL DEFAULT now()
);

ALTER TABLE idxstock.sector OWNER TO pakaiwa_app;
ALTER TABLE idxstock.sub_sector OWNER TO pakaiwa_app;
