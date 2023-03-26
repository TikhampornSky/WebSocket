CREATE TABLE "chatrooms" (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL UNIQUE,
)