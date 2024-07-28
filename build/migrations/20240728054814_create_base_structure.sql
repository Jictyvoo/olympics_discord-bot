-- Create "competitors" table
CREATE TABLE competitors (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, name text NOT NULL, code text NOT NULL, country_id integer NOT NULL, CONSTRAINT competitors_country_infos_competitors FOREIGN KEY (country_id) REFERENCES country_infos (id) ON DELETE NO ACTION);
-- Create index "competitors_country_id" to table: "competitors"
CREATE INDEX competitors_country_id ON competitors (country_id);
-- Create index "competitors_code" to table: "competitors"
CREATE INDEX competitors_code ON competitors (code);
-- Create index "competitors_name_code" to table: "competitors"
CREATE UNIQUE INDEX competitors_name_code ON competitors (name, code);
-- Create index "competitors_country_id_name_code" to table: "competitors"
CREATE UNIQUE INDEX competitors_country_id_name_code ON competitors (country_id, name, code);
-- Create "country_infos" table
CREATE TABLE country_infos (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, created_at datetime NOT NULL, updated_at datetime NOT NULL, code text NOT NULL, name text NOT NULL, code_num text NOT NULL, iso_code_len2 text NULL, iso_code_len3 text NOT NULL, ioc_code text NOT NULL, population integer NULL, area_km2 real NULL, gdp_usd text NULL);
-- Create index "country_infos_ioc_code_key" to table: "country_infos"
CREATE UNIQUE INDEX country_infos_ioc_code_key ON country_infos (ioc_code);
-- Create index "countryinfo_id" to table: "country_infos"
CREATE UNIQUE INDEX countryinfo_id ON country_infos (id);
-- Create index "countryinfo_ioc_code" to table: "country_infos"
CREATE INDEX countryinfo_ioc_code ON country_infos (ioc_code);
-- Create "olympic_disciplines" table
CREATE TABLE olympic_disciplines (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, name text NOT NULL, description text NULL);
-- Create index "olympic_disciplines_name_key" to table: "olympic_disciplines"
CREATE UNIQUE INDEX olympic_disciplines_name_key ON olympic_disciplines (name);
-- Create index "olympicdiscipline_id" to table: "olympic_disciplines"
CREATE UNIQUE INDEX olympicdiscipline_id ON olympic_disciplines (id);
-- Create index "olympicdiscipline_name" to table: "olympic_disciplines"
CREATE UNIQUE INDEX olympicdiscipline_name ON olympic_disciplines (name);
-- Create "olympic_events" table
CREATE TABLE olympic_events (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, event_name text NOT NULL, phase text NOT NULL, gender integer NOT NULL, start_at datetime NOT NULL, end_at datetime NOT NULL, status text NOT NULL, discipline_id integer NOT NULL, CONSTRAINT olympic_events_olympic_disciplines_olympic_events FOREIGN KEY (discipline_id) REFERENCES olympic_disciplines (id) ON DELETE NO ACTION);
-- Create index "olympicevent_event_name_discipline_id_phase_gender" to table: "olympic_events"
CREATE UNIQUE INDEX olympicevent_event_name_discipline_id_phase_gender ON olympic_events (event_name, discipline_id, phase, gender);
-- Create "results" table
CREATE TABLE results (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, position text NULL, mark text NULL, medal_type text NULL, irm text NOT NULL, competitor_id integer NOT NULL, event_id integer NOT NULL, CONSTRAINT results_competitors_results FOREIGN KEY (competitor_id) REFERENCES competitors (id) ON DELETE NO ACTION, CONSTRAINT results_olympic_events_results FOREIGN KEY (event_id) REFERENCES olympic_events (id) ON DELETE NO ACTION);
-- Create index "results_competitor_id_event_id" to table: "results"
CREATE UNIQUE INDEX results_competitor_id_event_id ON results (competitor_id, event_id);
