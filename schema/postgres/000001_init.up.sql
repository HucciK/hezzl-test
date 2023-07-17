CREATE TABLE campaigns (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE items(
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER DEFAULT 1,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    priority SERIAL,
    removed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX items_idx ON items(id, campaign_id, name);
INSERT INTO campaigns(name) VALUES ('Первая запись');