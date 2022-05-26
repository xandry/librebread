CREATE TABLE sms_messages (
    "id" BLOB NOT NULL,
    "time" datetime NOT NULL,
    "from" TEXT NOT NULL,
    "to" TEXT NOT NULL,
    "text" TEXT NOT NULL,
    "provider" TEXT NOT NULL,
    CONSTRAINT sms_message_PK PRIMARY KEY (id)
);
CREATE INDEX sms_time_IDX ON sms_messages (time DESC);
