ALTER TABLE products DROP FOREIGN KEY products_ibfk_1;

ALTER TABLE products
    ADD CONSTRAINT products_ibfk_1 FOREIGN KEY (category_id)
        REFERENCES categories (id)
        ON UPDATE RESTRICT
        ON DELETE CASCADE;

CREATE INDEX name_idx ON categories (name)