-- Create "notified_events" table
CREATE TABLE notified_events (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, event_sha256 text NOT NULL, status text NOT NULL, notified_at datetime NULL, event_id integer NOT NULL, CONSTRAINT notified_events_olympic_events_notified_events FOREIGN KEY (event_id) REFERENCES olympic_events (id) ON DELETE NO ACTION);
-- Create index "notified_events_event_sha256_key" to table: "notified_events"
CREATE UNIQUE INDEX notified_events_event_sha256_key ON notified_events (event_sha256);
