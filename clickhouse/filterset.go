package clickhouse

import (
	"fmt"
	"github.com/CarrotVegeta/aggrid/constant"
	"github.com/CarrotVegeta/aggrid/interfaces"
	"github.com/CarrotVegeta/aggrid/utils"
)

type FilterSet struct {
	Type constant.OperatorType `json:"type"`
	QF   *utils.QueryFilter
}

func (ft *FilterSet) New() interfaces.FilterTypeSqlHandler {
	return &FilterSet{QF: &utils.QueryFilter{}}
}
func (ft *FilterSet) BuildSql(k string, v any, t constant.OperatorType, f ...constant.F) (*utils.QueryFilter, error) {
	switch t {
	case constant.InRange:
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
