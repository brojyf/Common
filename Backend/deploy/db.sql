-- Basic Settings
SET NAMES utf8mb4;
SET time_zone = '+00:00';

-- Database
DROP DATABASE IF EXISTS app;
CREATE DATABASE app CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE app;

-- Tables
CREATE TABLE users (
                       id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                       email           VARCHAR(255)  NOT NULL,
                       password_hash   VARCHAR(255)  NOT NULL,
                       username        VARCHAR(32)   NULL,
                       profile_photo   VARCHAR(1024) NULL,

                       token_version   INT UNSIGNED NOT NULL DEFAULT 1,
                       is_deleted      TINYINT(1) NOT NULL DEFAULT 0,

                       created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

                       CONSTRAINT uk_users_email UNIQUE (email)
) ENGINE=InnoDB;

CREATE TABLE user_devices (
                              device_id    BINARY(16) PRIMARY KEY,
                              user_id      BIGINT UNSIGNED NOT NULL,

                              push_token   VARCHAR(255) NULL,
                              last_seen_at TIMESTAMP NULL,
                              revoked_at   TIMESTAMP NULL,

                              created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                              CONSTRAINT fk_ud_user FOREIGN KEY (user_id) REFERENCES users(id)
                                  ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE INDEX idx_ud_user ON user_devices(user_id);
CREATE INDEX idx_ud_last_seen ON user_devices(last_seen_at);

CREATE TABLE sessions (
                          session_id     BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                          user_id        BIGINT UNSIGNED NOT NULL,
                          device_id      BINARY(16) NOT NULL,              -- 关键：NOT NULL
                          rtk_hash       VARBINARY(32) NOT NULL,           -- SHA-256 二进制32字节

                          token_version  INT UNSIGNED NOT NULL,
                          expires_at     TIMESTAMP NOT NULL,
                          revoked_at     TIMESTAMP NULL,
                          created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

                          UNIQUE KEY uk_sessions_user_device (user_id, device_id),   -- 用逗号，不是分号
                          CONSTRAINT fk_sess_user   FOREIGN KEY (user_id)   REFERENCES users(id) ON DELETE CASCADE,
                          CONSTRAINT fk_sess_device FOREIGN KEY (device_id) REFERENCES user_devices(device_id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE INDEX idx_sessions_user         ON sessions(user_id);
CREATE INDEX idx_sessions_expires      ON sessions(expires_at);
CREATE INDEX idx_sessions_revoked      ON sessions(revoked_at);
CREATE INDEX idx_sessions_user_valid   ON sessions(user_id, revoked_at, expires_at);

CREATE TABLE friend_requests (
                                 id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                 requester_id  BIGINT UNSIGNED NOT NULL,
                                 addressee_id  BIGINT UNSIGNED NOT NULL,
                                 status        ENUM('pending','accepted','rejected','cancelled','expired') NOT NULL DEFAULT 'pending',
                                 message       VARCHAR(200) NULL,

                                 created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                 decided_at    TIMESTAMP NULL,

                                 CONSTRAINT fk_fr_req FOREIGN KEY (requester_id) REFERENCES users(id) ON DELETE CASCADE,
                                 CONSTRAINT fk_fr_add FOREIGN KEY (addressee_id) REFERENCES users(id) ON DELETE CASCADE,
                                 CONSTRAINT uk_fr_unique UNIQUE (requester_id, addressee_id),
                                 CONSTRAINT chk_fr_self CHECK (requester_id <> addressee_id)
) ENGINE=InnoDB;

CREATE INDEX idx_fr_addressee_pending ON friend_requests(addressee_id, status);

CREATE TABLE friendships (
                             user_id     BIGINT UNSIGNED NOT NULL,
                             friend_id   BIGINT UNSIGNED NOT NULL,
                             created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                             PRIMARY KEY (user_id, friend_id),
                             CONSTRAINT chk_friend_order CHECK (user_id < friend_id),
                             CONSTRAINT fk_fs_user   FOREIGN KEY (user_id)   REFERENCES users(id) ON DELETE CASCADE,
                             CONSTRAINT fk_fs_friend FOREIGN KEY (friend_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE `groups` (
                        id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                        owner_id      BIGINT UNSIGNED NOT NULL,
                        name          VARCHAR(100) NOT NULL,
                        description   VARCHAR(500) NULL,
                        avatar        VARCHAR(1024) NULL,

                        created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

                        CONSTRAINT fk_grp_owner FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE INDEX idx_groups_owner ON `groups`(owner_id);

CREATE TABLE group_members (
                               group_id     BIGINT UNSIGNED NOT NULL,
                               user_id      BIGINT UNSIGNED NOT NULL,
                               joined_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

                               PRIMARY KEY (group_id, user_id),
                               CONSTRAINT fk_gm_group FOREIGN KEY (group_id) REFERENCES `groups`(id) ON DELETE CASCADE,
                               CONSTRAINT fk_gm_user  FOREIGN KEY (user_id)  REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE INDEX idx_gm_user ON group_members(user_id);

CREATE TABLE group_invites (
                               id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                               group_id     BIGINT UNSIGNED NOT NULL,
                               inviter_id   BIGINT UNSIGNED NOT NULL,
                               invitee_id   BIGINT UNSIGNED NOT NULL,
                               status       ENUM('pending','accepted','rejected','expired','revoked') NOT NULL DEFAULT 'pending',

                               created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               decided_at   TIMESTAMP NULL,

                               CONSTRAINT fk_gi_group   FOREIGN KEY (group_id)   REFERENCES `groups`(id) ON DELETE CASCADE,
                               CONSTRAINT fk_gi_inviter FOREIGN KEY (inviter_id) REFERENCES users(id)  ON DELETE CASCADE,
                               CONSTRAINT fk_gi_invitee FOREIGN KEY (invitee_id) REFERENCES users(id)  ON DELETE CASCADE,
                               CONSTRAINT uk_gi_unique UNIQUE (group_id, invitee_id)
) ENGINE=InnoDB;
