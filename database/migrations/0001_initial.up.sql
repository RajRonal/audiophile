CREATE TABLE   users
(
    user_id            UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    first_name          TEXT NOT NULL,
    last_name      TEXT NOT NULL,
    email         TEXT NOT NULL UNIQUE ,
    contact_number  TEXT NOT NULL UNIQUE ,
    username      TEXT NOT NULL UNIQUE,
    password       TEXT NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at   TIMESTAMP WITH TIME ZONE
);
CREATE TABLE sessions
(
    session_id   UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    expired_at     TIMESTAMP with time zone,
    id       UUID ,
    FOREIGN KEY (id)
        references users(user_id),
    archived_at timestamp  with time zone DEFAUlT NULL
);