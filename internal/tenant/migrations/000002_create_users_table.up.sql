-- dev.users definition

CREATE TABLE `users` (
                         `id` binary(16) NOT NULL COMMENT 'uuidv7',
                         `tenant_id` binary(16) NOT NULL COMMENT '租户id',
                         `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名,登录标识符',
                         `password` char(60) NOT NULL COMMENT 'bcrypt哈希密码',
                         `created_at` timestamp NOT NULL COMMENT '创建时间',
                         `updated_at` timestamp NOT NULL COMMENT '更新时间',
                         `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `users_unique` (`name`),
                         KEY `users_name_password_IDX` (`name`,`password`) USING BTREE,
                         KEY `users_tenant_id_IDX` (`tenant_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;