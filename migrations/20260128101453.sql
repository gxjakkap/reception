-- Modify "pendings" table
ALTER TABLE "pendings" ADD COLUMN "data" jsonb NULL DEFAULT '{}';
