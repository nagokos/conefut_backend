package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/nagokos/connefut_backend/ent/validation"
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
			Validate(validation.CheckStringLen(60)).
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(60)",
			}),
		field.Enum("type").
			Values(
				"unnecessary",
				"opponent",
				"individual",
				"teammate",
				"joining",
				"coaching",
				"others",
			).
			Default("unnecessary"),
		field.Enum("level").
			Values(
				"unnecessary",
				"enjoy",
				"beginner",
				"middle",
				"expert",
				"open",
			).
			Default("unnecessary"),
		field.String("place").
			Optional(),
		field.Time("start_at").
			Optional(),
		field.String("content").
			Validate(validation.CheckStringLen(10000)).
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(10000)",
			}).
			Optional().
			Comment("募集詳細"),
		field.Float("locationLat").
			Optional().
			Comment("会場の緯度"),
		field.Float("locationLng").
			Optional().
			Comment("会場の経度"),
		field.Int("capacity").
			Optional(),
		field.Time("closing_at").
			Optional().
			Comment("募集期限"),
		field.Enum("status").
			Values(
				"draft",
				"published",
				"closed",
			).
			Default("draft").
			Comment("募集のステータス"),
	}
}

// Edges of the Recruitment.
func (Recruitment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("stocks", Stock.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}).
			StorageKey(edge.Column("recruitment_id")),
		edge.From("user", User.Type).
			Unique().
			Required().
			Ref("recruitments"),
		edge.From("prefecture", Prefecture.Type).
			Unique().
			Ref("recruitments"),
		edge.From("competition", Competition.Type).
			Unique().
			Ref("recruitments"),
	}
}
