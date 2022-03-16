ALTER TABLE "recruitment_tags" 
  ADD CONSTRAINT "recruitment_tags_recruitments_recruitment_tags" 
    FOREIGN KEY("recruitment_id") 
    REFERENCES "recruitments"("id") 
    ON DELETE CASCADE, 
  ADD CONSTRAINT "recruitment_tags_tags_recruitment_tags" 
    FOREIGN KEY("tag_id") 
    REFERENCES "tags"("id") 
    ON DELETE CASCADE;
