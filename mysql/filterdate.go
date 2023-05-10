package mysql

import (
	"fmt"
	"github.com/CarrotVegeta/aggrid/constant"
	"github.com/CarrotVegeta/aggrid/interfaces"
	"github.com/CarrotVegeta/aggrid/utils"
	"time"
)

type FilterDate struct {
	QF *utils.QueryFilter
}

//type F func(k string, v any) *aggrid.QueryFilter

func (fd *FilterDate) New() interfaces.FilterTypeSqlHandler {
	return &FilterDate{QF: &utils.QueryFilter{}}
}
func (fd *FilterDate) BuildSql(k string, v any, t constant.OperatorType, f ...constant.F) (*utils.QueryFilter, error) {
	switch t {
	case constant.LessThan:
		fd.LessThan(k, v)
	case constant.Equals:
		fd.Equals(k, v, f...)
	case constant.NotEqual:
		fd.NotEqual(k, v)
	case constant.GreaterThan:
		fd.GreaterThan(k, v, f...)
	case constant.InRange:
		fd.InRange(k, v)
	case constant.Blank:
		fd.Blank(k)
	case constant.NotBlank:
		fd.NotBlank(k)
	default:
		return nil, fmt.Errorf("filter type is invalid : %v", t)
	}

	return fd.QF, nil
}

func (fd *FilterDate) Blank(k string) *FilterDate {
	fd.QF.And(k + "IS  NULL ")
	return fd
}
func (fd *FilterDate) NotBlank(k string) *FilterDate {
	fd.QF.And(k + "IS NOT NULL ")
	return fd
}

// InRange 范围内
func (fd *FilterDate) InRange(k string, v any, f ...constant.F) *FilterDate {
	if len(f) > 0 {
		qf := f[0](k, v)
		fd.QF.And(qf.Query, qf.Args...)
		return fd
	}
	arr := v.([]string)
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", arr[0], time.Local)
	lTime, _ := time.ParseInLocation("2006-01-02 15:04:05", arr[1], time.Local)
	endTime := lTime.AddDate(0, 0, 1).UnixMilli()
	fd.QF.And(fmt.Sprintf("%s >= ? AND %s <= ? ", k, k), startTime.UnixMilli(), endTime)
	return fd
}

// GreaterThan 大于
func (fd *FilterDate) GreaterThan(k string, v any, f ...constant.F) *FilterDate {
	if len(f) > 0 {
		qf := f[0](k, v)
		fd.QF.And(qf.Query, qf.Args...)
		return fd
	}
	tt, _ := time.Parse("2006-01-02 15:04:05", v.(string))
	query := fmt.Sprintf("%s > ? ", k)
	fd.QF.And(query, tt.UnixMilli())
	return fd
}

// NotEqual 不等于
func (fd *FilterDate) NotEqual(k string, v any, f ...constant.F) *FilterDate {
	if len(f) > 0 {
		qf := f[0](k, v)
		fd.QF.And(qf.Query, qf.Args...)
		return fd
	}
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", v.(string), time.Local)
	lTime, _ := time.ParseInLocation("2006-01-02 15:04:05", v.(string), time.Local)
	endTime := lTime.AddDate(0, 0, 1).UnixMilli()
	fd.QF.And(fmt.Sprintf("%s > ? AND %s < ? ", k, k), endTime, startTime.UnixMilli())
	return fd
}

// Equals 等于
func (fd *FilterDate) Equals(k string, v any, f ...constant.F) *FilterDate {
	if len(f) > 0 {
		qf := f[0](k, v)
		fd.QF.And(qf.Query, qf.Args...)
		return fd
	}
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", v.(string), time.Local)
	lTime, _ := time.ParseInLocation("2006-01-02 15:04:05", v.(string), time.Local)
	endTime := lTime.AddDate(0, 0, 1).UnixMilli()
	fd.QF.And(fmt.Sprintf("%s >= ? AND %s <= ? ", k, k), startTime, endTime)
	return fd
}

// LessThan 小于
func (fd *FilterDate) LessThan(k string, v any, f ...constant.F) *FilterDate {
	if len(f) > 0 {
		qf := f[0](k, v)
		fd.QF.And(qf.Query, qf.Args...)
		return fd
	}
	tt, _ := time.Parse("2006-01-02 15:04:05", v.(string))
	query := fmt.Sprintf("%s < ? ", k)
	fd.QF.And(query, tt)
	return fd
}
