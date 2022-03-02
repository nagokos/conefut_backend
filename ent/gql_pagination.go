// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/nagokos/connefut_backend/ent/competition"
	"github.com/nagokos/connefut_backend/ent/prefecture"
	"github.com/nagokos/connefut_backend/ent/recruitment"
	"github.com/nagokos/connefut_backend/ent/stock"
	"github.com/nagokos/connefut_backend/ent/user"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vmihailenco/msgpack/v5"
)

// OrderDirection defines the directions in which to order a list of items.
type OrderDirection string

const (
	// OrderDirectionAsc specifies an ascending order.
	OrderDirectionAsc OrderDirection = "ASC"
	// OrderDirectionDesc specifies a descending order.
	OrderDirectionDesc OrderDirection = "DESC"
)

// Validate the order direction value.
func (o OrderDirection) Validate() error {
	if o != OrderDirectionAsc && o != OrderDirectionDesc {
		return fmt.Errorf("%s is not a valid OrderDirection", o)
	}
	return nil
}

// String implements fmt.Stringer interface.
func (o OrderDirection) String() string {
	return string(o)
}

// MarshalGQL implements graphql.Marshaler interface.
func (o OrderDirection) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(o.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (o *OrderDirection) UnmarshalGQL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("order direction %T must be a string", val)
	}
	*o = OrderDirection(str)
	return o.Validate()
}

func (o OrderDirection) reverse() OrderDirection {
	if o == OrderDirectionDesc {
		return OrderDirectionAsc
	}
	return OrderDirectionDesc
}

func (o OrderDirection) orderFunc(field string) OrderFunc {
	if o == OrderDirectionDesc {
		return Desc(field)
	}
	return Asc(field)
}

func cursorsToPredicates(direction OrderDirection, after, before *Cursor, field, idField string) []func(s *sql.Selector) {
	var predicates []func(s *sql.Selector)
	if after != nil {
		if after.Value != nil {
			var predicate func([]string, ...interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.CompositeGT
			} else {
				predicate = sql.CompositeLT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.Columns(field, idField),
					after.Value, after.ID,
				))
			})
		} else {
			var predicate func(string, interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.GT
			} else {
				predicate = sql.LT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.C(idField),
					after.ID,
				))
			})
		}
	}
	if before != nil {
		if before.Value != nil {
			var predicate func([]string, ...interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.CompositeLT
			} else {
				predicate = sql.CompositeGT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.Columns(field, idField),
					before.Value, before.ID,
				))
			})
		} else {
			var predicate func(string, interface{}) *sql.Predicate
			if direction == OrderDirectionAsc {
				predicate = sql.LT
			} else {
				predicate = sql.GT
			}
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(predicate(
					s.C(idField),
					before.ID,
				))
			})
		}
	}
	return predicates
}

// PageInfo of a connection type.
type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *Cursor `json:"startCursor"`
	EndCursor       *Cursor `json:"endCursor"`
}

// Cursor of an edge type.
type Cursor struct {
	ID    string `msgpack:"i"`
	Value Value  `msgpack:"v,omitempty"`
}

// MarshalGQL implements graphql.Marshaler interface.
func (c Cursor) MarshalGQL(w io.Writer) {
	quote := []byte{'"'}
	w.Write(quote)
	defer w.Write(quote)
	wc := base64.NewEncoder(base64.RawStdEncoding, w)
	defer wc.Close()
	_ = msgpack.NewEncoder(wc).Encode(c)
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (c *Cursor) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("%T is not a string", v)
	}
	if err := msgpack.NewDecoder(
		base64.NewDecoder(
			base64.RawStdEncoding,
			strings.NewReader(s),
		),
	).Decode(c); err != nil {
		return fmt.Errorf("cannot decode cursor: %w", err)
	}
	return nil
}

const errInvalidPagination = "INVALID_PAGINATION"

func validateFirstLast(first, last *int) (err *gqlerror.Error) {
	switch {
	case first != nil && last != nil:
		err = &gqlerror.Error{
			Message: "Passing both `first` and `last` to paginate a connection is not supported.",
		}
	case first != nil && *first < 0:
		err = &gqlerror.Error{
			Message: "`first` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	case last != nil && *last < 0:
		err = &gqlerror.Error{
			Message: "`last` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	}
	return err
}

func getCollectedField(ctx context.Context, path ...string) *graphql.CollectedField {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil
	}
	oc := graphql.GetOperationContext(ctx)
	field := fc.Field

walk:
	for _, name := range path {
		for _, f := range graphql.CollectFields(oc, field.Selections, nil) {
			if f.Name == name {
				field = f
				continue walk
			}
		}
		return nil
	}
	return &field
}

func hasCollectedField(ctx context.Context, path ...string) bool {
	if graphql.GetFieldContext(ctx) == nil {
		return true
	}
	return getCollectedField(ctx, path...) != nil
}

const (
	edgesField      = "edges"
	nodeField       = "node"
	pageInfoField   = "pageInfo"
	totalCountField = "totalCount"
)

// CompetitionEdge is the edge representation of Competition.
type CompetitionEdge struct {
	Node   *Competition `json:"node"`
	Cursor Cursor       `json:"cursor"`
}

// CompetitionConnection is the connection containing edges to Competition.
type CompetitionConnection struct {
	Edges      []*CompetitionEdge `json:"edges"`
	PageInfo   PageInfo           `json:"pageInfo"`
	TotalCount int                `json:"totalCount"`
}

// CompetitionPaginateOption enables pagination customization.
type CompetitionPaginateOption func(*competitionPager) error

// WithCompetitionOrder configures pagination ordering.
func WithCompetitionOrder(order *CompetitionOrder) CompetitionPaginateOption {
	if order == nil {
		order = DefaultCompetitionOrder
	}
	o := *order
	return func(pager *competitionPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultCompetitionOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithCompetitionFilter configures pagination filter.
func WithCompetitionFilter(filter func(*CompetitionQuery) (*CompetitionQuery, error)) CompetitionPaginateOption {
	return func(pager *competitionPager) error {
		if filter == nil {
			return errors.New("CompetitionQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type competitionPager struct {
	order  *CompetitionOrder
	filter func(*CompetitionQuery) (*CompetitionQuery, error)
}

func newCompetitionPager(opts []CompetitionPaginateOption) (*competitionPager, error) {
	pager := &competitionPager{}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultCompetitionOrder
	}
	return pager, nil
}

func (p *competitionPager) applyFilter(query *CompetitionQuery) (*CompetitionQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *competitionPager) toCursor(c *Competition) Cursor {
	return p.order.Field.toCursor(c)
}

func (p *competitionPager) applyCursors(query *CompetitionQuery, after, before *Cursor) *CompetitionQuery {
	for _, predicate := range cursorsToPredicates(
		p.order.Direction, after, before,
		p.order.Field.field, DefaultCompetitionOrder.Field.field,
	) {
		query = query.Where(predicate)
	}
	return query
}

func (p *competitionPager) applyOrder(query *CompetitionQuery, reverse bool) *CompetitionQuery {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	query = query.Order(direction.orderFunc(p.order.Field.field))
	if p.order.Field != DefaultCompetitionOrder.Field {
		query = query.Order(direction.orderFunc(DefaultCompetitionOrder.Field.field))
	}
	return query
}

// Paginate executes the query and returns a relay based cursor connection to Competition.
func (c *CompetitionQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...CompetitionPaginateOption,
) (*CompetitionConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newCompetitionPager(opts)
	if err != nil {
		return nil, err
	}

	if c, err = pager.applyFilter(c); err != nil {
		return nil, err
	}

	conn := &CompetitionConnection{Edges: []*CompetitionEdge{}}
	if !hasCollectedField(ctx, edgesField) || first != nil && *first == 0 || last != nil && *last == 0 {
		if hasCollectedField(ctx, totalCountField) ||
			hasCollectedField(ctx, pageInfoField) {
			count, err := c.Count(ctx)
			if err != nil {
				return nil, err
			}
			conn.TotalCount = count
			conn.PageInfo.HasNextPage = first != nil && count > 0
			conn.PageInfo.HasPreviousPage = last != nil && count > 0
		}
		return conn, nil
	}

	if (after != nil || first != nil || before != nil || last != nil) && hasCollectedField(ctx, totalCountField) {
		count, err := c.Clone().Count(ctx)
		if err != nil {
			return nil, err
		}
		conn.TotalCount = count
	}

	c = pager.applyCursors(c, after, before)
	c = pager.applyOrder(c, last != nil)
	var limit int
	if first != nil {
		limit = *first + 1
	} else if last != nil {
		limit = *last + 1
	}
	if limit > 0 {
		c = c.Limit(limit)
	}

	if field := getCollectedField(ctx, edgesField, nodeField); field != nil {
		c = c.collectField(graphql.GetOperationContext(ctx), *field)
	}

	nodes, err := c.All(ctx)
	if err != nil || len(nodes) == 0 {
		return conn, err
	}

	if len(nodes) == limit {
		conn.PageInfo.HasNextPage = first != nil
		conn.PageInfo.HasPreviousPage = last != nil
		nodes = nodes[:len(nodes)-1]
	}

	var nodeAt func(int) *Competition
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *Competition {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *Competition {
			return nodes[i]
		}
	}

	conn.Edges = make([]*CompetitionEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		conn.Edges[i] = &CompetitionEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}

	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor
	if conn.TotalCount == 0 {
		conn.TotalCount = len(nodes)
	}

	return conn, nil
}

// CompetitionOrderField defines the ordering field of Competition.
type CompetitionOrderField struct {
	field    string
	toCursor func(*Competition) Cursor
}

// CompetitionOrder defines the ordering of Competition.
type CompetitionOrder struct {
	Direction OrderDirection         `json:"direction"`
	Field     *CompetitionOrderField `json:"field"`
}

// DefaultCompetitionOrder is the default ordering of Competition.
var DefaultCompetitionOrder = &CompetitionOrder{
	Direction: OrderDirectionAsc,
	Field: &CompetitionOrderField{
		field: competition.FieldID,
		toCursor: func(c *Competition) Cursor {
			return Cursor{ID: c.ID}
		},
	},
}

// ToEdge converts Competition into CompetitionEdge.
func (c *Competition) ToEdge(order *CompetitionOrder) *CompetitionEdge {
	if order == nil {
		order = DefaultCompetitionOrder
	}
	return &CompetitionEdge{
		Node:   c,
		Cursor: order.Field.toCursor(c),
	}
}

// PrefectureEdge is the edge representation of Prefecture.
type PrefectureEdge struct {
	Node   *Prefecture `json:"node"`
	Cursor Cursor      `json:"cursor"`
}

// PrefectureConnection is the connection containing edges to Prefecture.
type PrefectureConnection struct {
	Edges      []*PrefectureEdge `json:"edges"`
	PageInfo   PageInfo          `json:"pageInfo"`
	TotalCount int               `json:"totalCount"`
}

// PrefecturePaginateOption enables pagination customization.
type PrefecturePaginateOption func(*prefecturePager) error

// WithPrefectureOrder configures pagination ordering.
func WithPrefectureOrder(order *PrefectureOrder) PrefecturePaginateOption {
	if order == nil {
		order = DefaultPrefectureOrder
	}
	o := *order
	return func(pager *prefecturePager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultPrefectureOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithPrefectureFilter configures pagination filter.
func WithPrefectureFilter(filter func(*PrefectureQuery) (*PrefectureQuery, error)) PrefecturePaginateOption {
	return func(pager *prefecturePager) error {
		if filter == nil {
			return errors.New("PrefectureQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type prefecturePager struct {
	order  *PrefectureOrder
	filter func(*PrefectureQuery) (*PrefectureQuery, error)
}

func newPrefecturePager(opts []PrefecturePaginateOption) (*prefecturePager, error) {
	pager := &prefecturePager{}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultPrefectureOrder
	}
	return pager, nil
}

func (p *prefecturePager) applyFilter(query *PrefectureQuery) (*PrefectureQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *prefecturePager) toCursor(pr *Prefecture) Cursor {
	return p.order.Field.toCursor(pr)
}

func (p *prefecturePager) applyCursors(query *PrefectureQuery, after, before *Cursor) *PrefectureQuery {
	for _, predicate := range cursorsToPredicates(
		p.order.Direction, after, before,
		p.order.Field.field, DefaultPrefectureOrder.Field.field,
	) {
		query = query.Where(predicate)
	}
	return query
}

func (p *prefecturePager) applyOrder(query *PrefectureQuery, reverse bool) *PrefectureQuery {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	query = query.Order(direction.orderFunc(p.order.Field.field))
	if p.order.Field != DefaultPrefectureOrder.Field {
		query = query.Order(direction.orderFunc(DefaultPrefectureOrder.Field.field))
	}
	return query
}

// Paginate executes the query and returns a relay based cursor connection to Prefecture.
func (pr *PrefectureQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...PrefecturePaginateOption,
) (*PrefectureConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newPrefecturePager(opts)
	if err != nil {
		return nil, err
	}

	if pr, err = pager.applyFilter(pr); err != nil {
		return nil, err
	}

	conn := &PrefectureConnection{Edges: []*PrefectureEdge{}}
	if !hasCollectedField(ctx, edgesField) || first != nil && *first == 0 || last != nil && *last == 0 {
		if hasCollectedField(ctx, totalCountField) ||
			hasCollectedField(ctx, pageInfoField) {
			count, err := pr.Count(ctx)
			if err != nil {
				return nil, err
			}
			conn.TotalCount = count
			conn.PageInfo.HasNextPage = first != nil && count > 0
			conn.PageInfo.HasPreviousPage = last != nil && count > 0
		}
		return conn, nil
	}

	if (after != nil || first != nil || before != nil || last != nil) && hasCollectedField(ctx, totalCountField) {
		count, err := pr.Clone().Count(ctx)
		if err != nil {
			return nil, err
		}
		conn.TotalCount = count
	}

	pr = pager.applyCursors(pr, after, before)
	pr = pager.applyOrder(pr, last != nil)
	var limit int
	if first != nil {
		limit = *first + 1
	} else if last != nil {
		limit = *last + 1
	}
	if limit > 0 {
		pr = pr.Limit(limit)
	}

	if field := getCollectedField(ctx, edgesField, nodeField); field != nil {
		pr = pr.collectField(graphql.GetOperationContext(ctx), *field)
	}

	nodes, err := pr.All(ctx)
	if err != nil || len(nodes) == 0 {
		return conn, err
	}

	if len(nodes) == limit {
		conn.PageInfo.HasNextPage = first != nil
		conn.PageInfo.HasPreviousPage = last != nil
		nodes = nodes[:len(nodes)-1]
	}

	var nodeAt func(int) *Prefecture
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *Prefecture {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *Prefecture {
			return nodes[i]
		}
	}

	conn.Edges = make([]*PrefectureEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		conn.Edges[i] = &PrefectureEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}

	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor
	if conn.TotalCount == 0 {
		conn.TotalCount = len(nodes)
	}

	return conn, nil
}

// PrefectureOrderField defines the ordering field of Prefecture.
type PrefectureOrderField struct {
	field    string
	toCursor func(*Prefecture) Cursor
}

// PrefectureOrder defines the ordering of Prefecture.
type PrefectureOrder struct {
	Direction OrderDirection        `json:"direction"`
	Field     *PrefectureOrderField `json:"field"`
}

// DefaultPrefectureOrder is the default ordering of Prefecture.
var DefaultPrefectureOrder = &PrefectureOrder{
	Direction: OrderDirectionAsc,
	Field: &PrefectureOrderField{
		field: prefecture.FieldID,
		toCursor: func(pr *Prefecture) Cursor {
			return Cursor{ID: pr.ID}
		},
	},
}

// ToEdge converts Prefecture into PrefectureEdge.
func (pr *Prefecture) ToEdge(order *PrefectureOrder) *PrefectureEdge {
	if order == nil {
		order = DefaultPrefectureOrder
	}
	return &PrefectureEdge{
		Node:   pr,
		Cursor: order.Field.toCursor(pr),
	}
}

// RecruitmentEdge is the edge representation of Recruitment.
type RecruitmentEdge struct {
	Node   *Recruitment `json:"node"`
	Cursor Cursor       `json:"cursor"`
}

// RecruitmentConnection is the connection containing edges to Recruitment.
type RecruitmentConnection struct {
	Edges      []*RecruitmentEdge `json:"edges"`
	PageInfo   PageInfo           `json:"pageInfo"`
	TotalCount int                `json:"totalCount"`
}

// RecruitmentPaginateOption enables pagination customization.
type RecruitmentPaginateOption func(*recruitmentPager) error

// WithRecruitmentOrder configures pagination ordering.
func WithRecruitmentOrder(order *RecruitmentOrder) RecruitmentPaginateOption {
	if order == nil {
		order = DefaultRecruitmentOrder
	}
	o := *order
	return func(pager *recruitmentPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultRecruitmentOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithRecruitmentFilter configures pagination filter.
func WithRecruitmentFilter(filter func(*RecruitmentQuery) (*RecruitmentQuery, error)) RecruitmentPaginateOption {
	return func(pager *recruitmentPager) error {
		if filter == nil {
			return errors.New("RecruitmentQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type recruitmentPager struct {
	order  *RecruitmentOrder
	filter func(*RecruitmentQuery) (*RecruitmentQuery, error)
}

func newRecruitmentPager(opts []RecruitmentPaginateOption) (*recruitmentPager, error) {
	pager := &recruitmentPager{}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultRecruitmentOrder
	}
	return pager, nil
}

func (p *recruitmentPager) applyFilter(query *RecruitmentQuery) (*RecruitmentQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *recruitmentPager) toCursor(r *Recruitment) Cursor {
	return p.order.Field.toCursor(r)
}

func (p *recruitmentPager) applyCursors(query *RecruitmentQuery, after, before *Cursor) *RecruitmentQuery {
	for _, predicate := range cursorsToPredicates(
		p.order.Direction, after, before,
		p.order.Field.field, DefaultRecruitmentOrder.Field.field,
	) {
		query = query.Where(predicate)
	}
	return query
}

func (p *recruitmentPager) applyOrder(query *RecruitmentQuery, reverse bool) *RecruitmentQuery {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	query = query.Order(direction.orderFunc(p.order.Field.field))
	if p.order.Field != DefaultRecruitmentOrder.Field {
		query = query.Order(direction.orderFunc(DefaultRecruitmentOrder.Field.field))
	}
	return query
}

// Paginate executes the query and returns a relay based cursor connection to Recruitment.
func (r *RecruitmentQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...RecruitmentPaginateOption,
) (*RecruitmentConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newRecruitmentPager(opts)
	if err != nil {
		return nil, err
	}

	if r, err = pager.applyFilter(r); err != nil {
		return nil, err
	}

	conn := &RecruitmentConnection{Edges: []*RecruitmentEdge{}}
	if !hasCollectedField(ctx, edgesField) || first != nil && *first == 0 || last != nil && *last == 0 {
		if hasCollectedField(ctx, totalCountField) ||
			hasCollectedField(ctx, pageInfoField) {
			count, err := r.Count(ctx)
			if err != nil {
				return nil, err
			}
			conn.TotalCount = count
			conn.PageInfo.HasNextPage = first != nil && count > 0
			conn.PageInfo.HasPreviousPage = last != nil && count > 0
		}
		return conn, nil
	}

	if (after != nil || first != nil || before != nil || last != nil) && hasCollectedField(ctx, totalCountField) {
		count, err := r.Clone().Count(ctx)
		if err != nil {
			return nil, err
		}
		conn.TotalCount = count
	}

	r = pager.applyCursors(r, after, before)
	r = pager.applyOrder(r, last != nil)
	var limit int
	if first != nil {
		limit = *first + 1
	} else if last != nil {
		limit = *last + 1
	}
	if limit > 0 {
		r = r.Limit(limit)
	}

	if field := getCollectedField(ctx, edgesField, nodeField); field != nil {
		r = r.collectField(graphql.GetOperationContext(ctx), *field)
	}

	nodes, err := r.All(ctx)
	if err != nil || len(nodes) == 0 {
		return conn, err
	}

	if len(nodes) == limit {
		conn.PageInfo.HasNextPage = first != nil
		conn.PageInfo.HasPreviousPage = last != nil
		nodes = nodes[:len(nodes)-1]
	}

	var nodeAt func(int) *Recruitment
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *Recruitment {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *Recruitment {
			return nodes[i]
		}
	}

	conn.Edges = make([]*RecruitmentEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		conn.Edges[i] = &RecruitmentEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}

	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor
	if conn.TotalCount == 0 {
		conn.TotalCount = len(nodes)
	}

	return conn, nil
}

// RecruitmentOrderField defines the ordering field of Recruitment.
type RecruitmentOrderField struct {
	field    string
	toCursor func(*Recruitment) Cursor
}

// RecruitmentOrder defines the ordering of Recruitment.
type RecruitmentOrder struct {
	Direction OrderDirection         `json:"direction"`
	Field     *RecruitmentOrderField `json:"field"`
}

// DefaultRecruitmentOrder is the default ordering of Recruitment.
var DefaultRecruitmentOrder = &RecruitmentOrder{
	Direction: OrderDirectionAsc,
	Field: &RecruitmentOrderField{
		field: recruitment.FieldID,
		toCursor: func(r *Recruitment) Cursor {
			return Cursor{ID: r.ID}
		},
	},
}

// ToEdge converts Recruitment into RecruitmentEdge.
func (r *Recruitment) ToEdge(order *RecruitmentOrder) *RecruitmentEdge {
	if order == nil {
		order = DefaultRecruitmentOrder
	}
	return &RecruitmentEdge{
		Node:   r,
		Cursor: order.Field.toCursor(r),
	}
}

// StockEdge is the edge representation of Stock.
type StockEdge struct {
	Node   *Stock `json:"node"`
	Cursor Cursor `json:"cursor"`
}

// StockConnection is the connection containing edges to Stock.
type StockConnection struct {
	Edges      []*StockEdge `json:"edges"`
	PageInfo   PageInfo     `json:"pageInfo"`
	TotalCount int          `json:"totalCount"`
}

// StockPaginateOption enables pagination customization.
type StockPaginateOption func(*stockPager) error

// WithStockOrder configures pagination ordering.
func WithStockOrder(order *StockOrder) StockPaginateOption {
	if order == nil {
		order = DefaultStockOrder
	}
	o := *order
	return func(pager *stockPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultStockOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithStockFilter configures pagination filter.
func WithStockFilter(filter func(*StockQuery) (*StockQuery, error)) StockPaginateOption {
	return func(pager *stockPager) error {
		if filter == nil {
			return errors.New("StockQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type stockPager struct {
	order  *StockOrder
	filter func(*StockQuery) (*StockQuery, error)
}

func newStockPager(opts []StockPaginateOption) (*stockPager, error) {
	pager := &stockPager{}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultStockOrder
	}
	return pager, nil
}

func (p *stockPager) applyFilter(query *StockQuery) (*StockQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *stockPager) toCursor(s *Stock) Cursor {
	return p.order.Field.toCursor(s)
}

func (p *stockPager) applyCursors(query *StockQuery, after, before *Cursor) *StockQuery {
	for _, predicate := range cursorsToPredicates(
		p.order.Direction, after, before,
		p.order.Field.field, DefaultStockOrder.Field.field,
	) {
		query = query.Where(predicate)
	}
	return query
}

func (p *stockPager) applyOrder(query *StockQuery, reverse bool) *StockQuery {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	query = query.Order(direction.orderFunc(p.order.Field.field))
	if p.order.Field != DefaultStockOrder.Field {
		query = query.Order(direction.orderFunc(DefaultStockOrder.Field.field))
	}
	return query
}

// Paginate executes the query and returns a relay based cursor connection to Stock.
func (s *StockQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...StockPaginateOption,
) (*StockConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newStockPager(opts)
	if err != nil {
		return nil, err
	}

	if s, err = pager.applyFilter(s); err != nil {
		return nil, err
	}

	conn := &StockConnection{Edges: []*StockEdge{}}
	if !hasCollectedField(ctx, edgesField) || first != nil && *first == 0 || last != nil && *last == 0 {
		if hasCollectedField(ctx, totalCountField) ||
			hasCollectedField(ctx, pageInfoField) {
			count, err := s.Count(ctx)
			if err != nil {
				return nil, err
			}
			conn.TotalCount = count
			conn.PageInfo.HasNextPage = first != nil && count > 0
			conn.PageInfo.HasPreviousPage = last != nil && count > 0
		}
		return conn, nil
	}

	if (after != nil || first != nil || before != nil || last != nil) && hasCollectedField(ctx, totalCountField) {
		count, err := s.Clone().Count(ctx)
		if err != nil {
			return nil, err
		}
		conn.TotalCount = count
	}

	s = pager.applyCursors(s, after, before)
	s = pager.applyOrder(s, last != nil)
	var limit int
	if first != nil {
		limit = *first + 1
	} else if last != nil {
		limit = *last + 1
	}
	if limit > 0 {
		s = s.Limit(limit)
	}

	if field := getCollectedField(ctx, edgesField, nodeField); field != nil {
		s = s.collectField(graphql.GetOperationContext(ctx), *field)
	}

	nodes, err := s.All(ctx)
	if err != nil || len(nodes) == 0 {
		return conn, err
	}

	if len(nodes) == limit {
		conn.PageInfo.HasNextPage = first != nil
		conn.PageInfo.HasPreviousPage = last != nil
		nodes = nodes[:len(nodes)-1]
	}

	var nodeAt func(int) *Stock
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *Stock {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *Stock {
			return nodes[i]
		}
	}

	conn.Edges = make([]*StockEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		conn.Edges[i] = &StockEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}

	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor
	if conn.TotalCount == 0 {
		conn.TotalCount = len(nodes)
	}

	return conn, nil
}

// StockOrderField defines the ordering field of Stock.
type StockOrderField struct {
	field    string
	toCursor func(*Stock) Cursor
}

// StockOrder defines the ordering of Stock.
type StockOrder struct {
	Direction OrderDirection   `json:"direction"`
	Field     *StockOrderField `json:"field"`
}

// DefaultStockOrder is the default ordering of Stock.
var DefaultStockOrder = &StockOrder{
	Direction: OrderDirectionAsc,
	Field: &StockOrderField{
		field: stock.FieldID,
		toCursor: func(s *Stock) Cursor {
			return Cursor{ID: s.ID}
		},
	},
}

// ToEdge converts Stock into StockEdge.
func (s *Stock) ToEdge(order *StockOrder) *StockEdge {
	if order == nil {
		order = DefaultStockOrder
	}
	return &StockEdge{
		Node:   s,
		Cursor: order.Field.toCursor(s),
	}
}

// UserEdge is the edge representation of User.
type UserEdge struct {
	Node   *User  `json:"node"`
	Cursor Cursor `json:"cursor"`
}

// UserConnection is the connection containing edges to User.
type UserConnection struct {
	Edges      []*UserEdge `json:"edges"`
	PageInfo   PageInfo    `json:"pageInfo"`
	TotalCount int         `json:"totalCount"`
}

// UserPaginateOption enables pagination customization.
type UserPaginateOption func(*userPager) error

// WithUserOrder configures pagination ordering.
func WithUserOrder(order *UserOrder) UserPaginateOption {
	if order == nil {
		order = DefaultUserOrder
	}
	o := *order
	return func(pager *userPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultUserOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithUserFilter configures pagination filter.
func WithUserFilter(filter func(*UserQuery) (*UserQuery, error)) UserPaginateOption {
	return func(pager *userPager) error {
		if filter == nil {
			return errors.New("UserQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type userPager struct {
	order  *UserOrder
	filter func(*UserQuery) (*UserQuery, error)
}

func newUserPager(opts []UserPaginateOption) (*userPager, error) {
	pager := &userPager{}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultUserOrder
	}
	return pager, nil
}

func (p *userPager) applyFilter(query *UserQuery) (*UserQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *userPager) toCursor(u *User) Cursor {
	return p.order.Field.toCursor(u)
}

func (p *userPager) applyCursors(query *UserQuery, after, before *Cursor) *UserQuery {
	for _, predicate := range cursorsToPredicates(
		p.order.Direction, after, before,
		p.order.Field.field, DefaultUserOrder.Field.field,
	) {
		query = query.Where(predicate)
	}
	return query
}

func (p *userPager) applyOrder(query *UserQuery, reverse bool) *UserQuery {
	direction := p.order.Direction
	if reverse {
		direction = direction.reverse()
	}
	query = query.Order(direction.orderFunc(p.order.Field.field))
	if p.order.Field != DefaultUserOrder.Field {
		query = query.Order(direction.orderFunc(DefaultUserOrder.Field.field))
	}
	return query
}

// Paginate executes the query and returns a relay based cursor connection to User.
func (u *UserQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...UserPaginateOption,
) (*UserConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newUserPager(opts)
	if err != nil {
		return nil, err
	}

	if u, err = pager.applyFilter(u); err != nil {
		return nil, err
	}

	conn := &UserConnection{Edges: []*UserEdge{}}
	if !hasCollectedField(ctx, edgesField) || first != nil && *first == 0 || last != nil && *last == 0 {
		if hasCollectedField(ctx, totalCountField) ||
			hasCollectedField(ctx, pageInfoField) {
			count, err := u.Count(ctx)
			if err != nil {
				return nil, err
			}
			conn.TotalCount = count
			conn.PageInfo.HasNextPage = first != nil && count > 0
			conn.PageInfo.HasPreviousPage = last != nil && count > 0
		}
		return conn, nil
	}

	if (after != nil || first != nil || before != nil || last != nil) && hasCollectedField(ctx, totalCountField) {
		count, err := u.Clone().Count(ctx)
		if err != nil {
			return nil, err
		}
		conn.TotalCount = count
	}

	u = pager.applyCursors(u, after, before)
	u = pager.applyOrder(u, last != nil)
	var limit int
	if first != nil {
		limit = *first + 1
	} else if last != nil {
		limit = *last + 1
	}
	if limit > 0 {
		u = u.Limit(limit)
	}

	if field := getCollectedField(ctx, edgesField, nodeField); field != nil {
		u = u.collectField(graphql.GetOperationContext(ctx), *field)
	}

	nodes, err := u.All(ctx)
	if err != nil || len(nodes) == 0 {
		return conn, err
	}

	if len(nodes) == limit {
		conn.PageInfo.HasNextPage = first != nil
		conn.PageInfo.HasPreviousPage = last != nil
		nodes = nodes[:len(nodes)-1]
	}

	var nodeAt func(int) *User
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *User {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *User {
			return nodes[i]
		}
	}

	conn.Edges = make([]*UserEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		conn.Edges[i] = &UserEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}

	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor
	if conn.TotalCount == 0 {
		conn.TotalCount = len(nodes)
	}

	return conn, nil
}

// UserOrderField defines the ordering field of User.
type UserOrderField struct {
	field    string
	toCursor func(*User) Cursor
}

// UserOrder defines the ordering of User.
type UserOrder struct {
	Direction OrderDirection  `json:"direction"`
	Field     *UserOrderField `json:"field"`
}

// DefaultUserOrder is the default ordering of User.
var DefaultUserOrder = &UserOrder{
	Direction: OrderDirectionAsc,
	Field: &UserOrderField{
		field: user.FieldID,
		toCursor: func(u *User) Cursor {
			return Cursor{ID: u.ID}
		},
	},
}

// ToEdge converts User into UserEdge.
func (u *User) ToEdge(order *UserOrder) *UserEdge {
	if order == nil {
		order = DefaultUserOrder
	}
	return &UserEdge{
		Node:   u,
		Cursor: order.Field.toCursor(u),
	}
}
