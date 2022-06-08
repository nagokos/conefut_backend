CREATE TABLE IF NOT EXISTS "applicants"(
  "id" VARCHAR UNIQUE NOT NULL, 
  "management_status" VARCHAR NOT NULL DEFAULT 'backlog', 
  "recruitment_id" VARCHAR NULL, 
  "user_id" VARCHAR NULL, 
  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL, 
  "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL,
  PRIMARY KEY("id"),
  UNIQUE("user_id", "recruitment_id")
);
CREATE INDEX ON "applicants"("user_id");
CREATE INDEX ON "applicants"("recruitment_id");