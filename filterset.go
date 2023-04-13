package agtwo

import (
	"fmt"
)

type FilterSet struct {
	Type OperatorType `json:"type"`
	QF   *QueryFilter
}

func (ft *FilterSet) New() FilterTypeSqlHandler {
	return &FilterSet{QF: &QueryFilter{}}
}
func (ft *FilterSet) BuildSql(k string, v any, t OperatorType, f ...F) (*QueryFilter, error) {
	switch t {
	case InRange:
		ft.In(k, v)
	}
	return ft.QF, nil
}

func (ft *FilterSet) In(k string, v any) *FilterSet {
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
