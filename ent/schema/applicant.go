package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Applicant holds the schema definition for the Applicant entity.
type Applicant struct {
	ent.Schema
}

func (Applicant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
		UUIDMixin{},
	}
}

// Fields of the Applicant.
func (Applicant) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("management_status").
			Values(
				"backlog",
				"checked",
				"accepted",
				"rejected",
			).
			Default("backlog"),
		field.String("user_id"),
		field.String("recruitment_id"),
	}
}

// Edges of the Applicant.
func (Applicant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Unique().
			Field("user_id").
			Required().
			Ref("applicants"),
		edge.From("recruitment", Recruitment.Type).
			Unique().
			Field("recruitment_id").
			Required().
			Ref("applicants"),
	}
}

func (Applicant) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "recruitment_id").Unique(),
		index.Fields("user_id"),
		index.Fields("recruitment_id"),
	}
}
