CREATE TABLE IF NOT EXISTS "competitions"(
  "id" varchar UNIQUE NOT NULL, 
  "name" varchar UNIQUE NOT NULL, 
  "created_at" timestamp with time zone NOT NULL, 
  "updated_at" timestamp with time zone NOT NULL, 
  PRIMARY KEY("id")
);
