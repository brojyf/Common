-- Basic Settings
SET NAMES utf8mb4;
SET time_zone = '+00:00';

-- Database
DROP DATABASE IF EXISTS app;
CREATE DATABASE app CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE app;

-- Tables
CREATE TABLE users (
    -- Basics
    id              BIGINT UNSIGNED AUTO_INCREMENT,
    email           VARCHAR(255)  NOT NULL,
    password_hash   VARCHAR(255)  NOT NULL,
    username        VARCHAR(50)   NOT NULL DEFAULT "J",
     profile_photo   VARCHAR(1024) NOT NULL DEFAULT "http://localhost:8080/pfp",
    -- Record
    token_version   INT UNSIGNED NOT NULL DEFAULT 1,
    is_deleted      TINYINT(1) NOT NULL DEFAULT 0,
    -- Auto
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    -- Constraints
    PRIMARY KEY (id),
    CONSTRAINT uk_users_email UNIQUE (email)
) ENGINE=InnoDB;

CREATE TABLE user_devices (
    -- Basics
    device_id    BINARY(16) PRIMARY KEY,
    user_id      BIGINT UNSIGNED NOT NULL,
    -- Record
    push_token   VARCHAR(512) NULL,
    last_seen_at TIMESTAMP NULL,
    revoked_at   TIMESTAMP NULL,
    -- Auto
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Constraints
    UNIQUE KEY uk_user_device (user_id, device_id),
    CONSTRAINT fk_ud_user FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE INDEX idx_ud_user ON user_devices(user_id);
CREATE INDEX idx_ud_last_seen ON user_devices(last_seen_at);

CREATE TABLE sessions (
    -- Basics
    session_id     BIGINT UNSIGNED AUTO_INCREMENT,
    user_id        BIGINT UNSIGNED NOT NULL,
    device_id      BINARY(16) NOT NULL,
    rtk_hash       VARBINARY(32) NOT NULL,
    token_version  INT UNSIGNED NOT NULL,
    expires_at     TIMESTAMP NOT NULL,
    revoked_at     TIMESTAMP NULL,
    -- Auto
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    -- Constraints
    PRIMARY KEY (session_id),
    UNIQUE KEY uk_sessions_user_device (user_id, device_id),
    CONSTRAINT fk_sess_user_device FOREIGN KEY (user_id, device_id)
    REFERENCES user_devices (user_id, device_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB;
