-- Create "subscriptions" table
CREATE TABLE subscriptions (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  created_at datetime NOT NULL DEFAULT (DATETIME('now')),
  updated_at datetime NOT NULL DEFAULT (DATETIME('now')),
  guild_id text NOT NULL,
  user_id text NOT NULL,
  kind text NOT NULL,
  value text NULL
);
-- Create index "subscription_guild_id_user_id_kind_value" to table: "subscriptions"
CREATE UNIQUE INDEX subscription_guild_id_user_id_kind_value ON subscriptions (guild_id, user_id, kind, value);
-- Create index "subscription_guild_id" to table: "subscriptions"
CREATE INDEX subscription_guild_id ON subscriptions (guild_id);
