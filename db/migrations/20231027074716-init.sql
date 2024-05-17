
-- +migrate Up

CREATE TABLE "users" (
    "id" UUID PRIMARY KEY NOT NULL,
    "name" TEXT NOT NULL,
    "email" TEXT NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    "role" TEXT NOT NULL,
    "active" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    "deleted_at" TIMESTAMP
);

CREATE TABLE "files" (
    "id" UUID PRIMARY KEY NOT NULL,
    "name" TEXT NOT NULL,
    "type" TEXT NOT NULL,
    "ext" TEXT NOT NULL,
    "url" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    "deleted_at" TIMESTAMP
);

CREATE TABLE outbox_items (
    id TEXT PRIMARY KEY NOT NULL,
    idempotent_key TEXT NOT NULL,
    "version" INT NOT NULL DEFAULT 1,
    "status" TEXT NOT NULL,
    job_type TEXT NOT NULL,
    payload TEXT NOT NULL
);

CREATE INDEX outbox_items_idempotent_key ON outbox_items ("idempotent_key");

-- +migrate Down
DROP INDEX outbox_items_key;
DROP TABLE "outbox_items";
DROP TABLE "files";
DROP TABLE "users";
