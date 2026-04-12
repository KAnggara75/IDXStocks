CREATE TABLE IF NOT EXISTS idxstock.industry
(
    id             INT            PRIMARY KEY,
    name           VARCHAR(200)   NOT NULL,
    last_modified  TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS idxstock.sub_industry
(
    id             INT            PRIMARY KEY,
    name           VARCHAR(200)   NOT NULL,
    industry_id    INT            NOT NULL REFERENCES idxstock.industry(id),
    last_modified  TIMESTAMPTZ    NOT NULL DEFAULT now()
);

ALTER TABLE idxstock.industry OWNER TO pakaiwa_app;
ALTER TABLE idxstock.sub_industry OWNER TO pakaiwa_app;
