ALTER TABLE transportation_ordering_officers
    ADD COLUMN user_id uuid;
ALTER TABLE transportation_ordering_officers
    ADD CONSTRAINT transportation_ordering_officers_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id);
ALTER TABLE transportation_ordering_officers
    ADD CONSTRAINT transportation_ordering_officers_user_id_ukey UNIQUE (user_id);