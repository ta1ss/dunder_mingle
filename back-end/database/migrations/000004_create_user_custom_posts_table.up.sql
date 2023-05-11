CREATE TABLE IF NOT EXISTS "user_custom_posts" (
    "postId" INTEGER NOT NULL,
    "targetId" INTEGER NOT NULL,
    FOREIGN KEY (postId) REFERENCES user_posts (id) ON DELETE CASCADE,
    FOREIGN KEY (targetId) REFERENCES users (id) ON DELETE CASCADE
);