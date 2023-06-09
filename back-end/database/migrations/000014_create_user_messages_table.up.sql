CREATE TABLE IF NOT EXISTS "user_messages" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "senderId" INTEGER NOT NULL, 
    "targetId" INTEGER NOT NULL, 
    "body" INTEGER NOT NULL, 
    "messageRead" INTEGER NOT NULL DEFAULT 0,
    "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (senderId) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (targetId) REFERENCES users (id) ON DELETE CASCADE
);