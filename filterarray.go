package agtwo

import (
	"fmt"
)

type FilterArray struct {
	Type string `json:"type"`
	QF   *QueryFilter
}

func (ft *FilterArray) New() FilterTypeSqlHandler {
	return &FilterArray{QF: &QueryFilter{}}
}
func (ft *FilterArray) BuildSql(k string, v any, t OperatorType, f ...F) (*QueryFilter, error) {
	switch t {
	case Contains:
		ft.Contains(k, v)
	case StartsWith:
		ft.StartsWith(k, v)
	case EndsWith:
		ft.EndsWith(k, v)
	default:
		return nil, fmt.Errorf("filter type is invalid : %v", t)
	}
	return ft.QF, nil
}

// EndsWith 结尾
func (ft *FilterArray) EndsWith(k string, v any) *FilterArray {
	arr := v.([]any)
	if len(arr) == 0 {
		return ft
	}
	var query string
	var args []any
	for _, v := range arr {
		if query == "" {
			query = fmt.Sprintf("%s REGEXP ? ", k)
			args = append(args, fmt.Sprintf("%v$", v))
			continue
		}
		query += fmt.Sprintf("OR %s REGEXP ? ", k)
		args = append(args, fmt.Sprintf("%v$", v))
	}
	query = "(" + query + ")"
	ft.QF.And(query, args...)
	return ft
}

// StartsWith 前缀
func (ft *FilterArray) StartsWith(k string, v any) *FilterArray {
	arr := v.([]any)
	if len(arr) == 0 {
		return ft
	}
	var query string
	var args []any
	for _, v := range arr {
		if query == "" {
			query = fmt.Sprintf("%s REGEXP ? ", k)
			args = append(args, fmt.Sprintf("^%v", v))
			continue
		}
		query += fmt.Sprintf("OR %s REGEXP ? ", k)
		args = append(args, fmt.Sprintf("^%v", v))
	}
	query = "(" + query + ")"
	ft.QF.And(query, args...)
	return ft
}

// Contains 包含
func (ft *FilterArray) Contains(k string, v any) *FilterArray {
	array := v.([]any)
	var query string
	if len(array) == 1 {
		query = fmt.Sprintf("%s = ? ", k)
	} else if len(array) > 1 {
		query = fmt.Sprintf("%s IN (?) ", k)
	} else {
		return ft
	}
	ft.QF.And(query, v)
	return ft
}
