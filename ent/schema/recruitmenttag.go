package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// RecruitmentTag holds the schema definition for the RecruitmentTag entity.
type RecruitmentTag struct {
	ent.Schema
}

func (RecruitmentTag) Mixin() []ent.Mixin {
	return []ent.Mixin{
		UUIDMixin{},
		TimeMixin{},
	}
}

// Fields of the RecruitmentTag.
func (RecruitmentTag) Fields() []ent.Field {
	return []ent.Field{
		field.String("recruitment_id"),
		field.String("tag_id"),
	}
}

// Edges of the RecruitmentTag.
func (RecruitmentTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("recruitment", Recruitment.Type).
			Unique().
			Required().
			Field("recruitment_id").
			Ref("recruitment_tags"),
		edge.From("tag", Tag.Type).
			Unique().
			Field("tag_id").
			Required().
			Ref("recruitment_tags"),
	}
}

func (RecruitmentTag) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("recruitment_id", "tag_id").Unique(),
		index.Fields("tag_id"),
		index.Fields("recruitment_id"),
	}
}
