ALTER TABLE "messages"
  ADD COLUMN "applicant_id" VARCHAR NOT NULL,
  ADD FOREIGN KEY("applicant_id") 
    REFERENCES "applicants"("id") 
    ON DELETE CASCADE;