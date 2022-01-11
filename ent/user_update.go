// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nagokos/connefut_backend/ent/predicate"
	"github.com/nagokos/connefut_backend/ent/user"
)

// UserUpdate is the builder for updating User entities.
type UserUpdate struct {
	config
	hooks    []Hook
	mutation *UserMutation
}

// Where appends a list predicates to the UserUpdate builder.
func (uu *UserUpdate) Where(ps ...predicate.User) *UserUpdate {
	uu.mutation.Where(ps...)
	return uu
}

// SetUpdatedAt sets the "updated_at" field.
func (uu *UserUpdate) SetUpdatedAt(t time.Time) *UserUpdate {
	uu.mutation.SetUpdatedAt(t)
	return uu
}

// SetName sets the "name" field.
func (uu *UserUpdate) SetName(s string) *UserUpdate {
	uu.mutation.SetName(s)
	return uu
}

// SetEmail sets the "email" field.
func (uu *UserUpdate) SetEmail(s string) *UserUpdate {
	uu.mutation.SetEmail(s)
	return uu
}

// SetRole sets the "role" field.
func (uu *UserUpdate) SetRole(u user.Role) *UserUpdate {
	uu.mutation.SetRole(u)
	return uu
}

// SetNillableRole sets the "role" field if the given value is not nil.
func (uu *UserUpdate) SetNillableRole(u *user.Role) *UserUpdate {
	if u != nil {
		uu.SetRole(*u)
	}
	return uu
}

// SetAvatar sets the "avatar" field.
func (uu *UserUpdate) SetAvatar(s string) *UserUpdate {
	uu.mutation.SetAvatar(s)
	return uu
}

// SetNillableAvatar sets the "avatar" field if the given value is not nil.
func (uu *UserUpdate) SetNillableAvatar(s *string) *UserUpdate {
	if s != nil {
		uu.SetAvatar(*s)
	}
	return uu
}

// SetIntroduction sets the "introduction" field.
func (uu *UserUpdate) SetIntroduction(s string) *UserUpdate {
	uu.mutation.SetIntroduction(s)
	return uu
}

// SetNillableIntroduction sets the "introduction" field if the given value is not nil.
func (uu *UserUpdate) SetNillableIntroduction(s *string) *UserUpdate {
	if s != nil {
		uu.SetIntroduction(*s)
	}
	return uu
}

// ClearIntroduction clears the value of the "introduction" field.
func (uu *UserUpdate) ClearIntroduction() *UserUpdate {
	uu.mutation.ClearIntroduction()
	return uu
}

// SetEmailVerificationStatus sets the "email_verification_status" field.
func (uu *UserUpdate) SetEmailVerificationStatus(b bool) *UserUpdate {
	uu.mutation.SetEmailVerificationStatus(b)
	return uu
}

// SetNillableEmailVerificationStatus sets the "email_verification_status" field if the given value is not nil.
func (uu *UserUpdate) SetNillableEmailVerificationStatus(b *bool) *UserUpdate {
	if b != nil {
		uu.SetEmailVerificationStatus(*b)
	}
	return uu
}

// SetEmailVerificationToken sets the "email_verification_token" field.
func (uu *UserUpdate) SetEmailVerificationToken(s string) *UserUpdate {
	uu.mutation.SetEmailVerificationToken(s)
	return uu
}

// SetNillableEmailVerificationToken sets the "email_verification_token" field if the given value is not nil.
func (uu *UserUpdate) SetNillableEmailVerificationToken(s *string) *UserUpdate {
	if s != nil {
		uu.SetEmailVerificationToken(*s)
	}
	return uu
}

// ClearEmailVerificationToken clears the value of the "email_verification_token" field.
func (uu *UserUpdate) ClearEmailVerificationToken() *UserUpdate {
	uu.mutation.ClearEmailVerificationToken()
	return uu
}

// SetEmailVerificationTokenExpiresAt sets the "email_verification_token_expires_at" field.
func (uu *UserUpdate) SetEmailVerificationTokenExpiresAt(t time.Time) *UserUpdate {
	uu.mutation.SetEmailVerificationTokenExpiresAt(t)
	return uu
}

// SetNillableEmailVerificationTokenExpiresAt sets the "email_verification_token_expires_at" field if the given value is not nil.
func (uu *UserUpdate) SetNillableEmailVerificationTokenExpiresAt(t *time.Time) *UserUpdate {
	if t != nil {
		uu.SetEmailVerificationTokenExpiresAt(*t)
	}
	return uu
}

// ClearEmailVerificationTokenExpiresAt clears the value of the "email_verification_token_expires_at" field.
func (uu *UserUpdate) ClearEmailVerificationTokenExpiresAt() *UserUpdate {
	uu.mutation.ClearEmailVerificationTokenExpiresAt()
	return uu
}

// SetPasswordDigest sets the "password_digest" field.
func (uu *UserUpdate) SetPasswordDigest(s string) *UserUpdate {
	uu.mutation.SetPasswordDigest(s)
	return uu
}

// SetNillablePasswordDigest sets the "password_digest" field if the given value is not nil.
func (uu *UserUpdate) SetNillablePasswordDigest(s *string) *UserUpdate {
	if s != nil {
		uu.SetPasswordDigest(*s)
	}
	return uu
}

// ClearPasswordDigest clears the value of the "password_digest" field.
func (uu *UserUpdate) ClearPasswordDigest() *UserUpdate {
	uu.mutation.ClearPasswordDigest()
	return uu
}

// Mutation returns the UserMutation object of the builder.
func (uu *UserUpdate) Mutation() *UserMutation {
	return uu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (uu *UserUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	uu.defaults()
	if len(uu.hooks) == 0 {
		if err = uu.check(); err != nil {
			return 0, err
		}
		affected, err = uu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UserMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = uu.check(); err != nil {
				return 0, err
			}
			uu.mutation = mutation
			affected, err = uu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(uu.hooks) - 1; i >= 0; i-- {
			if uu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = uu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, uu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (uu *UserUpdate) SaveX(ctx context.Context) int {
	affected, err := uu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (uu *UserUpdate) Exec(ctx context.Context) error {
	_, err := uu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uu *UserUpdate) ExecX(ctx context.Context) {
	if err := uu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uu *UserUpdate) defaults() {
	if _, ok := uu.mutation.UpdatedAt(); !ok {
		v := user.UpdateDefaultUpdatedAt()
		uu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uu *UserUpdate) check() error {
	if v, ok := uu.mutation.Name(); ok {
		if err := user.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	if v, ok := uu.mutation.Email(); ok {
		if err := user.EmailValidator(v); err != nil {
			return &ValidationError{Name: "email", err: fmt.Errorf("ent: validator failed for field \"email\": %w", err)}
		}
	}
	if v, ok := uu.mutation.Role(); ok {
		if err := user.RoleValidator(v); err != nil {
			return &ValidationError{Name: "role", err: fmt.Errorf("ent: validator failed for field \"role\": %w", err)}
		}
	}
	if v, ok := uu.mutation.Introduction(); ok {
		if err := user.IntroductionValidator(v); err != nil {
			return &ValidationError{Name: "introduction", err: fmt.Errorf("ent: validator failed for field \"introduction\": %w", err)}
		}
	}
	return nil
}

func (uu *UserUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   user.Table,
			Columns: user.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: user.FieldID,
			},
		},
	}
	if ps := uu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: user.FieldUpdatedAt,
		})
	}
	if value, ok := uu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldName,
		})
	}
	if value, ok := uu.mutation.Email(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldEmail,
		})
	}
	if value, ok := uu.mutation.Role(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: user.FieldRole,
		})
	}
	if value, ok := uu.mutation.Avatar(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldAvatar,
		})
	}
	if value, ok := uu.mutation.Introduction(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldIntroduction,
		})
	}
	if uu.mutation.IntroductionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldIntroduction,
		})
	}
	if value, ok := uu.mutation.EmailVerificationStatus(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: user.FieldEmailVerificationStatus,
		})
	}
	if value, ok := uu.mutation.EmailVerificationToken(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldEmailVerificationToken,
		})
	}
	if uu.mutation.EmailVerificationTokenCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldEmailVerificationToken,
		})
	}
	if value, ok := uu.mutation.EmailVerificationTokenExpiresAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: user.FieldEmailVerificationTokenExpiresAt,
		})
	}
	if uu.mutation.EmailVerificationTokenExpiresAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: user.FieldEmailVerificationTokenExpiresAt,
		})
	}
	if value, ok := uu.mutation.PasswordDigest(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldPasswordDigest,
		})
	}
	if uu.mutation.PasswordDigestCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldPasswordDigest,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, uu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// UserUpdateOne is the builder for updating a single User entity.
type UserUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *UserMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (uuo *UserUpdateOne) SetUpdatedAt(t time.Time) *UserUpdateOne {
	uuo.mutation.SetUpdatedAt(t)
	return uuo
}

// SetName sets the "name" field.
func (uuo *UserUpdateOne) SetName(s string) *UserUpdateOne {
	uuo.mutation.SetName(s)
	return uuo
}

// SetEmail sets the "email" field.
func (uuo *UserUpdateOne) SetEmail(s string) *UserUpdateOne {
	uuo.mutation.SetEmail(s)
	return uuo
}

// SetRole sets the "role" field.
func (uuo *UserUpdateOne) SetRole(u user.Role) *UserUpdateOne {
	uuo.mutation.SetRole(u)
	return uuo
}

// SetNillableRole sets the "role" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableRole(u *user.Role) *UserUpdateOne {
	if u != nil {
		uuo.SetRole(*u)
	}
	return uuo
}

// SetAvatar sets the "avatar" field.
func (uuo *UserUpdateOne) SetAvatar(s string) *UserUpdateOne {
	uuo.mutation.SetAvatar(s)
	return uuo
}

// SetNillableAvatar sets the "avatar" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableAvatar(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetAvatar(*s)
	}
	return uuo
}

// SetIntroduction sets the "introduction" field.
func (uuo *UserUpdateOne) SetIntroduction(s string) *UserUpdateOne {
	uuo.mutation.SetIntroduction(s)
	return uuo
}

// SetNillableIntroduction sets the "introduction" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableIntroduction(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetIntroduction(*s)
	}
	return uuo
}

// ClearIntroduction clears the value of the "introduction" field.
func (uuo *UserUpdateOne) ClearIntroduction() *UserUpdateOne {
	uuo.mutation.ClearIntroduction()
	return uuo
}

// SetEmailVerificationStatus sets the "email_verification_status" field.
func (uuo *UserUpdateOne) SetEmailVerificationStatus(b bool) *UserUpdateOne {
	uuo.mutation.SetEmailVerificationStatus(b)
	return uuo
}

// SetNillableEmailVerificationStatus sets the "email_verification_status" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableEmailVerificationStatus(b *bool) *UserUpdateOne {
	if b != nil {
		uuo.SetEmailVerificationStatus(*b)
	}
	return uuo
}

// SetEmailVerificationToken sets the "email_verification_token" field.
func (uuo *UserUpdateOne) SetEmailVerificationToken(s string) *UserUpdateOne {
	uuo.mutation.SetEmailVerificationToken(s)
	return uuo
}

// SetNillableEmailVerificationToken sets the "email_verification_token" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableEmailVerificationToken(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetEmailVerificationToken(*s)
	}
	return uuo
}

// ClearEmailVerificationToken clears the value of the "email_verification_token" field.
func (uuo *UserUpdateOne) ClearEmailVerificationToken() *UserUpdateOne {
	uuo.mutation.ClearEmailVerificationToken()
	return uuo
}

// SetEmailVerificationTokenExpiresAt sets the "email_verification_token_expires_at" field.
func (uuo *UserUpdateOne) SetEmailVerificationTokenExpiresAt(t time.Time) *UserUpdateOne {
	uuo.mutation.SetEmailVerificationTokenExpiresAt(t)
	return uuo
}

// SetNillableEmailVerificationTokenExpiresAt sets the "email_verification_token_expires_at" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableEmailVerificationTokenExpiresAt(t *time.Time) *UserUpdateOne {
	if t != nil {
		uuo.SetEmailVerificationTokenExpiresAt(*t)
	}
	return uuo
}

// ClearEmailVerificationTokenExpiresAt clears the value of the "email_verification_token_expires_at" field.
func (uuo *UserUpdateOne) ClearEmailVerificationTokenExpiresAt() *UserUpdateOne {
	uuo.mutation.ClearEmailVerificationTokenExpiresAt()
	return uuo
}

// SetPasswordDigest sets the "password_digest" field.
func (uuo *UserUpdateOne) SetPasswordDigest(s string) *UserUpdateOne {
	uuo.mutation.SetPasswordDigest(s)
	return uuo
}

// SetNillablePasswordDigest sets the "password_digest" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillablePasswordDigest(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetPasswordDigest(*s)
	}
	return uuo
}

// ClearPasswordDigest clears the value of the "password_digest" field.
func (uuo *UserUpdateOne) ClearPasswordDigest() *UserUpdateOne {
	uuo.mutation.ClearPasswordDigest()
	return uuo
}

// Mutation returns the UserMutation object of the builder.
func (uuo *UserUpdateOne) Mutation() *UserMutation {
	return uuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (uuo *UserUpdateOne) Select(field string, fields ...string) *UserUpdateOne {
	uuo.fields = append([]string{field}, fields...)
	return uuo
}

// Save executes the query and returns the updated User entity.
func (uuo *UserUpdateOne) Save(ctx context.Context) (*User, error) {
	var (
		err  error
		node *User
	)
	uuo.defaults()
	if len(uuo.hooks) == 0 {
		if err = uuo.check(); err != nil {
			return nil, err
		}
		node, err = uuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UserMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = uuo.check(); err != nil {
				return nil, err
			}
			uuo.mutation = mutation
			node, err = uuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(uuo.hooks) - 1; i >= 0; i-- {
			if uuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = uuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, uuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (uuo *UserUpdateOne) SaveX(ctx context.Context) *User {
	node, err := uuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (uuo *UserUpdateOne) Exec(ctx context.Context) error {
	_, err := uuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uuo *UserUpdateOne) ExecX(ctx context.Context) {
	if err := uuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uuo *UserUpdateOne) defaults() {
	if _, ok := uuo.mutation.UpdatedAt(); !ok {
		v := user.UpdateDefaultUpdatedAt()
		uuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uuo *UserUpdateOne) check() error {
	if v, ok := uuo.mutation.Name(); ok {
		if err := user.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	if v, ok := uuo.mutation.Email(); ok {
		if err := user.EmailValidator(v); err != nil {
			return &ValidationError{Name: "email", err: fmt.Errorf("ent: validator failed for field \"email\": %w", err)}
		}
	}
	if v, ok := uuo.mutation.Role(); ok {
		if err := user.RoleValidator(v); err != nil {
			return &ValidationError{Name: "role", err: fmt.Errorf("ent: validator failed for field \"role\": %w", err)}
		}
	}
	if v, ok := uuo.mutation.Introduction(); ok {
		if err := user.IntroductionValidator(v); err != nil {
			return &ValidationError{Name: "introduction", err: fmt.Errorf("ent: validator failed for field \"introduction\": %w", err)}
		}
	}
	return nil
}

func (uuo *UserUpdateOne) sqlSave(ctx context.Context) (_node *User, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   user.Table,
			Columns: user.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: user.FieldID,
			},
		},
	}
	id, ok := uuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing User.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := uuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, user.FieldID)
		for _, f := range fields {
			if !user.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != user.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := uuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uuo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: user.FieldUpdatedAt,
		})
	}
	if value, ok := uuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldName,
		})
	}
	if value, ok := uuo.mutation.Email(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldEmail,
		})
	}
	if value, ok := uuo.mutation.Role(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: user.FieldRole,
		})
	}
	if value, ok := uuo.mutation.Avatar(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldAvatar,
		})
	}
	if value, ok := uuo.mutation.Introduction(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldIntroduction,
		})
	}
	if uuo.mutation.IntroductionCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldIntroduction,
		})
	}
	if value, ok := uuo.mutation.EmailVerificationStatus(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: user.FieldEmailVerificationStatus,
		})
	}
	if value, ok := uuo.mutation.EmailVerificationToken(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldEmailVerificationToken,
		})
	}
	if uuo.mutation.EmailVerificationTokenCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldEmailVerificationToken,
		})
	}
	if value, ok := uuo.mutation.EmailVerificationTokenExpiresAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: user.FieldEmailVerificationTokenExpiresAt,
		})
	}
	if uuo.mutation.EmailVerificationTokenExpiresAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: user.FieldEmailVerificationTokenExpiresAt,
		})
	}
	if value, ok := uuo.mutation.PasswordDigest(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldPasswordDigest,
		})
	}
	if uuo.mutation.PasswordDigestCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldPasswordDigest,
		})
	}
	_node = &User{config: uuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, uuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
