ALTER TABLE "stocks" 
  ADD CONSTRAINT "stocks_recruitments_stocks" 
    FOREIGN KEY("recruitment_id") 
    REFERENCES "recruitments"("id") 
    ON DELETE CASCADE, 
  ADD CONSTRAINT "stocks_users_stocks" 
    FOREIGN KEY("user_id") 
    REFERENCES "users"("id") 
    ON DELETE CASCADE;
