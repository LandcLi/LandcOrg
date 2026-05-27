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

-- Organization invitations
CREATE TABLE IF NOT EXISTS `org_invitations` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `code` VARCHAR(64) NOT NULL COMMENT 'Invitation code',
    `department_id` BIGINT UNSIGNED NOT NULL COMMENT 'Target department',
    `position_id` BIGINT UNSIGNED DEFAULT 0 COMMENT 'Target position, 0=not specified',
    `inviter_id` BIGINT UNSIGNED NOT NULL COMMENT 'Who created this invitation',
    `invitee_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '0=open invitation, >0=specific user',
    `status` TINYINT DEFAULT 0 COMMENT '0=pending, 1=accepted, 2=expired, 3=revoked',
    `max_uses` INT DEFAULT 1 COMMENT 'Max number of uses, -1=unlimited',
    `used_count` INT DEFAULT 0 COMMENT 'Current usage count',
    `expires_at` DATETIME NOT NULL COMMENT 'Expiration time',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_inviter_id` (`inviter_id`),
    KEY `idx_invitee_id` (`invitee_id`),
    KEY `idx_department_id` (`department_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Department-Role binding (members of this department auto-inherit these roles)
CREATE TABLE IF NOT EXISTS `department_roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `department_id` BIGINT UNSIGNED NOT NULL,
    `role_id` BIGINT UNSIGNED NOT NULL COMMENT 'Reference to role in LandcPermission',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_dept_role` (`department_id`, `role_id`),
    KEY `idx_department_id` (`department_id`),
    KEY `idx_role_id` (`role_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- Position-Role binding (users with this position auto-inherit these roles)
CREATE TABLE IF NOT EXISTS `position_roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `position_id` BIGINT UNSIGNED NOT NULL,
    `role_id` BIGINT UNSIGNED NOT NULL COMMENT 'Reference to role in LandcPermission',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_pos_role` (`position_id`, `role_id`),
    KEY `idx_position_id` (`position_id`),
    KEY `idx_role_id` (`role_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;