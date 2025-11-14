CREATE TABLE `tenants` (
                           `id` binary(16) NOT NULL COMMENT 'uuidv7',
                           `name` varchar(100) NOT NULL COMMENT '租客名',
                           `api_key` varchar(64) NOT NULL COMMENT '密钥,用于 API 调用',
                           `created_at` timestamp NOT NULL COMMENT '创建时间',
                           `updated_at` timestamp NOT NULL COMMENT '修改时间',
                           `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `tenants_name_deleted_at_unqiue` (`name`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

