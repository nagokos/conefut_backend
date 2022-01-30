CREATE TABLE IF NOT EXISTS "recruitments"(
  "id" varchar UNIQUE NOT NULL, 
  "title" varchar(60) NOT NULL, 
  "type" varchar NOT NULL DEFAULT 'opponent', 
  "place" varchar NULL, 
  "start_at" timestamp with time zone NULL, 
  "content" varchar(10000) NOT NULL, 
  "location_url" varchar NULL, 
  "capacity" bigint NOT NULL, 
  "user_id" varchar NOT NULL, 
  "closing_at" timestamp with time zone NOT NULL, 
  "created_at" timestamp with time zone NOT NULL, 
  "updated_at" timestamp with time zone NOT NULL, 
  PRIMARY KEY("id")
);
