package entschema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/jictyvoo/olhojogo/database/entschema/customixins"
	"github.com/jictyvoo/olhojogo/database/entschema/uuidfield"
)

// Fixture is a scheduled match/event between participants within a Stage.
type Fixture struct{ ent.Schema }

func (Fixture) Fields() []ent.Field {
	return []ent.Field{
		uuidfield.DeterministicID(),
		field.String("provider_id").NotEmpty(),
		field.String("external_key").NotEmpty(),
		uuidfield.Field("stage_id"),
		uuidfield.Optional("group_id"),
		uuidfield.Optional("venue_id"),
		field.String("name").NotEmpty(),
		field.Time("starts_at").Default(time.Now),
		field.Time("ends_at").Default(time.Now),
		field.String("status").NotEmpty().Default("scheduled"),
		field.String("checksum").MaxLen(64).Optional(),
	}
}

func (Fixture) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider_id", "external_key").Unique(),
		index.Fields("stage_id"),
		index.Fields("starts_at"),
		index.Fields("status"),
	}
}

func (Fixture) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("stage", Stage.Type).Ref("fixtures").Field("stage_id").Unique().Required(),
		edge.To("fixture_participants", FixtureParticipant.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("results", Result.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("alerts", Alert.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("discord_events", DiscordEvent.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (Fixture) Mixin() []ent.Mixin { return []ent.Mixin{customixins.TimestampsMixin{}} }
