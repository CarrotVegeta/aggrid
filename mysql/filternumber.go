package mysql

import (
	"fmt"
	"github.com/CarrotVegeta/aggrid/constant"
	"github.com/CarrotVegeta/aggrid/interfaces"
	"github.com/CarrotVegeta/aggrid/utils"
)

type FilterNumber struct {
	QF *utils.QueryFilter
}

func (fn *FilterNumber) New() interfaces.FilterTypeSqlHandler {
	return &FilterNumber{QF: &utils.QueryFilter{}}
}
func (fn *FilterNumber) BuildSql(k string, v any, t constant.OperatorType, f ...constant.F) (*utils.QueryFilter, error) {
	switch t {
	case constant.LessThan:
		fn.LessThan(k, v)
	case constant.LessThanOrEqual:
		fn.LessThanOrEqual(k, v)
	case constant.Equals:
		fn.Equals(k, v)
	case constant.NotEqual:
		fn.NotEqual(k, v)
	case constant.GreaterThan:
		fn.GreaterThan(k, v)
	case constant.GreaterThanOrEqual:
		fn.GreaterThanOrEqual(k, v)
	case constant.InRange:
		fn.InRange(k, v)
	case constant.Blank:
		fn.Blank(k)
	case constant.NotBlank:
		fn.NotBlank(k)
	default:
		return nil, fmt.Errorf("filter type is invalid : %v", t)
	}
	return fn.QF, nil
}
func (fn *FilterNumber) Blank(k string) *FilterNumber {
	fn.QF.And(k + " IS  NULL ")
	return fn
}
func (fn *FilterNumber) NotBlank(k string) *FilterNumber {
	fn.QF.And(k + " IS NOT NULL ")
	return fn
}
func (fn *FilterNumber) InRange(k string, v any) *FilterNumber {
	array := v.([]int64)
	fn.QF.And(fmt.Sprintf("%s >= ? AND %s <= ? ", k, k), array[0], array[1])
	return fn
}

// GreaterThanOrEqual 大于或等于
func (fn *FilterNumber) GreaterThanOrEqual(k string, v any) *FilterNumber {
	fn.QF.And(fmt.Sprintf("%s >= ? ", k), v)
	return fn
}

// GreaterThan 大于
func (fn *FilterNumber) GreaterThan(k string, v any) *FilterNumber {
	fn.QF.And(fmt.Sprintf("%s > ? ", k), v)
	return fn
}

// NotEqual 不等于
func (fn *FilterNumber) NotEqual(k string, v any) *FilterNumber {
	fn.QF.And(fmt.Sprintf("%s <> ? ", k), v)
	return fn
}

// Equals 等于
func (fn *FilterNumber) Equals(k string, v any) *FilterNumber {
	fn.QF.And(fmt.Sprintf("%s = ? ", k), v)
	return fn
}

// LessThanOrEqual 小于等于
func (fn *FilterNumber) LessThanOrEqual(k string, v any) *FilterNumber {
	query := fmt.Sprintf("%s <= ? ", k)
	fn.QF.And(query, v)
	return fn
}

// LessThan 小于
func (fn *FilterNumber) LessThan(k string, v any) *FilterNumber {
	query := fmt.Sprintf("%s < ? ", k)
	fn.QF.And(query, v)
	return fn
}
