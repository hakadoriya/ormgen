package ormopt

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"strings"
)

// Deprecated: Use ormcommon.QueryerContext instead.
type QueryerContext interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// NOTE: noopLogger needs var declare for default use.
//
//nolint:gochecknoglobals
var noopLogger = slog.New(slog.NewJSONHandler(io.Discard, nil))

// Deprecated: Use ormcommon.LoggerFromContext instead.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return noopLogger
	}
	if logger, ok := ctx.Value((*slog.Logger)(nil)).(*slog.Logger); ok {
		return logger
	}
	return noopLogger
}

// Deprecated: Use ormcommon.LoggerWithContext instead.
func LoggerWithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, (*slog.Logger)(nil), logger)
}

type (
	ResultOption interface {
		ApplyResultOption(c *QueryConfig)
		QueryOption
	}
	QueryOption interface {
		ApplyQueryOption(c *QueryConfig)
	}
)

type PlaceholderGenerator func(placeholderStartAt int) string

type QueryConfig struct {
	placeholderGenerator PlaceholderGenerator
	queryPrefix          string
	where                *where
	orderBy              *orderBy
	limit                int
}

func WithPlaceholderGenerator(f PlaceholderGenerator) ResultOption {
	return &withListOptionPlaceholderGenerator{
		placeholderGenerator: f,
	}
}

type withListOptionPlaceholderGenerator struct {
	placeholderGenerator PlaceholderGenerator
}

func (o *withListOptionPlaceholderGenerator) ApplyResultOption(c *QueryConfig) {
	c.placeholderGenerator = o.placeholderGenerator
}

func (o *withListOptionPlaceholderGenerator) ApplyQueryOption(c *QueryConfig) {
	o.ApplyResultOption(c)
}

func (c *QueryConfig) ToSQL(query string, placeholderStartAt int) (string, []interface{}) {
	var q string
	var args []interface{}
	if c.queryPrefix != "" {
		q += c.queryPrefix + " "
	}
	q += query
	if c.where != nil {
		q += " WHERE " + c.where.Root.toSQL(&placeholderStartAt, &args, c.placeholderGenerator)
	}
	if c.orderBy != nil {
		q += " ORDER BY " + strings.Join(c.orderBy.columns, ", ")
		if c.orderBy.desc {
			q += " DESC"
		}
	}
	if c.limit > 0 {
		q += " LIMIT " + c.placeholderGenerator(placeholderStartAt)
		placeholderStartAt++
		args = append(args, c.limit)
	}
	return q, args
}

// ===============================================================
// query prefix
// ===============================================================

func QueryPrefix(prefix string) ResultOption {
	return &withListOptionQueryPrefix{
		prefix: prefix,
	}
}

var (
	_ ResultOption = (*withListOptionQueryPrefix)(nil)
	_ QueryOption  = (*withListOptionQueryPrefix)(nil)
)

type withListOptionQueryPrefix struct {
	prefix string
}

func (o *withListOptionQueryPrefix) ApplyResultOption(c *QueryConfig) {
	c.queryPrefix += o.prefix
}

func (o *withListOptionQueryPrefix) ApplyQueryOption(c *QueryConfig) {
	o.ApplyResultOption(c)
}

// ===============================================================
// ORDER BY
// ===============================================================

func OrderBy(orderBy ...string) ResultOption {
	return &withListOptionOrderBy{
		columns: orderBy,
		desc:    false,
	}
}

func OrderByDesc(orderBy ...string) ResultOption {
	return &withListOptionOrderBy{
		columns: orderBy,
		desc:    true,
	}
}

type orderBy struct {
	columns []string
	desc    bool
}

var (
	_ ResultOption = (*withListOptionOrderBy)(nil)
	_ QueryOption  = (*withListOptionOrderBy)(nil)
)

type withListOptionOrderBy struct {
	columns []string
	desc    bool
}

func (o *withListOptionOrderBy) ApplyResultOption(c *QueryConfig) {
	c.orderBy = &orderBy{
		columns: o.columns,
		desc:    o.desc,
	}
}

func (o *withListOptionOrderBy) ApplyQueryOption(c *QueryConfig) {
	o.ApplyResultOption(c)
}

// ===============================================================
// LIMIT
// ===============================================================

func Limit(limit int) ResultOption {
	return &withListOptionLimit{
		limit: limit,
	}
}

var (
	_ ResultOption = (*withListOptionLimit)(nil)
	_ QueryOption  = (*withListOptionLimit)(nil)
)

type withListOptionLimit struct {
	limit int
}

func (o *withListOptionLimit) ApplyResultOption(c *QueryConfig) {
	c.limit = o.limit
}

func (o *withListOptionLimit) ApplyQueryOption(c *QueryConfig) {
	o.ApplyResultOption(c)
}

// ===============================================================
// WHERE
// ===============================================================

func Where(condition Condition) QueryOption {
	return &withListOptionWhere{
		where: &where{
			Root: condition,
		},
	}
}

type where struct {
	Root Condition
}

type Condition interface {
	toSQL(placeholderStartAt *int, queryArgs *[]interface{}, placeHolderGenerator PlaceholderGenerator) string
}

type conditionCompound struct {
	Op         string // "AND" or "OR"
	Conditions []Condition
}

func (cc *conditionCompound) toSQL(placeholderStartAt *int, queryArgs *[]interface{}, placeHolderGenerator PlaceholderGenerator) string {
	s := "("
	for i, c := range cc.Conditions {
		if i != 0 {
			s += " " + cc.Op + " "
		}
		s += c.toSQL(placeholderStartAt, queryArgs, placeHolderGenerator)
	}
	s += ")"
	return s
}

var _ QueryOption = (*withListOptionWhere)(nil)

type withListOptionWhere struct {
	where *where
}

func (o *withListOptionWhere) ApplyQueryOption(c *QueryConfig) {
	c.where = o.where
}

// ===============================================================
// AND, OR
// ===============================================================

func And(conditions ...Condition) Condition {
	return &conditionCompound{
		Op:         "AND",
		Conditions: conditions,
	}
}

func Or(conditions ...Condition) Condition {
	return &conditionCompound{
		Op:         "OR",
		Conditions: conditions,
	}
}

// ===============================================================
// = (equal)
// ===============================================================

func Equal(column string, value interface{}) Condition {
	return &conditionEqual{
		Column: column,
		Value:  value,
	}
}

type conditionEqual struct {
	Column string
	Value  interface{}
}

func (c *conditionEqual) toSQL(placeholderStartAt *int, queryArgs *[]interface{}, placeHolderGenerator PlaceholderGenerator) string {
	s := c.Column + " = " + placeHolderGenerator(*placeholderStartAt)
	*placeholderStartAt++
	*queryArgs = append(*queryArgs, c.Value)
	return s
}

// ===============================================================
// <> (not equal)
// ===============================================================

func NotEqual(column string, value interface{}) Condition {
	return &conditionNotEqual{
		Column: column,
		Value:  value,
	}
}

type conditionNotEqual struct {
	Column string
	Value  interface{}
}

func (c *conditionNotEqual) toSQL(placeholderStartAt *int, queryArgs *[]interface{}, placeHolderGenerator PlaceholderGenerator) string {
	s := c.Column + " <> " + placeHolderGenerator(*placeholderStartAt)
	*placeholderStartAt++
	*queryArgs = append(*queryArgs, c.Value)
	return s
}

// ===============================================================
// IN
// ===============================================================

func In(column string, values ...interface{}) Condition {
	return &conditionIn{
		Column: column,
		Values: values,
	}
}

type conditionIn struct {
	Column string
	Values []interface{}
}

func (c *conditionIn) toSQL(placeholderStartAt *int, queryArgs *[]interface{}, placeHolderGenerator PlaceholderGenerator) string {
	s := c.Column + " IN ("
	for i, v := range c.Values {
		if i != 0 {
			s += ", "
		}
		s += placeHolderGenerator(*placeholderStartAt)
		*placeholderStartAt++
		*queryArgs = append(*queryArgs, v)
	}
	s += ")"
	return s
}

// ===============================================================
// NOT IN
// ===============================================================

func NotIn(column string, values ...interface{}) Condition {
	return &conditionNotIn{
		Column: column,
		Values: values,
	}
}

type conditionNotIn struct {
	Column string
	Values []interface{}
}

func (c *conditionNotIn) toSQL(placeholderStartAt *int, queryArgs *[]interface{}, placeHolderGenerator PlaceholderGenerator) string {
	s := c.Column + " NOT IN ("
	for i, v := range c.Values {
		if i != 0 {
			s += ", "
		}
		s += placeHolderGenerator(*placeholderStartAt)
		*placeholderStartAt++
		*queryArgs = append(*queryArgs, v)
	}
	s += ")"
	return s
}

// ===============================================================
// > (greater than)
// ===============================================================

func GreaterThan(column string, value interface{}) Condition {
	return &conditionGreaterThan{
		Column: column,
		Value:  value,
	}
}

type conditionGreaterThan struct {
	Column string
	Value  interface{}
}

func (c *conditionGreaterThan) toSQL(placeholderStartAt *int, queryArgs *[]interface{}, placeHolderGenerator PlaceholderGenerator) string {
	s := c.Column + " > " + placeHolderGenerator(*placeholderStartAt)
	*placeholderStartAt++
	*queryArgs = append(*queryArgs, c.Value)
	return s
}

// ===============================================================
// >= (greater than or equal)
// ===============================================================

func GreaterThanOrEqual(column string, value interface{}) Condition {
	return &conditionGreaterThanOrEqual{
		Column: column,
		Value:  value,
	}
}

type conditionGreaterThanOrEqual struct {
	Column string
	Value  interface{}
}

func (c *conditionGreaterThanOrEqual) toSQL(placeholderStartAt *int, queryArgs *[]interface{}, placeHolderGenerator PlaceholderGenerator) string {
	s := c.Column + " >= " + placeHolderGenerator(*placeholderStartAt)
	*placeholderStartAt++
	*queryArgs = append(*queryArgs, c.Value)
	return s
}

// ===============================================================
// < (less than)
// ===============================================================

func LessThan(column string, value interface{}) Condition {
	return &conditionLessThan{
		Column: column,
		Value:  value,
	}
}

type conditionLessThan struct {
	Column string
	Value  interface{}
}

func (c *conditionLessThan) toSQL(placeholderStartAt *int, queryArgs *[]interface{}, placeHolderGenerator PlaceholderGenerator) string {
	s := c.Column + " < " + placeHolderGenerator(*placeholderStartAt)
	*placeholderStartAt++
	*queryArgs = append(*queryArgs, c.Value)
	return s
}

// ===============================================================
// <= (less than or equal)
// ===============================================================

func LessThanOrEqual(column string, value interface{}) Condition {
	return &conditionLessThanOrEqual{
		Column: column,
		Value:  value,
	}
}

type conditionLessThanOrEqual struct {
	Column string
	Value  interface{}
}

func (c *conditionLessThanOrEqual) toSQL(placeholderStartAt *int, queryArgs *[]interface{}, placeHolderGenerator PlaceholderGenerator) string {
	s := c.Column + " <= " + placeHolderGenerator(*placeholderStartAt)
	*placeholderStartAt++
	*queryArgs = append(*queryArgs, c.Value)
	return s
}

// ===============================================================
// LIKE
// ===============================================================

func Like(column string, value interface{}) Condition {
	return &conditionLike{
		Column: column,
		Value:  value,
	}
}

type conditionLike struct {
	Column string
	Value  interface{}
}

func (c *conditionLike) toSQL(placeholderStartAt *int, queryArgs *[]interface{}, placeHolderGenerator PlaceholderGenerator) string {
	s := c.Column + " LIKE " + placeHolderGenerator(*placeholderStartAt)
	*placeholderStartAt++
	*queryArgs = append(*queryArgs, c.Value)
	return s
}

// ===============================================================
// NOT LIKE
// ===============================================================

func NotLike(column string, value interface{}) Condition {
	return &conditionNotLike{
		Column: column,
		Value:  value,
	}
}

type conditionNotLike struct {
	Column string
	Value  interface{}
}

func (c *conditionNotLike) toSQL(placeholderStartAt *int, queryArgs *[]interface{}, placeHolderGenerator PlaceholderGenerator) string {
	s := c.Column + " NOT LIKE " + placeHolderGenerator(*placeholderStartAt)
	*placeholderStartAt++
	*queryArgs = append(*queryArgs, c.Value)
	return s
}

// ===============================================================
// IS NULL
// ===============================================================

func IsNull(column string) Condition {
	return &conditionIsNull{
		Column: column,
	}
}

type conditionIsNull struct {
	Column string
}

func (c *conditionIsNull) toSQL(_ *int, _ *[]interface{}, _ PlaceholderGenerator) string {
	return c.Column + " IS NULL"
}

// ===============================================================
// IS NOT NULL
// ===============================================================

func IsNotNull(column string) Condition {
	return &conditionIsNotNull{
		Column: column,
	}
}

type conditionIsNotNull struct {
	Column string
}

func (c *conditionIsNotNull) toSQL(_ *int, _ *[]interface{}, _ PlaceholderGenerator) string {
	return c.Column + " IS NOT NULL"
}

// ===============================================================
// BETWEEN
// ===============================================================

func Between(column string, value1, value2 interface{}) Condition {
	return &conditionBetween{
		Column: column,
		Value1: value1,
		Value2: value2,
	}
}

type conditionBetween struct {
	Column string
	Value1 interface{}
	Value2 interface{}
}

func (c *conditionBetween) toSQL(placeholderStartAt *int, queryArgs *[]interface{}, placeHolderGenerator PlaceholderGenerator) string {
	s := c.Column + " BETWEEN "
	s += placeHolderGenerator(*placeholderStartAt)
	*placeholderStartAt++
	*queryArgs = append(*queryArgs, c.Value1)
	s += " AND "
	s += placeHolderGenerator(*placeholderStartAt)
	*placeholderStartAt++
	*queryArgs = append(*queryArgs, c.Value2)
	return s
}

// ===============================================================
// NOT BETWEEN
// ===============================================================

func NotBetween(column string, value1, value2 interface{}) Condition {
	return &conditionNotBetween{
		Column: column,
		Value1: value1,
		Value2: value2,
	}
}

type conditionNotBetween struct {
	Column string
	Value1 interface{}
	Value2 interface{}
}

func (c *conditionNotBetween) toSQL(placeholderStartAt *int, queryArgs *[]interface{}, placeHolderGenerator PlaceholderGenerator) string {
	s := c.Column + " NOT BETWEEN "
	s += placeHolderGenerator(*placeholderStartAt)
	*placeholderStartAt++
	*queryArgs = append(*queryArgs, c.Value1)
	s += " AND "
	s += placeHolderGenerator(*placeholderStartAt)
	*placeholderStartAt++
	*queryArgs = append(*queryArgs, c.Value2)
	return s
}
