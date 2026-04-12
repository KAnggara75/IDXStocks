DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type t JOIN pg_namespace n ON t.typnamespace = n.oid WHERE t.typname = 'board' AND n.nspname = 'idxstock') THEN
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
    END IF;
END $$;
