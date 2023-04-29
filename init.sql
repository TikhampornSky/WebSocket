CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL UNIQUE,
    "email" varchar NOT NULL UNIQUE,
    "password" varchar NOT NULL
);

CREATE TYPE roomType AS ENUM ('public', 'private');

CREATE TABLE "chatrooms" (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL UNIQUE,
    "clients" BIGINT[] DEFAULT array[]::BIGINT[]
);

ALTER TABLE chatrooms ADD COLUMN category roomType DEFAULT 'public';