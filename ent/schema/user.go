package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"github.com/rs/xid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				return xid.New().String()
			}).
			NotEmpty().
			Immutable().
			Unique(),
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
		field.String("password_digest").
			Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
