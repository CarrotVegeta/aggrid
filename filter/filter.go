package filter

import (
	"github.com/CarrotVegeta/aggrid/filtermodel"
	"github.com/CarrotVegeta/aggrid/filtertype"
)

type Filter struct {
	Handler         filtertype.FilterHandler
	Operator        string          `json:"operator"`
	FilterModelType FilterModelType `json:"filterModelType"`
}

type FilterModelType string

const (
	FilterModelT                FilterModelType = "model"
	FilterModelConditionT       FilterModelType = "model_condition"
	FilterModelCustomConditionT FilterModelType = "custom_conditions"
)

func (f *Filter) H() filtertype.FilterHandler {
	if f.Operator != "" {
		switch f.FilterModelType {
		case FilterModelCustomConditionT:
			f.Handler = &filtermodel.FilterModelCustomCondition{}
		case FilterModelConditionT:
			f.Handler = &filtermodel.FilterModelCondition{}
		default:
			f.Handler = &filtermodel.FilterModelCondition{}
		}
		return f.Handler
	}
	f.Handler = &filtermodel.FilterModel{}
	return f.Handler
}
