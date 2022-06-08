CREATE TABLE IF NOT EXISTS "recruitment_tags"(
  "id" VARCHAR UNIQUE NOT NULL,
  "recruitment_id" VARCHAR NULL, 
  "tag_id" VARCHAR NULL, 
  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL, 
  "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL, 
  PRIMARY KEY("id"),
  UNIQUE("recruitment_id", "tag_id")
);
CREATE INDEX IF NOT EXISTS "recruitmenttag_tag_id" ON "recruitment_tags"("tag_id");
CREATE INDEX IF NOT EXISTS "recruitmenttag_recruitment_id" ON "recruitment_tags"("recruitment_id");