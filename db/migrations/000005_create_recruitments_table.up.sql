CREATE TABLE IF NOT EXISTS "recruitments"(
  "id" varchar UNIQUE NOT NULL,
  "title" varchar(60) NOT NULL,
  "type" varchar NOT NULL DEFAULT 'unnecessary',
  "level" varchar NOT NULL DEFAULT 'unnecessary',
  "place" varchar NULL,
  "start_at" timestamp with time zone NULL,
  "detail" varchar(10000) NULL,
  "location_url" varchar NULL,
  "capacity" bigint NULL,
  "closing_at" timestamp with time zone NULL,
  "competition_id" varchar NULL,
  "prefecture_id" varchar NULL,
  "user_id" varchar NOT NULL,
  "created_at" timestamp with time zone NOT NULL,
  "updated_at" timestamp with time zone NOT NULL,
  PRIMARY KEY("id")
);
