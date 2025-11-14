CREATE TABLE `links` (
                         `id` binary(16) NOT NULL COMMENT 'uuidv7',
                         `tenant_id` binary(16) NOT NULL COMMENT '租户id',
                         `user_id` binary(16) NOT NULL COMMENT '用户id,审计用',
                         `short_code` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '短链接',
                         `original_url` varchar(255) NOT NULL COMMENT '原链接',
                         `created_at` timestamp NOT NULL COMMENT '创建时间',
                         `expires_at` timestamp NULL COMMENT '过期时间',
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `links_tenant_id_short_code_UNIQUE` (`tenant_id`,`short_code`) USING BTREE,
                         KEY `links_tenant_id_original_url_IDX` (`tenant_id`,`original_url`) USING BTREE,
                         KEY `links_short_code_IDX` (`short_code`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;