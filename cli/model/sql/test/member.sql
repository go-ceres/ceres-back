CREATE TABLE `members` (
                                   `id` bigint(20) NOT NULL COMMENT '用户编号',
                                   `nickname` varchar(50) DEFAULT NULL COMMENT '用户昵称',
                                   `avatar` varchar(150) DEFAULT NULL COMMENT '头像',
                                   `password` varchar(150) NOT NULL COMMENT '用户密码',
                                   `email` varchar(50) DEFAULT NULL COMMENT '用户邮箱',
                                   `phone` varchar(11) DEFAULT NULL COMMENT '用户手机号',
                                   `describe` varchar(255) DEFAULT NULL COMMENT '用户简介',
                                   `status` tinyint(2) NOT NULL DEFAULT 1 COMMENT '用户状态',
                                   `createAt` int(11) DEFAULT NULL COMMENT '创建时间',
                                   `updateAt` int(11) DEFAULT NULL COMMENT '修改时间',
                                   `deleteAt` int(11) DEFAULT NULL COMMENT '删除时间',
                                   PRIMARY KEY (`id`),
                                   UNIQUE KEY `ceshi_uniqired` (`nickname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户表';
