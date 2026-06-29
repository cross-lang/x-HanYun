SET NAMES utf8mb4;
SET
FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `tb_op_log`
(
    `id`         bigint       NOT NULL AUTO_INCREMENT COMMENT '唯一标识',
    `user_id`    varchar(255) NOT NULL DEFAULT '' COMMENT '操作用户ID',
    `op_time`    bigint       NOT NULL DEFAULT 0 COMMENT '操作时间',
    `op_dip`     varchar(50)  NOT NULL DEFAULT '' COMMENT '操作设备IP',
    `op_dname`   varchar(255) NOT NULL DEFAULT '' COMMENT '操作设备名称',
    `op_detail`  json         NOT NULL COMMENT '操作详情',
    `ctime`      bigint       NOT NULL DEFAULT 0 COMMENT '创建时间',
    `mtime`      bigint       NOT NULL DEFAULT 0 COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT '操作日志表';
