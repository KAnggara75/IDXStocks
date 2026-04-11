-- Create board ENUM type
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'board') THEN
        CREATE TYPE board AS ENUM (
            'Main',
            'Ekonomi Baru',
            'Acceleration',
            'Development',
            'Watchlist',
            'A_SERIES',
            'B_SERIES',
            'C_SERIES',
            'PREFEREN'
        );
    END IF;
END $$;

-- Create stocks table
CREATE TABLE IF NOT EXISTS stocks (
    code           VARCHAR(10)    NOT NULL,
    name           VARCHAR(200)   NOT NULL,
    listing_date   DATE           NOT NULL,
    delisting_date DATE,
    shares         BIGINT         NOT NULL,
    board          board          NOT NULL DEFAULT 'Main',
    last_modified  TIMESTAMPTZ    NOT NULL DEFAULT now(),
    PRIMARY KEY (code),
    CONSTRAINT unique_code_name UNIQUE (code, name)
);
