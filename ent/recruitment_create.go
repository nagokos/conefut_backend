// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nagokos/connefut_backend/ent/competition"
	"github.com/nagokos/connefut_backend/ent/prefecture"
	"github.com/nagokos/connefut_backend/ent/recruitment"
	"github.com/nagokos/connefut_backend/ent/user"
)

// RecruitmentCreate is the builder for creating a Recruitment entity.
type RecruitmentCreate struct {
	config
	mutation *RecruitmentMutation
	hooks    []Hook
}

// SetCreatedAt sets the "created_at" field.
func (rc *RecruitmentCreate) SetCreatedAt(t time.Time) *RecruitmentCreate {
	rc.mutation.SetCreatedAt(t)
	return rc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (rc *RecruitmentCreate) SetNillableCreatedAt(t *time.Time) *RecruitmentCreate {
	if t != nil {
		rc.SetCreatedAt(*t)
	}
	return rc
}

// SetUpdatedAt sets the "updated_at" field.
func (rc *RecruitmentCreate) SetUpdatedAt(t time.Time) *RecruitmentCreate {
	rc.mutation.SetUpdatedAt(t)
	return rc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (rc *RecruitmentCreate) SetNillableUpdatedAt(t *time.Time) *RecruitmentCreate {
	if t != nil {
		rc.SetUpdatedAt(*t)
	}
	return rc
}

// SetTitle sets the "title" field.
func (rc *RecruitmentCreate) SetTitle(s string) *RecruitmentCreate {
	rc.mutation.SetTitle(s)
	return rc
}

// SetType sets the "type" field.
func (rc *RecruitmentCreate) SetType(r recruitment.Type) *RecruitmentCreate {
	rc.mutation.SetType(r)
	return rc
}

// SetNillableType sets the "type" field if the given value is not nil.
func (rc *RecruitmentCreate) SetNillableType(r *recruitment.Type) *RecruitmentCreate {
	if r != nil {
		rc.SetType(*r)
	}
	return rc
}

// SetLevel sets the "level" field.
func (rc *RecruitmentCreate) SetLevel(r recruitment.Level) *RecruitmentCreate {
	rc.mutation.SetLevel(r)
	return rc
}

// SetNillableLevel sets the "level" field if the given value is not nil.
func (rc *RecruitmentCreate) SetNillableLevel(r *recruitment.Level) *RecruitmentCreate {
	if r != nil {
		rc.SetLevel(*r)
	}
	return rc
}

// SetPlace sets the "place" field.
func (rc *RecruitmentCreate) SetPlace(s string) *RecruitmentCreate {
	rc.mutation.SetPlace(s)
	return rc
}

// SetNillablePlace sets the "place" field if the given value is not nil.
func (rc *RecruitmentCreate) SetNillablePlace(s *string) *RecruitmentCreate {
	if s != nil {
		rc.SetPlace(*s)
	}
	return rc
}

// SetStartAt sets the "start_at" field.
func (rc *RecruitmentCreate) SetStartAt(t time.Time) *RecruitmentCreate {
	rc.mutation.SetStartAt(t)
	return rc
}

// SetNillableStartAt sets the "start_at" field if the given value is not nil.
func (rc *RecruitmentCreate) SetNillableStartAt(t *time.Time) *RecruitmentCreate {
	if t != nil {
		rc.SetStartAt(*t)
	}
	return rc
}

// SetContent sets the "content" field.
func (rc *RecruitmentCreate) SetContent(s string) *RecruitmentCreate {
	rc.mutation.SetContent(s)
	return rc
}

// SetLocationURL sets the "Location_url" field.
func (rc *RecruitmentCreate) SetLocationURL(s string) *RecruitmentCreate {
	rc.mutation.SetLocationURL(s)
	return rc
}

// SetNillableLocationURL sets the "Location_url" field if the given value is not nil.
func (rc *RecruitmentCreate) SetNillableLocationURL(s *string) *RecruitmentCreate {
	if s != nil {
		rc.SetLocationURL(*s)
	}
	return rc
}

// SetCapacity sets the "capacity" field.
func (rc *RecruitmentCreate) SetCapacity(i int) *RecruitmentCreate {
	rc.mutation.SetCapacity(i)
	return rc
}

// SetNillableCapacity sets the "capacity" field if the given value is not nil.
func (rc *RecruitmentCreate) SetNillableCapacity(i *int) *RecruitmentCreate {
	if i != nil {
		rc.SetCapacity(*i)
	}
	return rc
}

// SetClosingAt sets the "closing_at" field.
func (rc *RecruitmentCreate) SetClosingAt(t time.Time) *RecruitmentCreate {
	rc.mutation.SetClosingAt(t)
	return rc
}

// SetID sets the "id" field.
func (rc *RecruitmentCreate) SetID(s string) *RecruitmentCreate {
	rc.mutation.SetID(s)
	return rc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (rc *RecruitmentCreate) SetNillableID(s *string) *RecruitmentCreate {
	if s != nil {
		rc.SetID(*s)
	}
	return rc
}

// SetUserID sets the "user" edge to the User entity by ID.
func (rc *RecruitmentCreate) SetUserID(id string) *RecruitmentCreate {
	rc.mutation.SetUserID(id)
	return rc
}

// SetUser sets the "user" edge to the User entity.
func (rc *RecruitmentCreate) SetUser(u *User) *RecruitmentCreate {
	return rc.SetUserID(u.ID)
}

// SetPrefectureID sets the "prefecture" edge to the Prefecture entity by ID.
func (rc *RecruitmentCreate) SetPrefectureID(id string) *RecruitmentCreate {
	rc.mutation.SetPrefectureID(id)
	return rc
}

// SetPrefecture sets the "prefecture" edge to the Prefecture entity.
func (rc *RecruitmentCreate) SetPrefecture(p *Prefecture) *RecruitmentCreate {
	return rc.SetPrefectureID(p.ID)
}

// SetCompetitionID sets the "competition" edge to the Competition entity by ID.
func (rc *RecruitmentCreate) SetCompetitionID(id string) *RecruitmentCreate {
	rc.mutation.SetCompetitionID(id)
	return rc
}

// SetCompetition sets the "competition" edge to the Competition entity.
func (rc *RecruitmentCreate) SetCompetition(c *Competition) *RecruitmentCreate {
	return rc.SetCompetitionID(c.ID)
}

// Mutation returns the RecruitmentMutation object of the builder.
func (rc *RecruitmentCreate) Mutation() *RecruitmentMutation {
	return rc.mutation
}

// Save creates the Recruitment in the database.
func (rc *RecruitmentCreate) Save(ctx context.Context) (*Recruitment, error) {
	var (
		err  error
		node *Recruitment
	)
	rc.defaults()
	if len(rc.hooks) == 0 {
		if err = rc.check(); err != nil {
			return nil, err
		}
		node, err = rc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*RecruitmentMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = rc.check(); err != nil {
				return nil, err
			}
			rc.mutation = mutation
			if node, err = rc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(rc.hooks) - 1; i >= 0; i-- {
			if rc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = rc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (rc *RecruitmentCreate) SaveX(ctx context.Context) *Recruitment {
	v, err := rc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rc *RecruitmentCreate) Exec(ctx context.Context) error {
	_, err := rc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rc *RecruitmentCreate) ExecX(ctx context.Context) {
	if err := rc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (rc *RecruitmentCreate) defaults() {
	if _, ok := rc.mutation.CreatedAt(); !ok {
		v := recruitment.DefaultCreatedAt()
		rc.mutation.SetCreatedAt(v)
	}
	if _, ok := rc.mutation.UpdatedAt(); !ok {
		v := recruitment.DefaultUpdatedAt()
		rc.mutation.SetUpdatedAt(v)
	}
	if _, ok := rc.mutation.GetType(); !ok {
		v := recruitment.DefaultType
		rc.mutation.SetType(v)
	}
	if _, ok := rc.mutation.ID(); !ok {
		v := recruitment.DefaultID()
		rc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rc *RecruitmentCreate) check() error {
	if _, ok := rc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "created_at"`)}
	}
	if _, ok := rc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "updated_at"`)}
	}
	if _, ok := rc.mutation.Title(); !ok {
		return &ValidationError{Name: "title", err: errors.New(`ent: missing required field "title"`)}
	}
	if v, ok := rc.mutation.Title(); ok {
		if err := recruitment.TitleValidator(v); err != nil {
			return &ValidationError{Name: "title", err: fmt.Errorf(`ent: validator failed for field "title": %w`, err)}
		}
	}
	if _, ok := rc.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "type"`)}
	}
	if v, ok := rc.mutation.GetType(); ok {
		if err := recruitment.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "type": %w`, err)}
		}
	}
	if v, ok := rc.mutation.Level(); ok {
		if err := recruitment.LevelValidator(v); err != nil {
			return &ValidationError{Name: "level", err: fmt.Errorf(`ent: validator failed for field "level": %w`, err)}
		}
	}
	if _, ok := rc.mutation.Content(); !ok {
		return &ValidationError{Name: "content", err: errors.New(`ent: missing required field "content"`)}
	}
	if v, ok := rc.mutation.Content(); ok {
		if err := recruitment.ContentValidator(v); err != nil {
			return &ValidationError{Name: "content", err: fmt.Errorf(`ent: validator failed for field "content": %w`, err)}
		}
	}
	if _, ok := rc.mutation.ClosingAt(); !ok {
		return &ValidationError{Name: "closing_at", err: errors.New(`ent: missing required field "closing_at"`)}
	}
	if v, ok := rc.mutation.ID(); ok {
		if err := recruitment.IDValidator(v); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`ent: validator failed for field "id": %w`, err)}
		}
	}
	if _, ok := rc.mutation.UserID(); !ok {
		return &ValidationError{Name: "user", err: errors.New("ent: missing required edge \"user\"")}
	}
	if _, ok := rc.mutation.PrefectureID(); !ok {
		return &ValidationError{Name: "prefecture", err: errors.New("ent: missing required edge \"prefecture\"")}
	}
	if _, ok := rc.mutation.CompetitionID(); !ok {
		return &ValidationError{Name: "competition", err: errors.New("ent: missing required edge \"competition\"")}
	}
	return nil
}

func (rc *RecruitmentCreate) sqlSave(ctx context.Context) (*Recruitment, error) {
	_node, _spec := rc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		_node.ID = _spec.ID.Value.(string)
	}
	return _node, nil
}

func (rc *RecruitmentCreate) createSpec() (*Recruitment, *sqlgraph.CreateSpec) {
	var (
		_node = &Recruitment{config: rc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: recruitment.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: recruitment.FieldID,
			},
		}
	)
	if id, ok := rc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := rc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: recruitment.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := rc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: recruitment.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if value, ok := rc.mutation.Title(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: recruitment.FieldTitle,
		})
		_node.Title = value
	}
	if value, ok := rc.mutation.GetType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: recruitment.FieldType,
		})
		_node.Type = value
	}
	if value, ok := rc.mutation.Level(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: recruitment.FieldLevel,
		})
		_node.Level = value
	}
	if value, ok := rc.mutation.Place(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: recruitment.FieldPlace,
		})
		_node.Place = value
	}
	if value, ok := rc.mutation.StartAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: recruitment.FieldStartAt,
		})
		_node.StartAt = value
	}
	if value, ok := rc.mutation.Content(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: recruitment.FieldContent,
		})
		_node.Content = value
	}
	if value, ok := rc.mutation.LocationURL(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: recruitment.FieldLocationURL,
		})
		_node.LocationURL = value
	}
	if value, ok := rc.mutation.Capacity(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: recruitment.FieldCapacity,
		})
		_node.Capacity = value
	}
	if value, ok := rc.mutation.ClosingAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: recruitment.FieldClosingAt,
		})
		_node.ClosingAt = value
	}
	if nodes := rc.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   recruitment.UserTable,
			Columns: []string{recruitment.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.user_id = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.PrefectureIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   recruitment.PrefectureTable,
			Columns: []string{recruitment.PrefectureColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: prefecture.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.prefecture_id = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rc.mutation.CompetitionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   recruitment.CompetitionTable,
			Columns: []string{recruitment.CompetitionColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: competition.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.competition_id = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// RecruitmentCreateBulk is the builder for creating many Recruitment entities in bulk.
type RecruitmentCreateBulk struct {
	config
	builders []*RecruitmentCreate
}

// Save creates the Recruitment entities in the database.
func (rcb *RecruitmentCreateBulk) Save(ctx context.Context) ([]*Recruitment, error) {
	specs := make([]*sqlgraph.CreateSpec, len(rcb.builders))
	nodes := make([]*Recruitment, len(rcb.builders))
	mutators := make([]Mutator, len(rcb.builders))
	for i := range rcb.builders {
		func(i int, root context.Context) {
			builder := rcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*RecruitmentMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, rcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, rcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rcb *RecruitmentCreateBulk) SaveX(ctx context.Context) []*Recruitment {
	v, err := rcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rcb *RecruitmentCreateBulk) Exec(ctx context.Context) error {
	_, err := rcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rcb *RecruitmentCreateBulk) ExecX(ctx context.Context) {
	if err := rcb.Exec(ctx); err != nil {
		panic(err)
	}
}
