CREATE TABLE products (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    description VARCHAR(255) NOT NULL,
    image VARCHAR(255) NOT NULL,
    created_at datetime NOT NULL,
    category_id BIGINT NOT NULL,
	FOREIGN KEY (category_id) REFERENCES categories(id),
	KEY `title_idx` (`title`)
);