package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Recruitment holds the schema definition for the Recruitment entity.
type Recruitment struct {
	ent.Schema
}

func (Recruitment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeMixin{},
	}
}

// Fields of the Recruitment.
func (Recruitment) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").
			MaxLen(60).
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(60)",
			}),
		field.Enum("type").
			Values(
				"opponent",
				"individual",
				"teammate",
				"joining",
				"coaching",
				"others",
			).
			Default("opponent"),
		field.Enum("level").
			Values(
				"enjoy",
				"beginner",
				"middle",
				"expert",
				"open",
			).
			Default("enjoy").
			Optional(),
		field.String("place").
			Optional(),
		field.Time("start_at").
			Optional(),
		field.String("content").
			MaxLen(10000).
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(10000)",
			}).
			Comment("募集詳細"),
		field.String("Location_url").
			Optional().
			Comment("会場の場所を埋め込むURL"),
		field.Int("capacity"),
		field.Time("closing_at").
			Comment("募集期限"),
	}
}

// Edges of the Recruitment.
func (Recruitment) Edges() []ent.Edge {
	return nil
}
