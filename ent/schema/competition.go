package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Competition holds the schema definition for the Competition entity.
type Competition struct {
	ent.Schema
}

func (Competition) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
		UUIDMixin{},
	}
}

// Fields of the Competition.
func (Competition) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
	}
}

// Edges of the Competition.
func (Competition) Edges() []ent.Edge {
	return nil
}
