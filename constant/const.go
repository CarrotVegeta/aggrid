package constant

import "github.com/CarrotVegeta/aggrid/utils"

type FilterMethodType string

const (
	FilterMethodWhere  FilterMethodType = "WHERE"
	FilterMethodHaving FilterMethodType = "Having"
)
const (
	AND = "AND"
	OR  = "OR"
)

type F func(k string, v any) *utils.QueryFilter

type OperatorType string

// 判断逻辑
const (
	LessThan           OperatorType = "lessThan"
	Equals             OperatorType = "equals"
	NotEqual           OperatorType = "notEqual"
	GreaterThanOrEqual OperatorType = "greaterThanOrEqual"
	GreaterThan        OperatorType = "greaterThan"
	InRange            OperatorType = "inRange"
	Blank              OperatorType = "blank"
	NotBlank           OperatorType = "notBlank"
	LessThanOrEqual    OperatorType = "lessThanOrEqual"
	Contains           OperatorType = "contains"
	NotContains        OperatorType = "notContains"
	StartsWith         OperatorType = "startsWith"
	EndsWith           OperatorType = "endsWith"
	HasAll             OperatorType = "hasAll"
	NotHasAll          OperatorType = "notHasAll"
	Regexp             OperatorType = "regexp"
	Between            OperatorType = "between"
	NoBetween          OperatorType = "noBetween"
)

const (
	DefaultStorage = MySql
	MySql          = "mysql"
	Clickhouse     = "clickhouse"
)

type FilterType string

// FilterType 类型
const (
	Text   FilterType = "text"
	Number FilterType = "number"
	Date   FilterType = "date"
	Array  FilterType = "array"
	Set    FilterType = "set"
)
