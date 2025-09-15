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
  token_version INT UNSIGNED NOT NULL DEFAULT 1,
  created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE user_devices (
  device_id BINARY(16) PRIMARY KEY,
  user_id   BIGINT UNSIGNED NOT NULL,
  revoked_at TIMESTAMP NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  KEY idx_user (user_id),
  UNIQUE KEY uk_user_device (user_id, device_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sessions (
    user_id           BIGINT UNSIGNED NOT NULL,
    device_id         BINARY(16)      NOT NULL,
    rtk_hash          VARBINARY(32)   NOT NULL,
    refresh_expires_at DATETIME       NOT NULL,
    revoked_at        DATETIME        NULL,
    created_at        TIMESTAMP       DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, device_id),
    KEY idx_user (user_id),
    CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_sessions_device FOREIGN KEY (device_id) REFERENCES user_devices(device_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;