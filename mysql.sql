CREATE TABLE messages
(
    id           INT AUTO_INCREMENT PRIMARY KEY,
    message_id   VARCHAR(255) NOT NULL,
    number       VARCHAR(255) NOT NULL,
    context      LONGTEXT,
    is_send      BOOLEAN   DEFAULT 0,
    is_processed BOOLEAN   DEFAULT 0,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
