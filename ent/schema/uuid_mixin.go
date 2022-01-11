package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/rs/xid"
)

type UUIDMixin struct {
	mixin.Schema
}

func (UUIDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(func() string {
				return xid.New().String()
			}).
			NotEmpty().
			Immutable().
			Unique(),
	}
}
