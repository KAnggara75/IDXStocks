DROP TYPE IF EXISTS idxstock.board CASCADE;

CREATE TYPE idxstock.board AS ENUM (
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