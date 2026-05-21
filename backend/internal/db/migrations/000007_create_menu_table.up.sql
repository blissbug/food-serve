CREATE TABLE IF NOT EXISTS menus (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    date DATE NOT NULL,

    status ENUM('draft', 'active') DEFAULT 'draft',

    kitchen_id BIGINT NOT NULL,

    order_open DATETIME NULL,

    order_close DATETIME NULL,

    created_by BIGINT UNSIGNED,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    ON UPDATE CURRENT_TIMESTAMP,

    updated_by BIGINT UNSIGNED,

    deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_kitchen
    FOREIGN KEY (kitchen_id)
    REFERENCES kitchens(id)
    ON DELETE CASCADE,

    CONSTRAINT unique_kitchen_date
    UNIQUE (kitchen_id, date)
);