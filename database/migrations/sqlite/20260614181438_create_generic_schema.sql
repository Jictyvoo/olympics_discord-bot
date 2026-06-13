-- Create "alerts" table
CREATE TABLE alerts (
  id blob NOT NULL,
  created_at datetime NOT NULL DEFAULT (DATETIME('now')),
  updated_at datetime NOT NULL DEFAULT (DATETIME('now')),
  kind text NOT NULL,
  fixture_id blob NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT alerts_fixtures_alerts FOREIGN KEY (fixture_id) REFERENCES fixtures (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "alert_fixture_id" to table: "alerts"
CREATE INDEX alert_fixture_id ON alerts (fixture_id);
-- Create "competitions" table
CREATE TABLE competitions (
  id blob NOT NULL,
  created_at datetime NOT NULL DEFAULT (DATETIME('now')),
  updated_at datetime NOT NULL DEFAULT (DATETIME('now')),
  provider_id text NOT NULL,
  external_key text NOT NULL,
  code text NULL,
  name text NOT NULL,
  discipline text NULL,
  PRIMARY KEY (id)
);
-- Create index "competition_provider_id_external_key" to table: "competitions"
CREATE UNIQUE INDEX competition_provider_id_external_key ON competitions (provider_id, external_key);
-- Create "countries" table
CREATE TABLE countries (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  iso2 text NOT NULL,
  iso3 text NOT NULL,
  ioc_code text NULL,
  name text NOT NULL,
  code_num integer NULL,
  population integer NULL,
  area_km2 real NULL,
  gdp_usd real NULL
);
-- Create index "countries_iso2_key" to table: "countries"
CREATE UNIQUE INDEX countries_iso2_key ON countries (iso2);
-- Create index "countries_iso3_key" to table: "countries"
CREATE UNIQUE INDEX countries_iso3_key ON countries (iso3);
-- Create index "country_ioc_code" to table: "countries"
CREATE INDEX country_ioc_code ON countries (ioc_code);
-- Create index "country_name" to table: "countries"
CREATE INDEX country_name ON countries (name);
-- Create "discord_events" table
CREATE TABLE discord_events (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  created_at datetime NOT NULL DEFAULT (DATETIME('now')),
  updated_at datetime NOT NULL DEFAULT (DATETIME('now')),
  guild_id text NOT NULL,
  discord_event_id text NOT NULL,
  status text NOT NULL DEFAULT 'scheduled',
  last_checksum text NULL,
  fixture_id blob NOT NULL,
  CONSTRAINT discord_events_fixtures_discord_events FOREIGN KEY (fixture_id) REFERENCES fixtures (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "discord_events_discord_event_id_key" to table: "discord_events"
CREATE UNIQUE INDEX discord_events_discord_event_id_key ON discord_events (discord_event_id);
-- Create index "discordevent_fixture_id_guild_id" to table: "discord_events"
CREATE UNIQUE INDEX discordevent_fixture_id_guild_id ON discord_events (fixture_id, guild_id);
-- Create index "discordevent_discord_event_id" to table: "discord_events"
CREATE UNIQUE INDEX discordevent_discord_event_id ON discord_events (discord_event_id);
-- Create "fixtures" table
CREATE TABLE fixtures (
  id blob NOT NULL,
  created_at datetime NOT NULL DEFAULT (DATETIME('now')),
  updated_at datetime NOT NULL DEFAULT (DATETIME('now')),
  provider_id text NOT NULL,
  external_key text NOT NULL,
  group_id blob NULL,
  venue_id blob NULL,
  name text NOT NULL,
  starts_at datetime NOT NULL,
  ends_at datetime NOT NULL,
  status text NOT NULL DEFAULT 'scheduled',
  checksum text NULL,
  stage_id blob NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fixtures_stages_fixtures FOREIGN KEY (stage_id) REFERENCES stages (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "fixture_provider_id_external_key" to table: "fixtures"
CREATE UNIQUE INDEX fixture_provider_id_external_key ON fixtures (provider_id, external_key);
-- Create index "fixture_stage_id" to table: "fixtures"
CREATE INDEX fixture_stage_id ON fixtures (stage_id);
-- Create index "fixture_starts_at" to table: "fixtures"
CREATE INDEX fixture_starts_at ON fixtures (starts_at);
-- Create index "fixture_status" to table: "fixtures"
CREATE INDEX fixture_status ON fixtures (status);
-- Create "fixture_participants" table
CREATE TABLE fixture_participants (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  participant_id blob NOT NULL,
  role text NULL,
  fixture_id blob NOT NULL,
  CONSTRAINT fixture_participants_fixtures_fixture_participants FOREIGN KEY (fixture_id) REFERENCES fixtures (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "fixtureparticipant_fixture_id_participant_id" to table: "fixture_participants"
CREATE UNIQUE INDEX fixtureparticipant_fixture_id_participant_id ON fixture_participants (fixture_id, participant_id);
-- Create index "fixtureparticipant_participant_id" to table: "fixture_participants"
CREATE INDEX fixtureparticipant_participant_id ON fixture_participants (participant_id);
-- Create "groups" table
CREATE TABLE groups (
  id blob NOT NULL,
  created_at datetime NOT NULL DEFAULT (DATETIME('now')),
  updated_at datetime NOT NULL DEFAULT (DATETIME('now')),
  provider_id text NOT NULL,
  external_key text NOT NULL,
  name text NOT NULL,
  stage_id blob NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT groups_stages_groups FOREIGN KEY (stage_id) REFERENCES stages (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "group_provider_id_external_key" to table: "groups"
CREATE UNIQUE INDEX group_provider_id_external_key ON groups (provider_id, external_key);
-- Create index "group_stage_id" to table: "groups"
CREATE INDEX group_stage_id ON groups (stage_id);
-- Create "notifications" table
CREATE TABLE notifications (
  id blob NOT NULL,
  created_at datetime NOT NULL DEFAULT (DATETIME('now')),
  updated_at datetime NOT NULL DEFAULT (DATETIME('now')),
  channel_id text NULL,
  message_id text NULL,
  status text NOT NULL DEFAULT 'pending',
  checksum text NULL,
  sent_at datetime NULL,
  alert_id blob NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT notifications_alerts_notifications FOREIGN KEY (alert_id) REFERENCES alerts (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "notification_alert_id" to table: "notifications"
CREATE INDEX notification_alert_id ON notifications (alert_id);
-- Create index "notification_checksum" to table: "notifications"
CREATE INDEX notification_checksum ON notifications (checksum);
-- Create "participants" table
CREATE TABLE participants (
  id blob NOT NULL,
  created_at datetime NOT NULL DEFAULT (DATETIME('now')),
  updated_at datetime NOT NULL DEFAULT (DATETIME('now')),
  provider_id text NOT NULL,
  external_key text NOT NULL,
  kind text NOT NULL,
  name text NOT NULL,
  code text NULL,
  country_iso text NULL,
  gender text NULL,
  PRIMARY KEY (id)
);
-- Create index "participant_provider_id_external_key" to table: "participants"
CREATE UNIQUE INDEX participant_provider_id_external_key ON participants (provider_id, external_key);
-- Create index "participant_country_iso" to table: "participants"
CREATE INDEX participant_country_iso ON participants (country_iso);
-- Create "results" table
CREATE TABLE results (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  created_at datetime NOT NULL DEFAULT (DATETIME('now')),
  updated_at datetime NOT NULL DEFAULT (DATETIME('now')),
  participant_id blob NOT NULL,
  position integer NULL,
  score text NULL,
  raw_mark text NULL,
  outcome text NULL,
  fixture_id blob NOT NULL,
  CONSTRAINT results_fixtures_results FOREIGN KEY (fixture_id) REFERENCES fixtures (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "result_fixture_id_participant_id" to table: "results"
CREATE UNIQUE INDEX result_fixture_id_participant_id ON results (fixture_id, participant_id);
-- Create index "result_participant_id" to table: "results"
CREATE INDEX result_participant_id ON results (participant_id);
-- Create "seasons" table
CREATE TABLE seasons (
  id blob NOT NULL,
  created_at datetime NOT NULL DEFAULT (DATETIME('now')),
  updated_at datetime NOT NULL DEFAULT (DATETIME('now')),
  provider_id text NOT NULL,
  external_key text NOT NULL,
  name text NOT NULL,
  starts_on datetime NOT NULL,
  ends_on datetime NOT NULL,
  competition_id blob NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT seasons_competitions_seasons FOREIGN KEY (competition_id) REFERENCES competitions (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "season_provider_id_external_key" to table: "seasons"
CREATE UNIQUE INDEX season_provider_id_external_key ON seasons (provider_id, external_key);
-- Create index "season_competition_id" to table: "seasons"
CREATE INDEX season_competition_id ON seasons (competition_id);
-- Create "stages" table
CREATE TABLE stages (
  id blob NOT NULL,
  created_at datetime NOT NULL DEFAULT (DATETIME('now')),
  updated_at datetime NOT NULL DEFAULT (DATETIME('now')),
  provider_id text NOT NULL,
  external_key text NOT NULL,
  name text NOT NULL,
  ord integer NOT NULL DEFAULT 0,
  season_id blob NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT stages_seasons_stages FOREIGN KEY (season_id) REFERENCES seasons (id) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "stage_provider_id_external_key" to table: "stages"
CREATE UNIQUE INDEX stage_provider_id_external_key ON stages (provider_id, external_key);
-- Create index "stage_season_id" to table: "stages"
CREATE INDEX stage_season_id ON stages (season_id);
-- Create "standings" table
CREATE TABLE standings (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  created_at datetime NOT NULL DEFAULT (DATETIME('now')),
  updated_at datetime NOT NULL DEFAULT (DATETIME('now')),
  stage_id blob NOT NULL,
  participant_id blob NOT NULL,
  rank integer NOT NULL DEFAULT 0,
  points integer NOT NULL DEFAULT 0
);
-- Create index "standing_stage_id_participant_id" to table: "standings"
CREATE UNIQUE INDEX standing_stage_id_participant_id ON standings (stage_id, participant_id);
-- Create index "standing_stage_id_rank" to table: "standings"
CREATE INDEX standing_stage_id_rank ON standings (stage_id, rank);
-- Create "sync_states" table
CREATE TABLE sync_states (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  provider_id text NOT NULL,
  scope text NOT NULL,
  cursor text NULL,
  last_synced_at datetime NULL,
  last_error text NULL
);
-- Create index "syncstate_provider_id_scope" to table: "sync_states"
CREATE UNIQUE INDEX syncstate_provider_id_scope ON sync_states (provider_id, scope);
-- Create "venues" table
CREATE TABLE venues (
  id blob NOT NULL,
  created_at datetime NOT NULL DEFAULT (DATETIME('now')),
  updated_at datetime NOT NULL DEFAULT (DATETIME('now')),
  provider_id text NOT NULL,
  external_key text NOT NULL,
  name text NOT NULL,
  city text NULL,
  country_iso text NULL,
  PRIMARY KEY (id)
);
-- Create index "venue_provider_id_external_key" to table: "venues"
CREATE UNIQUE INDEX venue_provider_id_external_key ON venues (provider_id, external_key);
