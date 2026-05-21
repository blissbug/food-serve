CREATE TABLE IF NOT EXISTS menu_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    menu_id BIGINT UNSIGNED NOT NULL,
    food_item_id BIGINT UNSIGNED NOT NULL,

    price DECIMAL(10,2) NOT NULL,

    is_available BOOLEAN DEFAULT TRUE,

    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_menu
    FOREIGN KEY (menu_id)
    REFERENCES menus(id)
    ON DELETE CASCADE,

    CONSTRAINT fk_food_item
    FOREIGN KEY (food_item_id)
    REFERENCES food_items(id)
    ON DELETE CASCADE,

    CONSTRAINT unique_menu_food
    UNIQUE(menu_id, food_item_id)
);