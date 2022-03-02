package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Stock holds the schema definition for the Stock entity.
type Stock struct {
	ent.Schema
}

func (Stock) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeMixin{},
	}
}

// Fields of the Stock.
func (Stock) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id"),
		field.String("recruitment_id"),
	}
}

// Edges of the Stock.
func (Stock) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Unique().
			Field("user_id").
			Required().
			Ref("stocks"),
		edge.From("recruitment", Recruitment.Type).
			Unique().
			Field("recruitment_id").
			Required().
			Ref("stocks"),
	}
}

func (Stock) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "recruitment_id").Unique(),
	}
}
