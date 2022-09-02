ALTER TABLE product DROP COLUMN quantity ;
CREATE TABLE payment
(
   payment_id UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
   user_id   uuid REFERENCES   users(user_id),
   payment_type  TEXT NOT NULL,
   payment_date   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
 ALTER TABLE sessions
    RENAME COLUMN id TO user_id;
ALTER TABLE order_details
    ADD COLUMN  payment_id uuid REFERENCES payment(payment_id);
CREATE TABLE user_address
(
    address_id UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    user_id  uuid REFERENCES users(user_id),
    address_line_1  TEXT NOT NULL ,
    landmark     TEXT NOT NULL ,
    city       TEXT NOT NULL,
    postal_code   INTEGER
);
ALTER TABLE USERS
    RENAME COLUMN username TO user_name;
