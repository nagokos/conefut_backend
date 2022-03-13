CREATE TABLE IF NOT EXISTS "applicants"(
  "id" varchar UNIQUE NOT NULL, 
  "management_status" varchar NOT NULL DEFAULT 'backlog', 
  "recruitment_id" varchar NULL, 
  "user_id" varchar NULL, 
  "created_at" timestamp with time zone NOT NULL, 
  "updated_at" timestamp with time zone NOT NULL,
  PRIMARY KEY("id")
);
CREATE UNIQUE INDEX IF NOT EXISTS "applicant_user_id_recruitment_id" ON "applicants"("user_id", "recruitment_id");
CREATE INDEX IF NOT EXISTS "applicant_user_id" ON "applicants"("user_id");
CREATE INDEX IF NOT EXISTS "applicant_recruitment_id" ON "applicants"("recruitment_id");