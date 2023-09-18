ALTER TABLE categories MODIFY name varchar(190);

CREATE INDEX name_idx ON categories (name)