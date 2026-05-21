CREATE TABLE IF NOT EXISTS food_item_images (
      id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

      food_item_id BIGINT UNSIGNED NOT NULL,

      image_url TEXT NOT NULL,

      original_destination TEXT NOT NULL,

      display_order INT DEFAULT 0,

      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

      FOREIGN KEY (food_item_id)
          REFERENCES food_items(id)
          ON DELETE CASCADE
);