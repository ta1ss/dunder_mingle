CREATE TABLE IF NOT EXISTS "users" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "email" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "firstName" TEXT NOT NULL,
    "lastName" TEXT NOT NULL,
    "dateOfBirth" TIMESTAMP NOT NULL,
    "img" TEXT NOT NULL DEFAULT "default_profile.png",
    "nickname" TEXT NOT NULL DEFAULT "",
    "about" TEXT NOT NULL DEFAULT "",
    "profilePublic" INTEGER NOT NULL DEFAULT 1,
    "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);