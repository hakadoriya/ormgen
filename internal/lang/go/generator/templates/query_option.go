package sandbox

import (
	"strings"
)

type (
	ResultOption interface {
		applyResultOption(c *queryConfig)
		QueryOption
	}
	QueryOption interface {
		applyQueryOption(c *queryConfig)
	}
)

type PlaceholderGen func(args ...interface{}) string
type PlaceholderGenGen func(placeholderStartAt *int, queryArgs *[]interface{}) PlaceholderGen

type queryConfig struct {
	QueryPrefix        string
	PlaceholderStartAt int
	PlaceholderGenGen  PlaceholderGenGen
	Where              *where
	OrderBy            *orderBy
	Limit              int
}

func (c *queryConfig) ToSQL(query string) (string, []interface{}) {
	var q string
	var args []interface{}
	placeholderGen := c.PlaceholderGenGen(&c.PlaceholderStartAt, &args)
	if c.QueryPrefix != "" {
		q += c.QueryPrefix + " "
	}
	q += query
	if c.Where != nil {
		q += " WHERE " + c.Where.Root.toSQL(placeholderGen)
	}
	if c.OrderBy != nil {
		q += " ORDER BY " + strings.Join(c.OrderBy.columns, ", ")
		if c.OrderBy.desc {
			q += " DESC"
		}
	}
	if c.Limit > 0 {
		q += " LIMIT " + placeholderGen(c.Limit)
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

func (o *withListOptionQueryPrefix) applyResultOption(c *queryConfig) {
	c.QueryPrefix += o.prefix
}

func (o *withListOptionQueryPrefix) applyQueryOption(c *queryConfig) {
	o.applyResultOption(c)
}

// ===============================================================
// ORDER BY
// ===============================================================

func OrderBy(orderBy ...string) ResultOption {
	return &withListOptionOrderBy{
		columns: orderBy,
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

func (o *withListOptionOrderBy) applyResultOption(c *queryConfig) {
	c.OrderBy = &orderBy{
		columns: o.columns,
		desc:    o.desc,
	}
}

func (o *withListOptionOrderBy) applyQueryOption(c *queryConfig) {
	o.applyResultOption(c)
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

func (o *withListOptionLimit) applyResultOption(c *queryConfig) {
	c.Limit = o.limit
}

func (o *withListOptionLimit) applyQueryOption(c *queryConfig) {
	o.applyResultOption(c)
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
	toSQL(placeHolderGen PlaceholderGen) string
}

type conditionCompound struct {
	Op         string // "AND" or "OR"
	Conditions []Condition
}

func (cc *conditionCompound) toSQL(placeHolderGen PlaceholderGen) string {
	s := "("
	for i, c := range cc.Conditions {
		if i != 0 {
			s += " " + cc.Op + " "
		}
		s += c.toSQL(placeHolderGen)
	}
	s += ")"
	return s
}

var (
	_ QueryOption = (*withListOptionWhere)(nil)
)

type withListOptionWhere struct {
	where *where
}

func (o *withListOptionWhere) applyQueryOption(c *queryConfig) {
	c.Where = o.where
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

func (c *conditionEqual) toSQL(placeHolderGen PlaceholderGen) string {
	return c.Column + " = " + placeHolderGen(c.Value)
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

func (c *conditionNotEqual) toSQL(placeHolderGen PlaceholderGen) string {
	return c.Column + " <> " + placeHolderGen(c.Value)
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

func (c *conditionIn) toSQL(placeHolderGen PlaceholderGen) string {
	s := c.Column + " IN ("
	for i, v := range c.Values {
		if i != 0 {
			s += ", "
		}
		s += placeHolderGen(v)
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

func (c *conditionNotIn) toSQL(placeHolderGen PlaceholderGen) string {
	s := c.Column + " NOT IN ("
	for i, v := range c.Values {
		if i != 0 {
			s += ", "
		}
		s += placeHolderGen(v)
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

func (c *conditionGreaterThan) toSQL(placeHolderGen PlaceholderGen) string {
	return c.Column + " > " + placeHolderGen(c.Value)
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

func (c *conditionGreaterThanOrEqual) toSQL(placeHolderGen PlaceholderGen) string {
	return c.Column + " >= " + placeHolderGen(c.Value)
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

func (c *conditionLessThan) toSQL(placeHolderGen PlaceholderGen) string {
	return c.Column + " < " + placeHolderGen(c.Value)
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

func (c *conditionLessThanOrEqual) toSQL(placeHolderGen PlaceholderGen) string {
	return c.Column + " <= " + placeHolderGen(c.Value)
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

func (c *conditionLike) toSQL(placeHolderGen PlaceholderGen) string {
	return c.Column + " LIKE " + placeHolderGen(c.Value)
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

func (c *conditionNotLike) toSQL(placeHolderGen PlaceholderGen) string {
	return c.Column + " NOT LIKE " + placeHolderGen(c.Value)
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

func (c *conditionIsNull) toSQL(placeHolderGen PlaceholderGen) string {
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

func (c *conditionIsNotNull) toSQL(placeHolderGen PlaceholderGen) string {
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

func (c *conditionBetween) toSQL(placeHolderGen PlaceholderGen) string {
	return c.Column + " BETWEEN " + placeHolderGen(c.Value1) + " AND " + placeHolderGen(c.Value2)
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

func (c *conditionNotBetween) toSQL(placeHolderGen PlaceholderGen) string {
	return c.Column + " NOT BETWEEN " + placeHolderGen(c.Value1) + " AND " + placeHolderGen(c.Value2)
}
