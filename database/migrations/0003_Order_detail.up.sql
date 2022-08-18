CREATE TABLE  cart_item
(
     cart_id UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
     session_id  uuid REFERENCES sessions(SESSION_ID),
     product_id  uuid REFERENCES product(product_id),
     coupon_id   uuid REFERENCES discount(coupon_id),
     quantity    INTEGER ,
     created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
     archived_at TIMESTAMP WITH TIME ZONE DEFAUlT NULL
);
CREATE TABLE order_details
(
    order_id  UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    user_id uuid REFERENCES users(user_id),
    total    DOUBLE PRECISION,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);