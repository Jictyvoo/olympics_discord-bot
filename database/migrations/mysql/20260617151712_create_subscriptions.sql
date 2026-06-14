-- Create "subscriptions" table
CREATE TABLE `subscriptions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `guild_id` varchar(32) NOT NULL,
  `user_id` varchar(32) NOT NULL,
  `kind` varchar(32) NOT NULL,
  `value` varchar(64) NULL,
  PRIMARY KEY (`id`),
  INDEX `subscription_guild_id` (`guild_id`),
  UNIQUE INDEX `subscription_guild_id_user_id_kind_value` (`guild_id`, `user_id`, `kind`, `value`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
