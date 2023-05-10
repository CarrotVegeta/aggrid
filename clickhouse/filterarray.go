package clickhouse

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
func (ft *FilterArray) BuildSql(k string, va any, t constant.OperatorType, f ...constant.F) (*utils.QueryFilter, error) {
	var v []any
	if len(va.([]any)) > 0 {
		switch va.([]any)[0].(type) {
		case float64:
			v = ft.parseFloat2Uint(va)
		default:
			v = va.([]any)
		}
	}
	switch t {
	case constant.Equals:
		ft.Equals(k, v)
	case constant.NotEqual:
		ft.NotEquals(k, v)
	case constant.Contains:
		ft.Contains(k, v)
	case constant.NotContains:
		ft.NotContains(k, v)
	case constant.Between:
		ft.Contains(k, v)
	case constant.NoBetween:
		ft.NotContains(k, v)
	case constant.StartsWith:
		ft.StartsWith(k, v)
	case constant.EndsWith:
		ft.EndsWith(k, v)
	case constant.HasAll:
		ft.HasAll(k, v)
	case constant.NotHasAll:
		ft.NotHasAll(k, v)
	case constant.Blank:
		ft.Blank(k, v)
	case constant.NotBlank:
		ft.NotBlank(k, v)
	case constant.Regexp:
		ft.Regexp(k, v)
	default:
		return nil, fmt.Errorf("filter type is invalid : %v", t)
	}
	return ft.QF, nil
}

// Regexp 相等
func (ft *FilterArray) Regexp(k string, v any) *FilterArray {
	query := fmt.Sprintf("match(%s,?) ", k)
	ft.QF.And(query, v)
	return nil
}

// Equals 相等
func (ft *FilterArray) Equals(k string, v []any) *FilterArray {
	var argarray string
	for range v {
		if argarray == "" {
			argarray = "[?"
			continue
		}
		argarray += ",?"
	}
	argarray += "]"
	query := fmt.Sprintf("arraySort(%s) = arraySort(%s) ", k, argarray)
	ft.QF.And(query, v...)
	return ft
}

// NotEquals 不相等
func (ft *FilterArray) NotEquals(k string, v []any) *FilterArray {
	var argarray string
	for range v {
		if argarray == "" {
			argarray = "[?"
			continue
		}
		argarray += ",?"
	}
	argarray += "]"
	query := fmt.Sprintf("arraySort(%s) <> arraySort(%s) ", k, argarray)
	ft.QF.And(query, v...)
	return ft
}

// EndsWith 结尾
func (ft *FilterArray) EndsWith(k string, va []any) *FilterArray {
	if len(va) == 0 {
		return ft
	}
	var query string
	var args []any
	for _, v := range va {
		if query == "" {
			query = fmt.Sprintf("match(%s,?) ", k)
			args = append(args, fmt.Sprintf("%v$", v))
			continue
		}
		query += fmt.Sprintf("OR match(%s,?) ", k)
		args = append(args, fmt.Sprintf("%v$", v))
	}
	query = "(" + query + ")"
	ft.QF.And(query, args...)
	return ft
}

// StartsWith 前缀
func (ft *FilterArray) StartsWith(k string, va []any) *FilterArray {
	if len(va) == 0 {
		return ft
	}
	var query string
	var args []any
	for _, v := range va {
		if query == "" {
			query = fmt.Sprintf("match(%s,?) ", k)
			args = append(args, fmt.Sprintf("^%v", v))
			continue
		}
		query += fmt.Sprintf("OR match(%s,?) ", k)
		args = append(args, fmt.Sprintf("^%v", v))
	}
	query = "(" + query + ")"
	ft.QF.And(query, args...)
	return ft
}
func (ft *FilterArray) parseFloat2Uint(v any) []any {
	var args []any
	for _, va := range v.([]any) {
		args = append(args, (uint)(va.(float64)))
	}
	return args
}

// Contains 包含
func (ft *FilterArray) Contains(k string, v []any) *FilterArray {
	var query string
	if len(v) == 1 {
		query = fmt.Sprintf("%s = ? ", k)
	} else if len(v) > 1 {
		query = fmt.Sprintf("%s IN (?) ", k)
	} else {
		return ft
	}
	ft.QF.And(query, v)
	return ft
}

// NotContains 不包含
func (ft *FilterArray) NotContains(k string, v []any) *FilterArray {
	var query string
	if len(v) == 1 {
		query = fmt.Sprintf("%s <> ? ", k)
	} else if len(v) > 1 {
		query = fmt.Sprintf("%s NOT IN (?) ", k)
	} else {
		return ft
	}
	ft.QF.And(query, v)
	return ft
}

// HasAll 包含
func (ft *FilterArray) HasAll(k string, v []any) *FilterArray {
	if len(v) == 0 {
		return ft
	}
	var argarray string
	for range v {
		if argarray == "" {
			argarray = "[?"
			continue
		}
		argarray += ",?"
	}
	argarray += "]"
	query := fmt.Sprintf("hasAll(%s,%s) ", k, argarray)
	ft.QF.And(query, v...)
	return ft
}

// NotHasAll 不包含
func (ft *FilterArray) NotHasAll(k string, v []any) *FilterArray {
	if len(v) == 0 {
		return ft
	}
	var argarray string
	for range v {
		if argarray == "" {
			argarray = "[?"
			continue
		}
		argarray += ",?"
	}
	argarray += "]"
	query := fmt.Sprintf("not hasAll(%s,%s) ", k, argarray)
	ft.QF.And(query, v...)
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
