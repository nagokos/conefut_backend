package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
				"member",
				"joining",
				"others",
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
		field.String("prefecture_id").
			Optional(),
		field.String("competition_id").
			Optional(),
		field.String("user_id"),
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
		edge.To("applicants", Applicant.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}).
			StorageKey(edge.Column("recruitment_id")),
		edge.From("user", User.Type).
			Unique().
			Field("user_id").
			Required().
			Ref("recruitments"),
		edge.From("prefecture", Prefecture.Type).
			Unique().
			Field("prefecture_id").
			Ref("recruitments"),
		edge.From("competition", Competition.Type).
			Unique().
			Field("competition_id").
			Ref("recruitments"),
	}
}

func (Recruitment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("prefecture_id"),
		index.Fields("competition_id"),
	}
}
