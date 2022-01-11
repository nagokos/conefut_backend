package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
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
			MaxLen(50),
		field.String("email").
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(100)",
			}).
			MaxLen(100).
			Unique(),
		field.Enum("role").
			Values("admin", "general").
			Default("general"),
		field.String("avatar").
			Default("https://abs.twimg.com/sticky/default_profile_images/default_profile.png"),
		field.String("introduction").
			Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "varchar(4000)",
			}).
			MaxLen(4000),
		field.Bool("email_verification_status").
			Default(false),
		field.String("email_verification_token").
			Optional().
			Unique(),
		field.Time("email_verification_token_expires_at").
			Optional(),
		field.String("password_digest").
			Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
