-- Create "guilds" table
CREATE TABLE "guilds" (
  "id" text NOT NULL,
  "guild_name" text NOT NULL,
  "joined_at" timestamptz NULL,
  "prefix" text NULL DEFAULT '&&',
  "settings" jsonb NULL DEFAULT '{}',
  PRIMARY KEY ("id")
);
-- Create "pendings" table
CREATE TABLE "pendings" (
  "id" bigserial NOT NULL,
  "guild_id" text NOT NULL,
  "user_id" text NOT NULL,
  "type" text NOT NULL,
  "next" text NOT NULL,
  "created_at" timestamptz NULL,
  "expired_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_pendings_guild_id" to table: "pendings"
CREATE INDEX "idx_pendings_guild_id" ON "pendings" ("guild_id");
-- Create index "idx_pendings_user_id" to table: "pendings"
CREATE INDEX "idx_pendings_user_id" ON "pendings" ("user_id");
-- Create "histories" table
CREATE TABLE "histories" (
  "id" bigserial NOT NULL,
  "guild_id" text NOT NULL,
  "user_id" text NOT NULL,
  "mod_id" text NOT NULL,
  "type" text NOT NULL,
  "reason" text NULL DEFAULT 'No reason provided',
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_guilds_history" FOREIGN KEY ("guild_id") REFERENCES "guilds" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_histories_guild_id" to table: "histories"
CREATE INDEX "idx_histories_guild_id" ON "histories" ("guild_id");
-- Create index "idx_histories_user_id" to table: "histories"
CREATE INDEX "idx_histories_user_id" ON "histories" ("user_id");
-- Create "reaction_roles" table
CREATE TABLE "reaction_roles" (
  "id" bigserial NOT NULL,
  "guild_id" text NOT NULL,
  "message_id" text NOT NULL,
  "emoji" text NOT NULL,
  "role_id" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_guilds_reaction_roles" FOREIGN KEY ("guild_id") REFERENCES "guilds" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_reaction_roles_guild_id" to table: "reaction_roles"
CREATE INDEX "idx_reaction_roles_guild_id" ON "reaction_roles" ("guild_id");
-- Create index "idx_reaction_roles_message_id" to table: "reaction_roles"
CREATE INDEX "idx_reaction_roles_message_id" ON "reaction_roles" ("message_id");
