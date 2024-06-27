CREATE
KEYSPACE IF NOT EXISTS chat_app
WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

CREATE TABLE IF NOT EXISTS chat_app.users
(
    id
    UUID
    PRIMARY
    KEY,
    username
    TEXT,
    password
    TEXT,
    created_at
    TIMESTAMP
);

CREATE TABLE IF NOT EXISTS chat_app.messages
(
    id
    UUID
    PRIMARY
    KEY,
    sender
    TEXT,
    recipient
    TEXT,
    content
    TEXT,
    timestamp
    TIMESTAMP
);
