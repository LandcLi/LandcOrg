-- LandcOrg Database Schema

CREATE TABLE IF NOT EXISTS `departments` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(200) NOT NULL COMMENT 'Department name',
    `code` VARCHAR(100) NOT NULL COMMENT 'Department unique code',
    `parent_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '0 means root department',
    `path` VARCHAR(500) DEFAULT '/' COMMENT 'Full path like /1/2/3/',
    `level` INT DEFAULT 1 COMMENT 'Depth level in tree',
    `sort` INT DEFAULT 0 COMMENT 'Sort order within same level',
    `leader` VARCHAR(100) DEFAULT '' COMMENT 'Department leader name',
    `description` VARCHAR(500) DEFAULT '' COMMENT 'Department description',
    `status` TINYINT DEFAULT 1 COMMENT '1=active, 0=disabled',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_path` (`path`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `positions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(200) NOT NULL COMMENT 'Position/title name',
    `code` VARCHAR(100) NOT NULL COMMENT 'Position unique code',
    `level` INT DEFAULT 1 COMMENT 'Position level/rank',
    `description` VARCHAR(500) DEFAULT '' COMMENT 'Position description',
    `status` TINYINT DEFAULT 1 COMMENT '1=active, 0=disabled',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `org_members` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'Reference to user in LandcUserManager',
    `department_id` BIGINT UNSIGNED NOT NULL COMMENT 'Reference to department',
    `position_id` BIGINT UNSIGNED DEFAULT 0 COMMENT 'Reference to position',
    `is_leader` TINYINT(1) DEFAULT 0 COMMENT 'Is department leader',
    `joined_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'When user joined this department',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_dept` (`user_id`, `department_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_department_id` (`department_id`),
    KEY `idx_position_id` (`position_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;