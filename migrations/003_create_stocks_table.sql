CREATE TABLE IF NOT EXISTS stocks
(
    id                  INTEGER        NULL,
    code                TEXT           NOT NULL,
    name                TEXT           NOT NULL,
    listing_date        TEXT           NOT NULL,
    delisting_date      TEXT,
    shares              INTEGER        NOT NULL DEFAULT 0,
    board               TEXT           NOT NULL DEFAULT 'MAIN',
    total_employees     TEXT,
    annual_dividend     REAL,
    general_information TEXT,
    founding_date       TEXT,
    sector_id           INTEGER,
    sub_sector_id       INTEGER,
    industry_id         INTEGER,
    sub_industry_id     INTEGER,
    last_modified       TEXT           NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (code),
    CONSTRAINT unique_code_name UNIQUE (code, name)
);

CREATE INDEX IF NOT EXISTS idx_stocks_code ON stocks (code);
CREATE INDEX IF NOT EXISTS idx_stocks_name ON stocks (name);
CREATE INDEX IF NOT EXISTS idx_stocks_board ON stocks (board);
