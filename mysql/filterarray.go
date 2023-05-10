package mysql

import (
	"fmt"
	"github.com/CarrotVegeta/aggrid/constant"
	"github.com/CarrotVegeta/aggrid/interfaces"
	"github.com/CarrotVegeta/aggrid/utils"
)

type FilterArray struct {
	Type string `json:"type"`
	QF   *utils.QueryFilter
}

func (ft *FilterArray) New() interfaces.FilterTypeSqlHandler {
	return &FilterArray{QF: &utils.QueryFilter{}}
}
func (ft *FilterArray) BuildSql(k string, v any, t constant.OperatorType, f ...constant.F) (*utils.QueryFilter, error) {
	switch t {
	case constant.StartsWith:
		ft.StartsWith(k, v)
	case constant.EndsWith:
		ft.EndsWith(k, v)
	case constant.Contains:
		ft.Contains(k, v)
	case constant.NotContains:
		ft.NotContains(k, v)
	case constant.HasAll:
		ft.HasAll(k, v)
	case constant.NotHasAll:
		ft.NotHasAll(k, v)
	case constant.Blank:
		ft.Blank(k, v)
	case constant.NotBlank:
		ft.NotBlank(k, v)
	case constant.Regexp:

	default:
		return nil, fmt.Errorf("filter type is invalid : %v", t)
	}
	return ft.QF, nil
}

// Regexp 相等
func (ft *FilterArray) Regexp(k string, v any) *FilterArray {
	return nil
}

// Equals 相等
func (ft *FilterArray) Equals(k string, v any) *FilterArray {
	return nil
}

// NotEquals 不相等
func (ft *FilterArray) NotEquals(k string, v any) *FilterArray {
	return nil
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

// NotContains 不包含
func (ft *FilterArray) NotContains(k string, v any) *FilterArray {
	for _, v := range v.([]any) {
		query := fmt.Sprintf("%s <> ? ", k)
		ft.QF.Or(query, v)
	}
	return ft
}

// HasAll 包含
func (ft *FilterArray) HasAll(k string, v any) *FilterArray {
	query := fmt.Sprintf("hasAll(%s,?) ", k)
	ft.QF.And(query, v)
	return ft
}

// NotHasAll 不包含
func (ft *FilterArray) NotHasAll(k string, v any) *FilterArray {
	query := fmt.Sprintf("hasAll(%s,?) = 0 ", k)
	ft.QF.And(query, v)
	return ft
}

// NotBlank 不为空
func (ft *FilterArray) NotBlank(k string, v any) *FilterArray {
	query := fmt.Sprintf("notEmpty(%s) ", k)
	ft.QF.And(query)
	return ft
}

// Blank 为空
func (ft *FilterArray) Blank(k string, v any) *FilterArray {
	query := fmt.Sprintf("empty(%s)  ", k)
	ft.QF.And(query)
	return ft
}
