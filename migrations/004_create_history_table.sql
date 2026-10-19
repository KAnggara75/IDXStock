CREATE TABLE IF NOT EXISTS history (
    code                  TEXT NOT NULL,
    date                  TEXT NOT NULL,
    previous              REAL,
    open_price            REAL,
    first_trade           REAL,
    high                  REAL,
    low                   REAL,
    close                 REAL,
    change                REAL,
    volume                REAL,
    value                 REAL,
    frequency             REAL,
    index_individual      REAL,
    offer                 REAL,
    offer_volume          REAL,
    bid                   REAL,
    bid_volume            REAL,
    listed_shares         REAL,
    tradeble_shares       REAL,
    weight_for_index      REAL,
    foreign_sell          REAL,
    foreign_buy           REAL,
    delisting_date        TEXT,
    non_regular_volume    REAL,
    non_regular_value     REAL,
    non_regular_frequency REAL,
    last_modified         TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (code, date),
    CONSTRAINT fk_history_stock_code FOREIGN KEY (code) REFERENCES stocks (code) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_history_date ON history (date);
CREATE INDEX IF NOT EXISTS idx_history_date_desc ON history (date DESC);
CREATE INDEX IF NOT EXISTS idx_history_covering ON history (code, date DESC);
