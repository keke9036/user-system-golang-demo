-- 建数据库
CREATE DATABASE IF NOT EXISTS user;
USE user;

-- 建表
CREATE TABLE if NOT EXISTS `user_tab`
(
    id          bigint unsigned auto_increment NOT NULL COMMENT '主键',
    user_id     bigint unsigned NOT NULL COMMENT 'user id',
    user_name   varchar(64)                    NOT NULL COMMENT '用户名',
    password    varchar(64)                    NOT NULL COMMENT '密码',
    nick_name   varchar(64)                    NOT NULL COMMENT '昵称',
    avatar_url  varchar(2048)                  NULL COMMENT '头像url',
    create_time DATETIME                       NOT NULL COMMENT '创建时间',
    modify_time datetime                       NOT NULL COMMENT '最近修改时间',
    PRIMARY KEY (id),
    UNIQUE KEY `uniq_username` (user_name),
    UNIQUE KEY `uniq_userid` (user_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
    COMMENT ='entry task用户表';
