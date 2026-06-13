-- Create "standings" table
CREATE TABLE `standings` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `stage_id` binary(16) NOT NULL,
  `participant_id` binary(16) NOT NULL,
  `rank` bigint NOT NULL DEFAULT 0,
  `points` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `standing_stage_id_participant_id` (`stage_id`, `participant_id`),
  INDEX `standing_stage_id_rank` (`stage_id`, `rank`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "competitions" table
CREATE TABLE `competitions` (
  `id` binary(16) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `provider_id` varchar(255) NOT NULL,
  `external_key` varchar(255) NOT NULL,
  `code` varchar(255) NULL,
  `name` varchar(255) NOT NULL,
  `discipline` varchar(255) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `competition_provider_id_external_key` (`provider_id`, `external_key`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "countries" table
CREATE TABLE `countries` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `iso2` varchar(2) NOT NULL,
  `iso3` varchar(3) NOT NULL,
  `ioc_code` varchar(3) NULL,
  `name` varchar(255) NOT NULL,
  `code_num` bigint NULL,
  `population` bigint NULL,
  `area_km2` double NULL,
  `gdp_usd` double NULL,
  PRIMARY KEY (`id`),
  INDEX `country_ioc_code` (`ioc_code`),
  INDEX `country_name` (`name`),
  UNIQUE INDEX `iso2` (`iso2`),
  UNIQUE INDEX `iso3` (`iso3`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "venues" table
CREATE TABLE `venues` (
  `id` binary(16) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `provider_id` varchar(255) NOT NULL,
  `external_key` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `city` varchar(255) NULL,
  `country_iso` varchar(255) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `venue_provider_id_external_key` (`provider_id`, `external_key`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "participants" table
CREATE TABLE `participants` (
  `id` binary(16) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `provider_id` varchar(255) NOT NULL,
  `external_key` varchar(255) NOT NULL,
  `kind` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `code` varchar(255) NULL,
  `country_iso` varchar(255) NULL,
  `gender` varchar(255) NULL,
  PRIMARY KEY (`id`),
  INDEX `participant_country_iso` (`country_iso`),
  UNIQUE INDEX `participant_provider_id_external_key` (`provider_id`, `external_key`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "sync_states" table
CREATE TABLE `sync_states` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `provider_id` varchar(255) NOT NULL,
  `scope` varchar(255) NOT NULL,
  `cursor` varchar(255) NULL,
  `last_synced_at` timestamp NULL,
  `last_error` varchar(255) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `syncstate_provider_id_scope` (`provider_id`, `scope`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "seasons" table
CREATE TABLE `seasons` (
  `id` binary(16) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `provider_id` varchar(255) NOT NULL,
  `external_key` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `starts_on` timestamp NOT NULL,
  `ends_on` timestamp NOT NULL,
  `competition_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `season_competition_id` (`competition_id`),
  UNIQUE INDEX `season_provider_id_external_key` (`provider_id`, `external_key`),
  CONSTRAINT `seasons_competitions_seasons` FOREIGN KEY (`competition_id`) REFERENCES `competitions` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "stages" table
CREATE TABLE `stages` (
  `id` binary(16) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `provider_id` varchar(255) NOT NULL,
  `external_key` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `ord` bigint NOT NULL DEFAULT 0,
  `season_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `stage_provider_id_external_key` (`provider_id`, `external_key`),
  INDEX `stage_season_id` (`season_id`),
  CONSTRAINT `stages_seasons_stages` FOREIGN KEY (`season_id`) REFERENCES `seasons` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "fixtures" table
CREATE TABLE `fixtures` (
  `id` binary(16) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `provider_id` varchar(255) NOT NULL,
  `external_key` varchar(255) NOT NULL,
  `group_id` binary(16) NULL,
  `venue_id` binary(16) NULL,
  `name` varchar(255) NOT NULL,
  `starts_at` timestamp NOT NULL,
  `ends_at` timestamp NOT NULL,
  `status` varchar(255) NOT NULL DEFAULT "scheduled",
  `checksum` varchar(64) NULL,
  `stage_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `fixture_provider_id_external_key` (`provider_id`, `external_key`),
  INDEX `fixture_stage_id` (`stage_id`),
  INDEX `fixture_starts_at` (`starts_at`),
  INDEX `fixture_status` (`status`),
  CONSTRAINT `fixtures_stages_fixtures` FOREIGN KEY (`stage_id`) REFERENCES `stages` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "alerts" table
CREATE TABLE `alerts` (
  `id` binary(16) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `kind` varchar(255) NOT NULL,
  `fixture_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `alert_fixture_id` (`fixture_id`),
  CONSTRAINT `alerts_fixtures_alerts` FOREIGN KEY (`fixture_id`) REFERENCES `fixtures` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "discord_events" table
CREATE TABLE `discord_events` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `guild_id` varchar(255) NOT NULL,
  `discord_event_id` varchar(255) NOT NULL,
  `status` varchar(255) NOT NULL DEFAULT "scheduled",
  `last_checksum` varchar(64) NULL,
  `fixture_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `discord_event_id` (`discord_event_id`),
  UNIQUE INDEX `discordevent_discord_event_id` (`discord_event_id`),
  UNIQUE INDEX `discordevent_fixture_id_guild_id` (`fixture_id`, `guild_id`),
  CONSTRAINT `discord_events_fixtures_discord_events` FOREIGN KEY (`fixture_id`) REFERENCES `fixtures` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "fixture_participants" table
CREATE TABLE `fixture_participants` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `participant_id` binary(16) NOT NULL,
  `role` varchar(255) NULL,
  `fixture_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `fixtureparticipant_fixture_id_participant_id` (`fixture_id`, `participant_id`),
  INDEX `fixtureparticipant_participant_id` (`participant_id`),
  CONSTRAINT `fixture_participants_fixtures_fixture_participants` FOREIGN KEY (`fixture_id`) REFERENCES `fixtures` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "groups" table
CREATE TABLE `groups` (
  `id` binary(16) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `provider_id` varchar(255) NOT NULL,
  `external_key` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `stage_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `group_provider_id_external_key` (`provider_id`, `external_key`),
  INDEX `group_stage_id` (`stage_id`),
  CONSTRAINT `groups_stages_groups` FOREIGN KEY (`stage_id`) REFERENCES `stages` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "notifications" table
CREATE TABLE `notifications` (
  `id` binary(16) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `channel_id` varchar(255) NULL,
  `message_id` varchar(255) NULL,
  `status` varchar(255) NOT NULL DEFAULT "pending",
  `checksum` varchar(64) NULL,
  `sent_at` timestamp NULL,
  `alert_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `notification_alert_id` (`alert_id`),
  INDEX `notification_checksum` (`checksum`),
  CONSTRAINT `notifications_alerts_notifications` FOREIGN KEY (`alert_id`) REFERENCES `alerts` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "results" table
CREATE TABLE `results` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT (now()),
  `updated_at` timestamp NOT NULL DEFAULT (now()),
  `participant_id` binary(16) NOT NULL,
  `position` bigint NULL,
  `score` varchar(255) NULL,
  `raw_mark` varchar(255) NULL,
  `outcome` varchar(255) NULL,
  `fixture_id` binary(16) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `result_fixture_id_participant_id` (`fixture_id`, `participant_id`),
  INDEX `result_participant_id` (`participant_id`),
  CONSTRAINT `results_fixtures_results` FOREIGN KEY (`fixture_id`) REFERENCES `fixtures` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
