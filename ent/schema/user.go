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

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
		UUIDMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(50)",
			}).
			Validate(validation.CheckStringLen(50)),
		field.String("email").
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(100)",
			}).
			Validate(validation.CheckStringLen(100)).
			Unique(),
		field.Enum("role").
			Values(
				"admin",
				"general",
			).
			Default("general"),
		field.String("avatar").
			Default("https://abs.twimg.com/sticky/default_profile_images/default_profile.png"),
		field.String("introduction").
			Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(4000)",
			}).
			Validate(validation.CheckStringLen(4000)),
		field.Enum("email_verification_status").
			Values(
				"unnecessary",
				"pending",
				"verified",
			).
			Default("pending"),
		field.String("email_verification_token").
			Optional(),
		field.Time("email_verification_token_expires_at").
			Optional(),
		field.String("password_digest").
			Optional(),
		field.Time("last_sign_in_at").
			Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("recruitments", Recruitment.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}).
			StorageKey(edge.Column("user_id")),
		edge.To("stocks", Stock.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}).
			StorageKey(edge.Column("user_id")),
		edge.To("applicants", Applicant.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}).
			StorageKey(edge.Column("user_id")),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email_verification_token"),
	}
}
