CREATE TABLE  product_categories
(
    category_id  UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    category_name  TEXT NOT NULL ,
    category_description TEXT NOT NULL ,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at timestamp  with time zone DEFAUlT NULL
);
CREATE TABLE product_inventory
(
    inventory_id UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    quantity     integer,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at timestamp  with time zone DEFAUlT NULL

);
CREATE TABLE product
(
    product_id UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    product_name  TEXT NOT NULL ,
    product_description TEXT NOT NULL ,
    category_id     uuid REFERENCES product_categories(category_id),
    inventory_id     uuid REFERENCES product_inventory(inventory_id),
    regular_price    DOUBLE PRECISION,
    discounted_price  DOUBLE PRECISION  default 0.0,
    quantity          INTEGER,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE DEFAUlT NULL
);
CREATE TABLE discount
(
    coupon_id    UUID PRIMARY KEY  DEFAULT gen_random_uuid(),
    coupon_name    TEXT NOT NULL ,
    coupon_description  TEXT,
    discount_percentage INTEGER,
    discount_status      BOOLEAN DEFAULT FALSE,
    start_date    TIMESTAMP WITH TIME ZONE,
    end_date    TIMESTAMP WITH TIME ZONE
);
