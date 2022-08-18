ALTER TABLE user_address
    ADD COLUMN archived_at timestamp  with time zone DEFAUlT NULL;
create unique index address_index ON  user_address (address_line_1,city,postal_code) where (archived_at is null);