package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Prefecture holds the schema definition for the Prefecture entity.
type Prefecture struct {
	ent.Schema
}

func (Prefecture) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
		UUIDMixin{},
	}
}

// Fields of the Prefecture.
func (Prefecture) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
	}
}

// Edges of the Prefecture.
func (Prefecture) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("recruitments", Recruitment.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Restrict,
			}).
			StorageKey(edge.Column("prefecture_id")),
	}
}
