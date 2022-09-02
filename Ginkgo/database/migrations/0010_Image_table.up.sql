CREATE TABLE image_details
(
    image_id  uuid,
    product_id uuid REFERENCES product(product_id)
);
