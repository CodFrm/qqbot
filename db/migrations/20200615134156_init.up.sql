CREATE TABLE `scenes`
(
    `id`   int unsigned                                                 NOT NULL AUTO_INCREMENT,
    `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '场景名',
    `stat` tinyint(1)                                                   NOT NULL DEFAULT '1' COMMENT '0删除 1 显示 2 隐藏',
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

CREATE TABLE `scenes_tag`
(
    `id`        int                                                           NOT NULL AUTO_INCREMENT,
    `scenes_id` int                                                           NOT NULL,
    `key`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '原tag',
    `value`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '映射tag',
    PRIMARY KEY (`id`),
    UNIQUE KEY `scenes_key` (`scenes_id`, `key`),
    KEY `scenes_id` (`scenes_id`),
    KEY `key` (`key`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;