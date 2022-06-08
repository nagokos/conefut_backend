ALTER TABLE "stocks" 
  ADD CONSTRAINT "stocks_recruitment_id" 
    FOREIGN KEY("recruitment_id") 
    REFERENCES "recruitments"("id") 
    ON DELETE CASCADE, 
  ADD CONSTRAINT "stocks_user_id" 
    FOREIGN KEY("user_id") 
    REFERENCES "users"("id") 
    ON DELETE CASCADE;
