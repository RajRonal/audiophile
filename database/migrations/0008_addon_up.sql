ALTER TABLE product_inventory ADD COLUMN category_id uuid REFERENCES product_categories(category_id);