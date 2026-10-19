CREATE TABLE IF NOT EXISTS brokers (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    investor_type TEXT,
    total_value INTEGER,
    net_value INTEGER,
    buy_value INTEGER,
    sell_value INTEGER,
    total_volume INTEGER,
    total_frequency INTEGER,
    broker_group TEXT,
    created_at TEXT DEFAULT (datetime('now'))
);
