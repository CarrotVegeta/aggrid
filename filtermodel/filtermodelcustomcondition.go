package filtermodel

import (
	"encoding/json"
	"fmt"
	"github.com/CarrotVegeta/aggrid/constant"
	"github.com/CarrotVegeta/aggrid/filtertype"
	"github.com/CarrotVegeta/aggrid/utils"
)

// FilterModelCustomCondition 多个查询条件
type FilterModelCustomCondition struct {
	Key        string              `json:"key"`
	FilterType constant.FilterType `json:"filterType"`
	Operator   string              `json:"operator"`
	Conditions []any               `json:"conditions"`
}

// Parse 单个查询条件的参数解析
func (f *FilterModelCustomCondition) Parse(k string, c []byte) error {
	f.Key = k
	if err := json.Unmarshal(c, f); err != nil {
		return fmt.Errorf("parse filter model failed:%v", err.Error())
	}
	return nil
}
func (f *FilterModelCustomCondition) GetFilterAndType(p []byte) (any, constant.OperatorType, error) {
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

type FilterCondition struct {
	FilterType constant.FilterType `json:"filterType"`
	Type       string              `json:"type"`
	Filter     string              `json:"filter"`
}

// BuildQuery 单个查询条件的sql生成
func (f *FilterModelCustomCondition) BuildQuery(service *filtertype.FilterTypeSqlService) (*utils.QueryFilter, error) {
	qf := &utils.QueryFilter{}
	for _, v := range f.Conditions {
		condBs1, _ := json.Marshal(v)
		condition := &FilterCondition{}
		if err := json.Unmarshal(condBs1, &condition); err != nil {
			return nil, fmt.Errorf("unmarshal filter condition fail:%v", err)
		}
		f.FilterType = condition.FilterType
		f1, t1, err := f.GetFilterAndType(condBs1)
		if err != nil {
			return nil, err
		}
		h, err := service.NewFilterSqlHandler(condition.FilterType)
		if err != nil {
			return nil, err
		}
		qf1, err := h.New().BuildSql(f.Key, f1, t1)
		if err != nil {
			return nil, err
		}
		switch f.Operator {
		case constant.AND:
			qf.And(qf1.Query, qf1.Args...)
		case constant.OR:
			qf.Or(qf1.Query, qf1.Args...)
		default:
			return nil, fmt.Errorf("invalid operatof : %v", f.Operator)
		}
	}
	return qf, nil

}
