CREATE TABLE IF NOT EXISTS "stocks"(
  "id" varchar UNIQUE NOT NULL, 
  "recruitment_id" varchar NULL, 
  "user_id" varchar NULL, 
  "created_at" timestamp with time zone NOT NULL, 
  "updated_at" timestamp with time zone NOT NULL, 
  PRIMARY KEY("id")
);
CREATE UNIQUE INDEX IF NOT EXISTS "stock_user_id_recruitment_id" ON "stocks"("user_id", "recruitment_id");
