CREATE TABLE IF NOT EXISTS "recruitments"(
  "id" VARCHAR UNIQUE NOT NULL,
  "title" VARCHAR(60) NOT NULL,
  "type" VARCHAR NOT NULL DEFAULT 'unnecessary',
  "level" VARCHAR NOT NULL DEFAULT 'unnecessary',
  "place" VARCHAR NULL,
  "start_at" TIMESTAMP WITH TIME ZONE NULL,
  "content" VARCHAR(10000) NULL,
  "location_url" VARCHAR NULL,
  "capacity" BIGINT NULL,
  "closing_at" TIMESTAMP WITH TIME ZONE NULL,
  "competition_id" VARCHAR NULL,
  "prefecture_id" VARCHAR NULL,
  "user_id" VARCHAR NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
  "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL,
  PRIMARY KEY("id")
);
