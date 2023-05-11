CREATE TABLE IF NOT EXISTS "group_events" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "groupId" INTEGER NOT NULL,
    "createdBy" INTEGER NOT NULL,
    "title" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "dateStart" TIMESTAMP NOT NULL,
    "dateEnd" TIMESTAMP NOT NULL,
    "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "creatorName" TEXT NOT NULL,
    "img" TEXT NOT NULL,
    "groupName" TEXT NOT NULL,
    FOREIGN KEY (groupId) REFERENCES groups (id) ON DELETE CASCADE
);