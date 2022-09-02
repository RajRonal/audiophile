CREATE TABLE roles
(
    user_id     UUID NOT NULL REFERENCES users (user_id),
    user_role       TEXT              DEFAULT 'user',
    user_name    text not null REFERENCES users(user_name),
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);