CREATE TABLE IF NOT EXISTS "entries"(
  "id" VARCHAR NOT NULL,
  "room_id" VARCHAR NOT NULL,
  "user_id" VARCHAR NULL,
  PRIMARY KEY("id"),
  FOREIGN KEY("room_id") 
    REFERENCES "rooms"("id")
    ON DELETE CASCADE,
  FOREIGN KEY("user_id")
    REFERENCES "users"("id")
    ON DELETE SET NULL,
  UNIQUE("room_id", "user_id")
);
CREATE INDEX ON "entries"("room_id");
CREATE INDEX ON "entries"("user_id");
