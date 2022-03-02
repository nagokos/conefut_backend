// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"log"

	"github.com/nagokos/connefut_backend/ent/migrate"

	"github.com/nagokos/connefut_backend/ent/competition"
	"github.com/nagokos/connefut_backend/ent/prefecture"
	"github.com/nagokos/connefut_backend/ent/recruitment"
	"github.com/nagokos/connefut_backend/ent/stock"
	"github.com/nagokos/connefut_backend/ent/user"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Competition is the client for interacting with the Competition builders.
	Competition *CompetitionClient
	// Prefecture is the client for interacting with the Prefecture builders.
	Prefecture *PrefectureClient
	// Recruitment is the client for interacting with the Recruitment builders.
	Recruitment *RecruitmentClient
	// Stock is the client for interacting with the Stock builders.
	Stock *StockClient
	// User is the client for interacting with the User builders.
	User *UserClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Competition = NewCompetitionClient(c.config)
	c.Prefecture = NewPrefectureClient(c.config)
	c.Recruitment = NewRecruitmentClient(c.config)
	c.Stock = NewStockClient(c.config)
	c.User = NewUserClient(c.config)
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:         ctx,
		config:      cfg,
		Competition: NewCompetitionClient(cfg),
		Prefecture:  NewPrefectureClient(cfg),
		Recruitment: NewRecruitmentClient(cfg),
		Stock:       NewStockClient(cfg),
		User:        NewUserClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		config:      cfg,
		Competition: NewCompetitionClient(cfg),
		Prefecture:  NewPrefectureClient(cfg),
		Recruitment: NewRecruitmentClient(cfg),
		Stock:       NewStockClient(cfg),
		User:        NewUserClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Competition.
//		Query().
//		Count(ctx)
//
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Competition.Use(hooks...)
	c.Prefecture.Use(hooks...)
	c.Recruitment.Use(hooks...)
	c.Stock.Use(hooks...)
	c.User.Use(hooks...)
}

// CompetitionClient is a client for the Competition schema.
type CompetitionClient struct {
	config
}

// NewCompetitionClient returns a client for the Competition from the given config.
func NewCompetitionClient(c config) *CompetitionClient {
	return &CompetitionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `competition.Hooks(f(g(h())))`.
func (c *CompetitionClient) Use(hooks ...Hook) {
	c.hooks.Competition = append(c.hooks.Competition, hooks...)
}

// Create returns a create builder for Competition.
func (c *CompetitionClient) Create() *CompetitionCreate {
	mutation := newCompetitionMutation(c.config, OpCreate)
	return &CompetitionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Competition entities.
func (c *CompetitionClient) CreateBulk(builders ...*CompetitionCreate) *CompetitionCreateBulk {
	return &CompetitionCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Competition.
func (c *CompetitionClient) Update() *CompetitionUpdate {
	mutation := newCompetitionMutation(c.config, OpUpdate)
	return &CompetitionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *CompetitionClient) UpdateOne(co *Competition) *CompetitionUpdateOne {
	mutation := newCompetitionMutation(c.config, OpUpdateOne, withCompetition(co))
	return &CompetitionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *CompetitionClient) UpdateOneID(id string) *CompetitionUpdateOne {
	mutation := newCompetitionMutation(c.config, OpUpdateOne, withCompetitionID(id))
	return &CompetitionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Competition.
func (c *CompetitionClient) Delete() *CompetitionDelete {
	mutation := newCompetitionMutation(c.config, OpDelete)
	return &CompetitionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *CompetitionClient) DeleteOne(co *Competition) *CompetitionDeleteOne {
	return c.DeleteOneID(co.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *CompetitionClient) DeleteOneID(id string) *CompetitionDeleteOne {
	builder := c.Delete().Where(competition.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &CompetitionDeleteOne{builder}
}

// Query returns a query builder for Competition.
func (c *CompetitionClient) Query() *CompetitionQuery {
	return &CompetitionQuery{
		config: c.config,
	}
}

// Get returns a Competition entity by its id.
func (c *CompetitionClient) Get(ctx context.Context, id string) (*Competition, error) {
	return c.Query().Where(competition.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *CompetitionClient) GetX(ctx context.Context, id string) *Competition {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryRecruitments queries the recruitments edge of a Competition.
func (c *CompetitionClient) QueryRecruitments(co *Competition) *RecruitmentQuery {
	query := &RecruitmentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := co.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(competition.Table, competition.FieldID, id),
			sqlgraph.To(recruitment.Table, recruitment.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, competition.RecruitmentsTable, competition.RecruitmentsColumn),
		)
		fromV = sqlgraph.Neighbors(co.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *CompetitionClient) Hooks() []Hook {
	return c.hooks.Competition
}

// PrefectureClient is a client for the Prefecture schema.
type PrefectureClient struct {
	config
}

// NewPrefectureClient returns a client for the Prefecture from the given config.
func NewPrefectureClient(c config) *PrefectureClient {
	return &PrefectureClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `prefecture.Hooks(f(g(h())))`.
func (c *PrefectureClient) Use(hooks ...Hook) {
	c.hooks.Prefecture = append(c.hooks.Prefecture, hooks...)
}

// Create returns a create builder for Prefecture.
func (c *PrefectureClient) Create() *PrefectureCreate {
	mutation := newPrefectureMutation(c.config, OpCreate)
	return &PrefectureCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Prefecture entities.
func (c *PrefectureClient) CreateBulk(builders ...*PrefectureCreate) *PrefectureCreateBulk {
	return &PrefectureCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Prefecture.
func (c *PrefectureClient) Update() *PrefectureUpdate {
	mutation := newPrefectureMutation(c.config, OpUpdate)
	return &PrefectureUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *PrefectureClient) UpdateOne(pr *Prefecture) *PrefectureUpdateOne {
	mutation := newPrefectureMutation(c.config, OpUpdateOne, withPrefecture(pr))
	return &PrefectureUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *PrefectureClient) UpdateOneID(id string) *PrefectureUpdateOne {
	mutation := newPrefectureMutation(c.config, OpUpdateOne, withPrefectureID(id))
	return &PrefectureUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Prefecture.
func (c *PrefectureClient) Delete() *PrefectureDelete {
	mutation := newPrefectureMutation(c.config, OpDelete)
	return &PrefectureDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *PrefectureClient) DeleteOne(pr *Prefecture) *PrefectureDeleteOne {
	return c.DeleteOneID(pr.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *PrefectureClient) DeleteOneID(id string) *PrefectureDeleteOne {
	builder := c.Delete().Where(prefecture.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &PrefectureDeleteOne{builder}
}

// Query returns a query builder for Prefecture.
func (c *PrefectureClient) Query() *PrefectureQuery {
	return &PrefectureQuery{
		config: c.config,
	}
}

// Get returns a Prefecture entity by its id.
func (c *PrefectureClient) Get(ctx context.Context, id string) (*Prefecture, error) {
	return c.Query().Where(prefecture.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *PrefectureClient) GetX(ctx context.Context, id string) *Prefecture {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryRecruitments queries the recruitments edge of a Prefecture.
func (c *PrefectureClient) QueryRecruitments(pr *Prefecture) *RecruitmentQuery {
	query := &RecruitmentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := pr.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(prefecture.Table, prefecture.FieldID, id),
			sqlgraph.To(recruitment.Table, recruitment.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, prefecture.RecruitmentsTable, prefecture.RecruitmentsColumn),
		)
		fromV = sqlgraph.Neighbors(pr.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *PrefectureClient) Hooks() []Hook {
	return c.hooks.Prefecture
}

// RecruitmentClient is a client for the Recruitment schema.
type RecruitmentClient struct {
	config
}

// NewRecruitmentClient returns a client for the Recruitment from the given config.
func NewRecruitmentClient(c config) *RecruitmentClient {
	return &RecruitmentClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `recruitment.Hooks(f(g(h())))`.
func (c *RecruitmentClient) Use(hooks ...Hook) {
	c.hooks.Recruitment = append(c.hooks.Recruitment, hooks...)
}

// Create returns a create builder for Recruitment.
func (c *RecruitmentClient) Create() *RecruitmentCreate {
	mutation := newRecruitmentMutation(c.config, OpCreate)
	return &RecruitmentCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Recruitment entities.
func (c *RecruitmentClient) CreateBulk(builders ...*RecruitmentCreate) *RecruitmentCreateBulk {
	return &RecruitmentCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Recruitment.
func (c *RecruitmentClient) Update() *RecruitmentUpdate {
	mutation := newRecruitmentMutation(c.config, OpUpdate)
	return &RecruitmentUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *RecruitmentClient) UpdateOne(r *Recruitment) *RecruitmentUpdateOne {
	mutation := newRecruitmentMutation(c.config, OpUpdateOne, withRecruitment(r))
	return &RecruitmentUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *RecruitmentClient) UpdateOneID(id string) *RecruitmentUpdateOne {
	mutation := newRecruitmentMutation(c.config, OpUpdateOne, withRecruitmentID(id))
	return &RecruitmentUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Recruitment.
func (c *RecruitmentClient) Delete() *RecruitmentDelete {
	mutation := newRecruitmentMutation(c.config, OpDelete)
	return &RecruitmentDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *RecruitmentClient) DeleteOne(r *Recruitment) *RecruitmentDeleteOne {
	return c.DeleteOneID(r.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *RecruitmentClient) DeleteOneID(id string) *RecruitmentDeleteOne {
	builder := c.Delete().Where(recruitment.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &RecruitmentDeleteOne{builder}
}

// Query returns a query builder for Recruitment.
func (c *RecruitmentClient) Query() *RecruitmentQuery {
	return &RecruitmentQuery{
		config: c.config,
	}
}

// Get returns a Recruitment entity by its id.
func (c *RecruitmentClient) Get(ctx context.Context, id string) (*Recruitment, error) {
	return c.Query().Where(recruitment.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *RecruitmentClient) GetX(ctx context.Context, id string) *Recruitment {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryStocks queries the stocks edge of a Recruitment.
func (c *RecruitmentClient) QueryStocks(r *Recruitment) *StockQuery {
	query := &StockQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := r.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(recruitment.Table, recruitment.FieldID, id),
			sqlgraph.To(stock.Table, stock.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, recruitment.StocksTable, recruitment.StocksColumn),
		)
		fromV = sqlgraph.Neighbors(r.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryUser queries the user edge of a Recruitment.
func (c *RecruitmentClient) QueryUser(r *Recruitment) *UserQuery {
	query := &UserQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := r.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(recruitment.Table, recruitment.FieldID, id),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, recruitment.UserTable, recruitment.UserColumn),
		)
		fromV = sqlgraph.Neighbors(r.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryPrefecture queries the prefecture edge of a Recruitment.
func (c *RecruitmentClient) QueryPrefecture(r *Recruitment) *PrefectureQuery {
	query := &PrefectureQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := r.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(recruitment.Table, recruitment.FieldID, id),
			sqlgraph.To(prefecture.Table, prefecture.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, recruitment.PrefectureTable, recruitment.PrefectureColumn),
		)
		fromV = sqlgraph.Neighbors(r.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryCompetition queries the competition edge of a Recruitment.
func (c *RecruitmentClient) QueryCompetition(r *Recruitment) *CompetitionQuery {
	query := &CompetitionQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := r.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(recruitment.Table, recruitment.FieldID, id),
			sqlgraph.To(competition.Table, competition.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, recruitment.CompetitionTable, recruitment.CompetitionColumn),
		)
		fromV = sqlgraph.Neighbors(r.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *RecruitmentClient) Hooks() []Hook {
	return c.hooks.Recruitment
}

// StockClient is a client for the Stock schema.
type StockClient struct {
	config
}

// NewStockClient returns a client for the Stock from the given config.
func NewStockClient(c config) *StockClient {
	return &StockClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `stock.Hooks(f(g(h())))`.
func (c *StockClient) Use(hooks ...Hook) {
	c.hooks.Stock = append(c.hooks.Stock, hooks...)
}

// Create returns a create builder for Stock.
func (c *StockClient) Create() *StockCreate {
	mutation := newStockMutation(c.config, OpCreate)
	return &StockCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Stock entities.
func (c *StockClient) CreateBulk(builders ...*StockCreate) *StockCreateBulk {
	return &StockCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Stock.
func (c *StockClient) Update() *StockUpdate {
	mutation := newStockMutation(c.config, OpUpdate)
	return &StockUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *StockClient) UpdateOne(s *Stock) *StockUpdateOne {
	mutation := newStockMutation(c.config, OpUpdateOne, withStock(s))
	return &StockUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *StockClient) UpdateOneID(id string) *StockUpdateOne {
	mutation := newStockMutation(c.config, OpUpdateOne, withStockID(id))
	return &StockUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Stock.
func (c *StockClient) Delete() *StockDelete {
	mutation := newStockMutation(c.config, OpDelete)
	return &StockDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *StockClient) DeleteOne(s *Stock) *StockDeleteOne {
	return c.DeleteOneID(s.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *StockClient) DeleteOneID(id string) *StockDeleteOne {
	builder := c.Delete().Where(stock.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &StockDeleteOne{builder}
}

// Query returns a query builder for Stock.
func (c *StockClient) Query() *StockQuery {
	return &StockQuery{
		config: c.config,
	}
}

// Get returns a Stock entity by its id.
func (c *StockClient) Get(ctx context.Context, id string) (*Stock, error) {
	return c.Query().Where(stock.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *StockClient) GetX(ctx context.Context, id string) *Stock {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryUser queries the user edge of a Stock.
func (c *StockClient) QueryUser(s *Stock) *UserQuery {
	query := &UserQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(stock.Table, stock.FieldID, id),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, stock.UserTable, stock.UserColumn),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryRecruitment queries the recruitment edge of a Stock.
func (c *StockClient) QueryRecruitment(s *Stock) *RecruitmentQuery {
	query := &RecruitmentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := s.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(stock.Table, stock.FieldID, id),
			sqlgraph.To(recruitment.Table, recruitment.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, stock.RecruitmentTable, stock.RecruitmentColumn),
		)
		fromV = sqlgraph.Neighbors(s.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *StockClient) Hooks() []Hook {
	return c.hooks.Stock
}

// UserClient is a client for the User schema.
type UserClient struct {
	config
}

// NewUserClient returns a client for the User from the given config.
func NewUserClient(c config) *UserClient {
	return &UserClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `user.Hooks(f(g(h())))`.
func (c *UserClient) Use(hooks ...Hook) {
	c.hooks.User = append(c.hooks.User, hooks...)
}

// Create returns a create builder for User.
func (c *UserClient) Create() *UserCreate {
	mutation := newUserMutation(c.config, OpCreate)
	return &UserCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of User entities.
func (c *UserClient) CreateBulk(builders ...*UserCreate) *UserCreateBulk {
	return &UserCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for User.
func (c *UserClient) Update() *UserUpdate {
	mutation := newUserMutation(c.config, OpUpdate)
	return &UserUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *UserClient) UpdateOne(u *User) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne, withUser(u))
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *UserClient) UpdateOneID(id string) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne, withUserID(id))
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for User.
func (c *UserClient) Delete() *UserDelete {
	mutation := newUserMutation(c.config, OpDelete)
	return &UserDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a delete builder for the given entity.
func (c *UserClient) DeleteOne(u *User) *UserDeleteOne {
	return c.DeleteOneID(u.ID)
}

// DeleteOneID returns a delete builder for the given id.
func (c *UserClient) DeleteOneID(id string) *UserDeleteOne {
	builder := c.Delete().Where(user.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &UserDeleteOne{builder}
}

// Query returns a query builder for User.
func (c *UserClient) Query() *UserQuery {
	return &UserQuery{
		config: c.config,
	}
}

// Get returns a User entity by its id.
func (c *UserClient) Get(ctx context.Context, id string) (*User, error) {
	return c.Query().Where(user.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *UserClient) GetX(ctx context.Context, id string) *User {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryRecruitments queries the recruitments edge of a User.
func (c *UserClient) QueryRecruitments(u *User) *RecruitmentQuery {
	query := &RecruitmentQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := u.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(user.Table, user.FieldID, id),
			sqlgraph.To(recruitment.Table, recruitment.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, user.RecruitmentsTable, user.RecruitmentsColumn),
		)
		fromV = sqlgraph.Neighbors(u.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryStocks queries the stocks edge of a User.
func (c *UserClient) QueryStocks(u *User) *StockQuery {
	query := &StockQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := u.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(user.Table, user.FieldID, id),
			sqlgraph.To(stock.Table, stock.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, user.StocksTable, user.StocksColumn),
		)
		fromV = sqlgraph.Neighbors(u.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *UserClient) Hooks() []Hook {
	return c.hooks.User
}
