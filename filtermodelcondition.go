package agtwo

import (
	"encoding/json"
	"fmt"
)

// FilterModelCondition 多个查询条件
type FilterModelCondition struct {
	Key        string `json:"key"`
	FilterType string `json:"filterType"`
	Operator   string `json:"operator"`
	Condition1 any    `json:"condition1"`
	Condition2 any    `json:"condition2"`
}

// Parse 单个查询条件的参数解析
func (f *FilterModelCondition) Parse(k string, c []byte) error {
	f.Key = k
	if err := json.Unmarshal(c, f); err != nil {
		return fmt.Errorf("parse filter model failed:%v", err.Error())
	}
	return nil
}
func (f *FilterModelCondition) GetFilterAndType(p []byte) (any, string, error) {
	handler, err := NewFilterTypeHandler(f.FilterType, p)
	if err != nil {
		return nil, "", err
	}
	filter, err := handler.GetFilter()
	if err != nil {
		return nil, "", nil
	}
	return filter, handler.GetType(), nil
}

const (
	AND = "AND"
	OR  = "OR"
)

// BuildQuery 单个查询条件的sql生成
func (f *FilterModelCondition) BuildQuery() (*QueryFilter, error) {
	condBs1, _ := json.Marshal(f.Condition1)
	f1, t1, err := f.GetFilterAndType(condBs1)
	if err != nil {
		return nil, err
	}
	h, err := NewFilterSqlHandler(f.FilterType)
	if err != nil {
		return nil, err
	}
	qf1, err := h.New().BuildSql(f.Key, f1, t1)
	if err != nil {
		return nil, err
	}
	condBs2, _ := json.Marshal(f.Condition2)
	f2, t2, err := f.GetFilterAndType(condBs2)
	if err != nil {
		return nil, err
	}
	qf2, err := h.New().BuildSql(f.Key, f2, t2)
	if err != nil {
		return nil, err
	}
	switch f.Operator {
	case AND:
		qf1.And(qf2.Query, qf2.Args...)
	case OR:
		qf1.Or(qf2.Query, qf2.Args...)
	default:
		return nil, fmt.Errorf("invalid operatof : %v", f.Operator)
	}
	return qf1, nil

}
