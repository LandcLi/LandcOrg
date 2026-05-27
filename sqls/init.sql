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

-- Default organization structure
INSERT INTO
    `departments` (
        `name`,
        `code`,
        `parent_id`,
        `path`,
        `level`,
        `sort`,
        `description`
    )
VALUES (
        '总公司',
        'root',
        0,
        '/',
        1,
        0,
        'Root organization'
    ),
    (
        '技术部',
        'tech',
        1,
        '/1/',
        2,
        1,
        'Technology department'
    ),
    (
        '产品部',
        'product',
        1,
        '/1/',
        2,
        2,
        'Product department'
    ),
    (
        '运营部',
        'operation',
        1,
        '/1/',
        2,
        3,
        'Operation department'
    ),
    (
        '后端组',
        'tech_backend',
        2,
        '/1/2/',
        3,
        1,
        'Backend development team'
    ),
    (
        '前端组',
        'tech_frontend',
        2,
        '/1/2/',
        3,
        2,
        'Frontend development team'
    ),
    (
        '测试组',
        'tech_qa',
        2,
        '/1/2/',
        3,
        3,
        'QA team'
    );

-- Default positions
INSERT INTO
    `positions` (
        `name`,
        `code`,
        `level`,
        `description`
    )
VALUES (
        'CEO',
        'ceo',
        10,
        'Chief Executive Officer'
    ),
    (
        'CTO',
        'cto',
        9,
        'Chief Technology Officer'
    ),
    (
        '总监',
        'director',
        8,
        'Director'
    ),
    ('经理', 'manager', 7, 'Manager'),
    (
        '高级工程师',
        'senior_engineer',
        6,
        'Senior Engineer'
    ),
    (
        '工程师',
        'engineer',
        5,
        'Engineer'
    ),
    (
        '初级工程师',
        'junior_engineer',
        4,
        'Junior Engineer'
    ),
    ('实习生', 'intern', 1, 'Intern');