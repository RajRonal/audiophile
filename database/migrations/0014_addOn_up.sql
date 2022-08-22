CREATE TABLE cart_sessions
(
    session_id   UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    user_id       UUID ,
    FOREIGN KEY (user_id)
        references users(user_id),
    archived_at timestamp  with time zone DEFAUlT NULL
);
ALTER TABLE cart_item DROP COLUMN session_id;
ALTER TABLE cart_item ADD COLUMN  session_id uuid REFERENCES  cart_sessions(session_id);