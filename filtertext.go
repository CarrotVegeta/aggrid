package agtwo

import (
	"fmt"
)

type FilterText struct {
	Type string `json:"type"`
	QF   *QueryFilter
}

func (ft *FilterText) New() FilterTypeSqlHandler {
	return &FilterText{QF: &QueryFilter{}}
}
func (ft *FilterText) BuildSql(k string, v any, t string, f ...F) (*QueryFilter, error) {
	switch t {
	case Contains:
		ft.Contains(k, v)
	case NotContains:
		ft.Locate(k, v)
	case Equals:
		ft.Equals(k, v)
	case NotEqual:
		ft.NotEqual(k, v)
	case StartsWith:
		ft.StartsWith(k, v)
	case EndsWith:
		ft.EndsWith(k, v)
	case Blank:
		ft.Blank(k)
	case NotBlank:
		ft.NotBlank(k)
	default:
		return nil, fmt.Errorf("filter type is invalid : %v", t)
	}
	return ft.QF, nil
}
func (ft *FilterText) NotBlank(k string) *FilterText {
	ft.QF.And(fmt.Sprintf("%s <> '' OR %s IS NOT NULL ", k, k))
	return ft
}
func (ft *FilterText) Blank(k string) *FilterText {
	ft.QF.And(fmt.Sprintf("%s = '' OR %s is NULL ", k, k))
	return ft
}
func (ft *FilterText) EndsWith(k string, v any) *FilterText {
	ft.QF.And(fmt.Sprintf("RIGHT(%s,%d)= ? ", k, len(v.(string))), v)
	return ft
}
func (ft *FilterText) StartsWith(k string, v any) *FilterText {
	ft.QF.And(fmt.Sprintf("LEFT(%s,%d)= ? ", k, len(v.(string))), v)
	return ft
}
func (ft *FilterText) NotEqual(k string, v any) *FilterText {
	ft.QF.And(fmt.Sprintf("%s <> ? ", k), v)
	return ft
}

// Equals 等于
func (ft *FilterText) Equals(k string, v any) *FilterText {
	ft.QF.And(fmt.Sprintf("%s = ? ", k), v)
	return ft
}

// Locate 不包含
func (ft *FilterText) Locate(k string, v any) *FilterText {
	query := fmt.Sprintf("locate(%s,?) = 0 ", k)
	ft.QF.And(query, v)
	return ft
}

// Contains 包含
func (ft *FilterText) Contains(k string, v any) *FilterText {
	query := fmt.Sprintf("%s LIKE ? ", k)
	ft.QF.And(query, fmt.Sprintf("%%%s%%", v))
	return ft
}
