ALTER TABLE "recruitments"
  ADD CONSTRAINT "recruitments_competition_id"
    FOREIGN KEY("competition_id")
    REFERENCES "competitions"("id")
    ON DELETE RESTRICT,
  ADD CONSTRAINT "recruitments_prefecture_id"
    FOREIGN KEY("prefecture_id")
    REFERENCES "prefectures"("id")
    ON DELETE RESTRICT,
  ADD CONSTRAINT "recruitments_user_id"
    FOREIGN KEY("user_id")
    REFERENCES "users"("id")
    ON DELETE CASCADE;
