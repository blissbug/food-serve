CREATE TABLE IF NOT EXISTS kitchen_members (
   id BIGINT AUTO_INCREMENT PRIMARY KEY,

   kitchen_id BIGINT NOT NULL,
   user_id BIGINT NOT NULL,
   role VARCHAR(255) NOT NULL DEFAULT 'user',

    UNIQUE (kitchen_id, user_id),

    FOREIGN KEY (kitchen_id)
    REFERENCES kitchens(id)
    ON DELETE CASCADE,

    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);