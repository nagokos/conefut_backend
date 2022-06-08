CREATE TABLE IF NOT EXISTS "messages"(
  "id" VARCHAR NOT NULL,
  "content" VARCHAR(1000) NOT NULL,
  "room_id" VARCHAR NOT NULL,
  "user_id" VARCHAR NULL,
  PRIMARY KEY("id"),
  FOREIGN KEY("room_id") 
    REFERENCES "rooms"("id")
    ON DELETE CASCADE,
  FOREIGN KEY("user_id")
    REFERENCES "users"("id")
    ON DELETE SET NULL
);
CREATE INDEX ON "messages"("room_id");
CREATE INDEX ON "messages"("user_id");