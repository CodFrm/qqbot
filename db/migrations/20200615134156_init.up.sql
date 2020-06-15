CREATE TABLE `scenes`
(
    `id`   int unsigned                                                 NOT NULL AUTO_INCREMENT,
    `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '场景名',
    `stat` tinyint(1)                                                   NOT NULL DEFAULT '1' COMMENT '0删除 1 显示 2 隐藏',
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;
