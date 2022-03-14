CREATE TABLE IF NOT EXISTS "recruitment_tags"(
  "id" varchar UNIQUE NOT NULL,
  "recruitment_id" varchar NULL, 
  "tag_id" varchar NULL, 
  "created_at" timestamp with time zone NOT NULL, 
  "updated_at" timestamp with time zone NOT NULL, 
  PRIMARY KEY("id")
);
CREATE UNIQUE INDEX IF NOT EXISTS "recruitmenttag_recruitment_id_tag_id" ON "recruitment_tags"("recruitment_id", "tag_id");
CREATE INDEX IF NOT EXISTS "recruitmenttag_tag_id" ON "recruitment_tags"("tag_id");
CREATE INDEX IF NOT EXISTS "recruitmenttag_recruitment_id" ON "recruitment_tags"("recruitment_id");