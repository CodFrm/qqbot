CREATE TABLE `scenes_tag`
(
    `id`        int                                                          NOT NULL AUTO_INCREMENT,
    `scenes_id` int                                                          NOT NULL,
    `key`       varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '原tag',
    `value`     varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '映射tag',
    PRIMARY KEY (`id`),
    KEY `scenes_id` (`scenes_id`),
    KEY `key` (`key`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;