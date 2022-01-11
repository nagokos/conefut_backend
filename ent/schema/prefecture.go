package schema

import (
	"entgo.io/ent"
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
	return nil
}
