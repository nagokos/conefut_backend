ALTER TABLE "recruitments" ADD CONSTRAINT "recruitments_users_recruitments" FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE;
