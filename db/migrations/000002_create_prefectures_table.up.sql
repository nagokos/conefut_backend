CREATE TABLE IF NOT EXISTS "prefectures"(
  "id" varchar UNIQUE NOT NULL, 
  "name" varchar NOT NULL, 
  "created_at" timestamp with time zone NOT NULL, 
  "updated_at" timestamp with time zone NOT NULL, 
  PRIMARY KEY("id")
);
