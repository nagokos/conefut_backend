CREATE TABLE IF NOT EXISTS "stocks"(
  "id" VARCHAR UNIQUE NOT NULL, 
  "recruitment_id" VARCHAR NULL, 
  "user_id" VARCHAR NULL, 
  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL, 
  "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL, 
  PRIMARY KEY("id"),
  UNIQUE("user_id", "recruitment_id")
);
CREATE INDEX IF NOT EXISTS "stocks_user_id" ON "stocks"("user_id");
CREATE INDEX IF NOT EXISTS "stocks_recruitment_id" ON "stocks"("recruitment_id");
