ALTER TABLE "recruitments"
  ADD CONSTRAINT "recruitments_competitions_recruitments"
    FOREIGN KEY("competition_id")
    REFERENCES "competitions"("id")
    ON DELETE RESTRICT,
  ADD CONSTRAINT "recruitments_prefectures_recruitments"
    FOREIGN KEY("prefecture_id")
    REFERENCES "prefectures"("id")
    ON DELETE RESTRICT,
  ADD CONSTRAINT "recruitments_users_recruitments"
    FOREIGN KEY("user_id")
    REFERENCES "users"("id")
    ON DELETE CASCADE;
