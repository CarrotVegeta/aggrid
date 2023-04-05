package agtwo

import (
	"fmt"
)

type FilterHandler interface {
	Parse(k string, c []byte) error
	BuildQuery() (*QueryFilter, error)
}

type Filter struct {
	Handler  FilterHandler
	Operator string `json:"operator"`
}

func (f *Filter) H() FilterHandler {
	if f.Operator != "" {
		f.Handler = &FilterModelCondition{}
		return f.Handler
	}
	f.Handler = &FilterModel{}
	return f.Handler
}

type FilterTypeSqlHandler interface {
	New() FilterTypeSqlHandler
	BuildSql(k string, v any, t string, f ...F) (*QueryFilter, error)
}

// 判断逻辑
const (
	LessThan           = "lessThan"
	Equals             = "equals"
	NotEqual           = "notEqual"
	GreaterThanOrEqual = "greaterThanOrEqual"
	GreaterThan        = "greaterThan"
	InRange            = "inRange"
	Blank              = "blank"
	NotBlank           = "notBlank"
	LessThanOrEqual    = "lessThanOrEqual"
	Contains           = "contains"
	NotContains        = "notContains"
	StartsWith         = "startsWith"
	EndsWith           = "endsWith"
)

// FilterType 类型
const (
	Text   = "text"
	Number = "number"
	Date   = "date"
	Array  = "array"
)

var FilterTypeSqlHandlerM = map[string]FilterTypeSqlHandler{}

func init() {
	registerFilterTypeSqlHandler(Text, &FilterText{})
	registerFilterTypeSqlHandler(Number, &FilterNumber{})
	registerFilterTypeSqlHandler(Date, &FilterDate{})
	registerFilterTypeSqlHandler(Array, &FilterArray{})
}
func registerFilterTypeSqlHandler(k string, handler FilterTypeSqlHandler) {
	FilterTypeSqlHandlerM[k] = handler
}
func RegisterFilterTypeSqlHandler(k string, h FilterTypeSqlHandler) {
	registerFilterTypeSqlHandler(k, h)
}
func getFilterTypeSqlHandler(filterType string) (FilterTypeSqlHandler, error) {
	if _, ok := FilterTypeSqlHandlerM[filterType]; !ok {
		return nil, fmt.Errorf("invalid filter-type:%v", filterType)
	}
	return FilterTypeSqlHandlerM[filterType], nil
}
func NewFilterSqlHandler(filterType string) (FilterTypeSqlHandler, error) {
	h, err := getFilterTypeSqlHandler(filterType)
	if err != nil {
		return nil, err
	}
	return h.New(), nil
}
