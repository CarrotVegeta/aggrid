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
	BuildSql(k string, v any, t OperatorType, f ...F) (*QueryFilter, error)
}
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

var DefaultFilterTypeSqlM = map[FilterType]FilterTypeSqlHandler{}

func init() {
	registerFilterTypeSqlHandler(Text, &FilterText{})
	registerFilterTypeSqlHandler(Number, &FilterNumber{})
	registerFilterTypeSqlHandler(Date, &FilterDate{})
	registerFilterTypeSqlHandler(Array, &FilterArray{})
	registerFilterTypeSqlHandler(Set, &FilterSet{})
}
func registerFilterTypeSqlHandler(k FilterType, handler FilterTypeSqlHandler) {
	DefaultFilterTypeSqlM[k] = handler
}

type FilterTypeSqlService struct {
	FilterTypeSqlHandlerM map[FilterType]FilterTypeSqlHandler
}

func NewFilterTypeSqlService() *FilterTypeSqlService {
	f := &FilterTypeSqlService{
		FilterTypeSqlHandlerM: make(map[FilterType]FilterTypeSqlHandler),
	}
	f.initTypeSqlHandlerM()
	return f
}

func (f *FilterTypeSqlService) initTypeSqlHandlerM() {
	for k, v := range DefaultFilterTypeSqlM {
		f.FilterTypeSqlHandlerM[k] = v
	}
}
func (f *FilterTypeSqlService) registerFilterTypeSqlHandler(k FilterType, handler FilterTypeSqlHandler) {
	f.FilterTypeSqlHandlerM[k] = handler
}

func (f *FilterTypeSqlService) getFilterTypeSqlHandler(filterType FilterType) (FilterTypeSqlHandler, error) {
	if _, ok := f.FilterTypeSqlHandlerM[filterType]; !ok {
		return nil, fmt.Errorf("invalid filter-type:%v", filterType)
	}
	return f.FilterTypeSqlHandlerM[filterType].New(), nil
}
func (f *FilterTypeSqlService) NewFilterSqlHandler(filterType FilterType) (FilterTypeSqlHandler, error) {
	h, err := f.getFilterTypeSqlHandler(filterType)
	if err != nil {
		return nil, err
	}
	return h, nil
}
