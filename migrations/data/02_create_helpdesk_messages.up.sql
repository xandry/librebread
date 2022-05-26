CREATE TABLE helpdesk_messages (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "created_at" datetime NOT NULL,
    "title" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "type_id" INTEGER NOT NULL,
    "priority_id" INTEGER NOT NULL,
    "department_id" INTEGER NOT NULL
);
CREATE INDEX helpdesk_messages_created_at_IDX ON helpdesk_messages (created_at DESC);
