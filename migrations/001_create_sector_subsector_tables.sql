CREATE TABLE IF NOT EXISTS sector
(
    id             INTEGER        PRIMARY KEY,
    code           TEXT,
    name           TEXT           NOT NULL,
    name_en        TEXT,
    description    TEXT,
    last_modified  TEXT           NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS sub_sector
(
    id             INTEGER        PRIMARY KEY,
    code           TEXT,
    name           TEXT           NOT NULL,
    name_en        TEXT,
    description    TEXT,
    sector_id      INTEGER        NOT NULL REFERENCES sector(id),
    last_modified  TEXT           NOT NULL DEFAULT (datetime('now'))
);
