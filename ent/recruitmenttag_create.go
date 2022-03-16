// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nagokos/connefut_backend/ent/recruitment"
	"github.com/nagokos/connefut_backend/ent/recruitmenttag"
	"github.com/nagokos/connefut_backend/ent/tag"
)

// RecruitmentTagCreate is the builder for creating a RecruitmentTag entity.
type RecruitmentTagCreate struct {
	config
	mutation *RecruitmentTagMutation
	hooks    []Hook
}

// SetCreatedAt sets the "created_at" field.
func (rtc *RecruitmentTagCreate) SetCreatedAt(t time.Time) *RecruitmentTagCreate {
	rtc.mutation.SetCreatedAt(t)
	return rtc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (rtc *RecruitmentTagCreate) SetNillableCreatedAt(t *time.Time) *RecruitmentTagCreate {
	if t != nil {
		rtc.SetCreatedAt(*t)
	}
	return rtc
}

// SetUpdatedAt sets the "updated_at" field.
func (rtc *RecruitmentTagCreate) SetUpdatedAt(t time.Time) *RecruitmentTagCreate {
	rtc.mutation.SetUpdatedAt(t)
	return rtc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (rtc *RecruitmentTagCreate) SetNillableUpdatedAt(t *time.Time) *RecruitmentTagCreate {
	if t != nil {
		rtc.SetUpdatedAt(*t)
	}
	return rtc
}

// SetRecruitmentID sets the "recruitment_id" field.
func (rtc *RecruitmentTagCreate) SetRecruitmentID(s string) *RecruitmentTagCreate {
	rtc.mutation.SetRecruitmentID(s)
	return rtc
}

// SetTagID sets the "tag_id" field.
func (rtc *RecruitmentTagCreate) SetTagID(s string) *RecruitmentTagCreate {
	rtc.mutation.SetTagID(s)
	return rtc
}

// SetID sets the "id" field.
func (rtc *RecruitmentTagCreate) SetID(s string) *RecruitmentTagCreate {
	rtc.mutation.SetID(s)
	return rtc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (rtc *RecruitmentTagCreate) SetNillableID(s *string) *RecruitmentTagCreate {
	if s != nil {
		rtc.SetID(*s)
	}
	return rtc
}

// SetRecruitment sets the "recruitment" edge to the Recruitment entity.
func (rtc *RecruitmentTagCreate) SetRecruitment(r *Recruitment) *RecruitmentTagCreate {
	return rtc.SetRecruitmentID(r.ID)
}

// SetTag sets the "tag" edge to the Tag entity.
func (rtc *RecruitmentTagCreate) SetTag(t *Tag) *RecruitmentTagCreate {
	return rtc.SetTagID(t.ID)
}

// Mutation returns the RecruitmentTagMutation object of the builder.
func (rtc *RecruitmentTagCreate) Mutation() *RecruitmentTagMutation {
	return rtc.mutation
}

// Save creates the RecruitmentTag in the database.
func (rtc *RecruitmentTagCreate) Save(ctx context.Context) (*RecruitmentTag, error) {
	var (
		err  error
		node *RecruitmentTag
	)
	rtc.defaults()
	if len(rtc.hooks) == 0 {
		if err = rtc.check(); err != nil {
			return nil, err
		}
		node, err = rtc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*RecruitmentTagMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = rtc.check(); err != nil {
				return nil, err
			}
			rtc.mutation = mutation
			if node, err = rtc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(rtc.hooks) - 1; i >= 0; i-- {
			if rtc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = rtc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rtc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (rtc *RecruitmentTagCreate) SaveX(ctx context.Context) *RecruitmentTag {
	v, err := rtc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rtc *RecruitmentTagCreate) Exec(ctx context.Context) error {
	_, err := rtc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rtc *RecruitmentTagCreate) ExecX(ctx context.Context) {
	if err := rtc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (rtc *RecruitmentTagCreate) defaults() {
	if _, ok := rtc.mutation.CreatedAt(); !ok {
		v := recruitmenttag.DefaultCreatedAt()
		rtc.mutation.SetCreatedAt(v)
	}
	if _, ok := rtc.mutation.UpdatedAt(); !ok {
		v := recruitmenttag.DefaultUpdatedAt()
		rtc.mutation.SetUpdatedAt(v)
	}
	if _, ok := rtc.mutation.ID(); !ok {
		v := recruitmenttag.DefaultID()
		rtc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rtc *RecruitmentTagCreate) check() error {
	if _, ok := rtc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "created_at"`)}
	}
	if _, ok := rtc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "updated_at"`)}
	}
	if _, ok := rtc.mutation.RecruitmentID(); !ok {
		return &ValidationError{Name: "recruitment_id", err: errors.New(`ent: missing required field "recruitment_id"`)}
	}
	if _, ok := rtc.mutation.TagID(); !ok {
		return &ValidationError{Name: "tag_id", err: errors.New(`ent: missing required field "tag_id"`)}
	}
	if v, ok := rtc.mutation.ID(); ok {
		if err := recruitmenttag.IDValidator(v); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`ent: validator failed for field "id": %w`, err)}
		}
	}
	if _, ok := rtc.mutation.RecruitmentID(); !ok {
		return &ValidationError{Name: "recruitment", err: errors.New("ent: missing required edge \"recruitment\"")}
	}
	if _, ok := rtc.mutation.TagID(); !ok {
		return &ValidationError{Name: "tag", err: errors.New("ent: missing required edge \"tag\"")}
	}
	return nil
}

func (rtc *RecruitmentTagCreate) sqlSave(ctx context.Context) (*RecruitmentTag, error) {
	_node, _spec := rtc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rtc.driver, _spec); err != nil {
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

func (rtc *RecruitmentTagCreate) createSpec() (*RecruitmentTag, *sqlgraph.CreateSpec) {
	var (
		_node = &RecruitmentTag{config: rtc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: recruitmenttag.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeString,
				Column: recruitmenttag.FieldID,
			},
		}
	)
	if id, ok := rtc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := rtc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: recruitmenttag.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := rtc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: recruitmenttag.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if nodes := rtc.mutation.RecruitmentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   recruitmenttag.RecruitmentTable,
			Columns: []string{recruitmenttag.RecruitmentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: recruitment.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.RecruitmentID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := rtc.mutation.TagIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   recruitmenttag.TagTable,
			Columns: []string{recruitmenttag.TagColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeString,
					Column: tag.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.TagID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// RecruitmentTagCreateBulk is the builder for creating many RecruitmentTag entities in bulk.
type RecruitmentTagCreateBulk struct {
	config
	builders []*RecruitmentTagCreate
}

// Save creates the RecruitmentTag entities in the database.
func (rtcb *RecruitmentTagCreateBulk) Save(ctx context.Context) ([]*RecruitmentTag, error) {
	specs := make([]*sqlgraph.CreateSpec, len(rtcb.builders))
	nodes := make([]*RecruitmentTag, len(rtcb.builders))
	mutators := make([]Mutator, len(rtcb.builders))
	for i := range rtcb.builders {
		func(i int, root context.Context) {
			builder := rtcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*RecruitmentTagMutation)
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
					_, err = mutators[i+1].Mutate(root, rtcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rtcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, rtcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rtcb *RecruitmentTagCreateBulk) SaveX(ctx context.Context) []*RecruitmentTag {
	v, err := rtcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rtcb *RecruitmentTagCreateBulk) Exec(ctx context.Context) error {
	_, err := rtcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rtcb *RecruitmentTagCreateBulk) ExecX(ctx context.Context) {
	if err := rtcb.Exec(ctx); err != nil {
		panic(err)
	}
}