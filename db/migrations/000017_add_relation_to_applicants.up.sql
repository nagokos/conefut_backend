ALTER TABLE "applicants" 
  ADD CONSTRAINT "applicants_recruitments_applicants" 
    FOREIGN KEY("recruitment_id") 
    REFERENCES "recruitments"("id") 
    ON DELETE CASCADE, 
  ADD CONSTRAINT "applicants_users_applicants" 
    FOREIGN KEY("user_id") 
    REFERENCES "users"("id")
    ON DELETE CASCADE;
