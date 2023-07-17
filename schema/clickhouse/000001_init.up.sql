CREATE TABLE items(
    Id INTEGER,
    CampaignId INTEGER,
    Name TEXT NOT NULL,
    Description TEXT NOT NULL,
    Priority INTEGER,
    Removed BOOLEAN NOT NULL DEFAULT FALSE,
    EventTime TIMESTAMP NOT NULL DEFAULT now()
)
ENGINE = MergeTree()
PRIMARY KEY (Id)