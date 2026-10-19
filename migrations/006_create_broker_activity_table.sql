CREATE TABLE IF NOT EXISTS broker_activity (
    broker_code TEXT NOT NULL REFERENCES brokers(code),
    stock_code TEXT NOT NULL REFERENCES stocks(code),
    date TEXT NOT NULL,
    side TEXT NOT NULL CHECK (side IN ('buy', 'sell')),
    lot INTEGER,
    value INTEGER,
    avg_price REAL,
    freq INTEGER,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    PRIMARY KEY (broker_code, date, stock_code, side)
);

CREATE INDEX IF NOT EXISTS idx_broker_act_bd ON broker_activity (broker_code, date);
CREATE INDEX IF NOT EXISTS idx_broker_act_bsd ON broker_activity (broker_code, date, stock_code);
