-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_competitors" table
CREATE TABLE new_competitors (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, name text NOT NULL, code text NOT NULL, country_id integer NOT NULL, CONSTRAINT competitors_country_infos_competitors FOREIGN KEY (country_id) REFERENCES country_infos (id) ON DELETE RESTRICT);
-- Copy rows from old table "competitors" to new temporary table "new_competitors"
INSERT INTO new_competitors (id, name, code, country_id) SELECT id, name, code, country_id FROM competitors;
-- Drop "competitors" table after copying rows
DROP TABLE competitors;
-- Rename temporary table "new_competitors" to "competitors"
ALTER TABLE new_competitors RENAME TO competitors;
-- Create index "competitors_country_id" to table: "competitors"
CREATE INDEX competitors_country_id ON competitors (country_id);
-- Create index "competitors_code" to table: "competitors"
CREATE INDEX competitors_code ON competitors (code);
-- Create index "competitors_name_code" to table: "competitors"
CREATE UNIQUE INDEX competitors_name_code ON competitors (name, code);
-- Create index "competitors_country_id_name_code" to table: "competitors"
CREATE UNIQUE INDEX competitors_country_id_name_code ON competitors (country_id, name, code);
-- Create "new_country_infos" table
CREATE TABLE new_country_infos (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at datetime NOT NULL DEFAULT (DATETIME('now')), updated_at datetime NOT NULL DEFAULT (DATETIME('now')), code text NOT NULL, name text NOT NULL, code_num text NOT NULL, iso_code_len2 text NULL, iso_code_len3 text NOT NULL, ioc_code text NOT NULL, population integer NULL, area_km2 real NULL, gdp_usd text NULL);
-- Copy rows from old table "country_infos" to new temporary table "new_country_infos"
INSERT INTO new_country_infos (id, created_at, updated_at, code, name, code_num, iso_code_len2, iso_code_len3, ioc_code, population, area_km2, gdp_usd) SELECT id, IFNULL(created_at, (DATETIME('now'))) AS created_at, IFNULL(updated_at, (DATETIME('now'))) AS updated_at, code, name, code_num, iso_code_len2, iso_code_len3, ioc_code, population, area_km2, gdp_usd FROM country_infos;
-- Drop "country_infos" table after copying rows
DROP TABLE country_infos;
-- Rename temporary table "new_country_infos" to "country_infos"
ALTER TABLE new_country_infos RENAME TO country_infos;
-- Create index "country_infos_ioc_code_key" to table: "country_infos"
CREATE UNIQUE INDEX country_infos_ioc_code_key ON country_infos (ioc_code);
-- Create index "countryinfo_id" to table: "country_infos"
CREATE UNIQUE INDEX countryinfo_id ON country_infos (id);
-- Create index "countryinfo_ioc_code" to table: "country_infos"
CREATE INDEX countryinfo_ioc_code ON country_infos (ioc_code);
-- Create "new_olympic_events" table
CREATE TABLE new_olympic_events (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at datetime NOT NULL DEFAULT (DATETIME('now')), updated_at datetime NOT NULL DEFAULT (DATETIME('now')), event_name text NOT NULL, phase text NOT NULL, gender integer NOT NULL, session_code text NOT NULL, has_medal bool NOT NULL DEFAULT (false), start_at datetime NOT NULL, end_at datetime NOT NULL, status text NOT NULL, discipline_id integer NOT NULL, CONSTRAINT olympic_events_olympic_disciplines_olympic_events FOREIGN KEY (discipline_id) REFERENCES olympic_disciplines (id) ON DELETE RESTRICT);
-- Copy rows from old table "olympic_events" to new temporary table "new_olympic_events"
INSERT INTO new_olympic_events (id, event_name, phase, gender, session_code, start_at, end_at, status, discipline_id) SELECT id, event_name, phase, gender, session_code, start_at, end_at, status, discipline_id FROM olympic_events;
-- Drop "olympic_events" table after copying rows
DROP TABLE olympic_events;
-- Rename temporary table "new_olympic_events" to "olympic_events"
ALTER TABLE new_olympic_events RENAME TO olympic_events;
-- Create index "olympicevent_event_name_discipline_id_phase_gender_session_code" to table: "olympic_events"
CREATE UNIQUE INDEX olympicevent_event_name_discipline_id_phase_gender_session_code ON olympic_events (event_name, discipline_id, phase, gender, session_code);
-- Create "new_results" table
CREATE TABLE new_results (id uuid NOT NULL, created_at datetime NOT NULL DEFAULT (DATETIME('now')), updated_at datetime NOT NULL DEFAULT (DATETIME('now')), position text NULL, mark text NULL, medal_type text NULL, irm text NOT NULL, competitor_id integer NOT NULL, event_id integer NOT NULL, PRIMARY KEY (id), CONSTRAINT results_competitors_results FOREIGN KEY (competitor_id) REFERENCES competitors (id) ON DELETE CASCADE, CONSTRAINT results_olympic_events_results FOREIGN KEY (event_id) REFERENCES olympic_events (id) ON DELETE CASCADE);
-- Copy rows from old table "results" to new temporary table "new_results"
INSERT INTO new_results (id, position, mark, medal_type, irm, competitor_id, event_id) SELECT id, position, mark, medal_type, irm, competitor_id, event_id FROM results;
-- Drop "results" table after copying rows
DROP TABLE results;
-- Rename temporary table "new_results" to "results"
ALTER TABLE new_results RENAME TO results;
-- Create index "results_competitor_id_event_id" to table: "results"
CREATE UNIQUE INDEX results_competitor_id_event_id ON results (competitor_id, event_id);
-- Create "new_notified_events" table
CREATE TABLE new_notified_events (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, event_sha256 text NOT NULL, status text NOT NULL, notified_at datetime NULL, event_id integer NOT NULL, CONSTRAINT notified_events_olympic_events_notified_events FOREIGN KEY (event_id) REFERENCES olympic_events (id) ON DELETE CASCADE);
-- Copy rows from old table "notified_events" to new temporary table "new_notified_events"
INSERT INTO new_notified_events (id, event_sha256, status, notified_at, event_id) SELECT id, event_sha256, status, notified_at, event_id FROM notified_events;
-- Drop "notified_events" table after copying rows
DROP TABLE notified_events;
-- Rename temporary table "new_notified_events" to "notified_events"
ALTER TABLE new_notified_events RENAME TO notified_events;
-- Create index "notified_events_event_sha256_key" to table: "notified_events"
CREATE UNIQUE INDEX notified_events_event_sha256_key ON notified_events (event_sha256);
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
