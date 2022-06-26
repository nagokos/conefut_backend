CREATE TABLE IF NOT EXISTS "prefectures"(
  "id" BIGSERIAL UNIQUE, 
  "name" VARCHAR NOT NULL, 
  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL, 
  "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL, 
  PRIMARY KEY("id")
);
