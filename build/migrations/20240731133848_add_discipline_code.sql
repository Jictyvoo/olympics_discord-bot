-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_olympic_disciplines" table
CREATE TABLE new_olympic_disciplines (id integer NOT NULL PRIMARY KEY AUTOINCREMENT, name text NOT NULL, description text NULL, code text NOT NULL DEFAULT (''));
-- Copy rows from old table "olympic_disciplines" to new temporary table "new_olympic_disciplines"
INSERT INTO new_olympic_disciplines (id, name, description) SELECT id, name, description FROM olympic_disciplines;
-- Drop "olympic_disciplines" table after copying rows
DROP TABLE olympic_disciplines;
-- Rename temporary table "new_olympic_disciplines" to "olympic_disciplines"
ALTER TABLE new_olympic_disciplines RENAME TO olympic_disciplines;
-- Create index "olympic_disciplines_name_key" to table: "olympic_disciplines"
CREATE UNIQUE INDEX olympic_disciplines_name_key ON olympic_disciplines (name);
-- Create index "olympicdiscipline_id" to table: "olympic_disciplines"
CREATE UNIQUE INDEX olympicdiscipline_id ON olympic_disciplines (id);
-- Create index "olympicdiscipline_name" to table: "olympic_disciplines"
CREATE UNIQUE INDEX olympicdiscipline_name ON olympic_disciplines (name);
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
