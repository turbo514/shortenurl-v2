-- `default`.click_events_write definition
CREATE TABLE default.click_events_write
(
    `id` UUID,
    `link_id` UUID,
    `click_time` DateTime,
    `click_ip` IPv4,
    `user_agent` String,
    `referrer` String
)
    ENGINE = ReplacingMergeTree
PARTITION BY toYYYYMM(click_time)
ORDER BY id
SETTINGS index_granularity = 8192;