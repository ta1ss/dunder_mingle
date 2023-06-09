CREATE TABLE IF NOT EXISTS "followers" (
    "userId" INTEGER NOT NULL,
    "followerId" INTEGER NOT NULL,
    "requested" INTEGER NOT NULL DEFAULT 0,
    "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (userId) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (followerId) REFERENCES users (id) ON DELETE CASCADE
);