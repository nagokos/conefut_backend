CREATE TABLE IF NOT EXISTS "stocks"(
  "id" varchar UNIQUE NOT NULL, 
  "recruitment_id" varchar NULL, 
  "user_id" varchar NULL, 
  "created_at" timestamp with time zone NOT NULL, 
  "updated_at" timestamp with time zone NOT NULL, 
  PRIMARY KEY("id"),
  UNIQUE("user_id", "recruitment_id")
);
CREATE INDEX IF NOT EXISTS "stock_user_id" ON "stocks"("user_id");
CREATE INDEX IF NOT EXISTS "stock_recruitment_id" ON "stocks"("recruitment_id");
