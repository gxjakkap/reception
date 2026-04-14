-- Modify "guilds" table
ALTER TABLE "guilds" ADD COLUMN "deleted_at" timestamptz NULL;
-- Create index "idx_guilds_deleted_at" to table: "guilds"
CREATE INDEX "idx_guilds_deleted_at" ON "guilds" ("deleted_at");
