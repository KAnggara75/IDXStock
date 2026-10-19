CREATE TABLE IF NOT EXISTS industry
(
    id             INTEGER        PRIMARY KEY,
    name           TEXT           NOT NULL,
    last_modified  TEXT           NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS sub_industry
(
    id             INTEGER        PRIMARY KEY,
    name           TEXT           NOT NULL,
    industry_id    INTEGER        NOT NULL REFERENCES industry(id),
    last_modified  TEXT           NOT NULL DEFAULT (datetime('now'))
);
