-- Drop index "idx_reaction_roles_message_id" from table: "reaction_roles"
DROP INDEX "idx_reaction_roles_message_id";
-- Create index "idx_msg_emoji" to table: "reaction_roles"
CREATE INDEX "idx_msg_emoji" ON "reaction_roles" ("message_id", "emoji");
