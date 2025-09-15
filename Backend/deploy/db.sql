-- Database
DROP DATABASE IF EXISTS app;
CREATE DATABASE app;
USE app;

-- Tables
CREATE TABLE users (
  id         BIGINT UNSIGNED  AUTO_INCREMENT                PRIMARY KEY,
  username   VARCHAR(50)      DEFAULT "Trio",
  email      VARCHAR(255)     NOT NULL UNIQUE,
  password_hash VARCHAR(255)  NOT NULL,
  created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE user_devices (
  device_id BINARY(16) PRIMARY KEY,
  user_id   BIGINT UNSIGNED NOT NULL,
  last_seen TIMESTAMP NULL,
  revoked_at TIMESTAMP NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  KEY idx_user (user_id),
  UNIQUE KEY uk_user_device (user_id, device_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;