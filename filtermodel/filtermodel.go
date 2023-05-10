package filtermodel

import (
	"encoding/json"
	"fmt"
	"github.com/CarrotVegeta/aggrid/constant"
	"github.com/CarrotVegeta/aggrid/filtertype"
	"github.com/CarrotVegeta/aggrid/utils"
)

// FilterModel 单个查询条件
type FilterModel struct {
	FilterType constant.FilterType `json:"filterType" form:"filterType"`
	Key        string              //查询的字段名
	Type       string              `json:"type"`
	Param      []byte              //查询的参数
}

// Parse 单个查询条件的参数解析
func (f *FilterModel) Parse(k string, c []byte) error {
	f.Key = k
	if err := json.Unmarshal(c, f); err != nil {
		return fmt.Errorf("parse filter model failed:%v", err.Error())
	}
	f.Param = c
	return nil
}

// BuildQuery 单个查询条件的sql生成
func (f *FilterModel) BuildQuery(service *filtertype.FilterTypeSqlService) (*utils.QueryFilter, error) {
	handler, err := NewFilterTypeHandler(f.FilterType, f.Param)
	if err != nil {
		return nil, err
	}
	filter, err := handler.GetFilter()
	if err != nil {
		return nil, err
	}
	h, err := service.NewFilterSqlHandler(f.FilterType)
	if err != nil {
		return nil, err
	}
	return h.BuildSql(f.Key, filter, handler.GetType())
}

func init() {
	registerFilterType(constant.Text, &FilterTextModel{})
	registerFilterType(constant.Number, &FilterNumberModel{})
	registerFilterType(constant.Date, &FilterDateModel{})
	registerFilterType(constant.Array, &FilterArrayModel{})
	registerFilterType(constant.Set, &FilterSetModel{})
}

var FilterTypeHandlerM = make(map[constant.FilterType]FilterTypeHandler)

func RegisterFilterType(filterType constant.FilterType, h FilterTypeHandler) {
	registerFilterType(filterType, h)
}
func registerFilterType(filterType constant.FilterType, h FilterTypeHandler) {
	FilterTypeHandlerM[filterType] = h
}
func getFilterTypeHandler(filterType constant.FilterType) (FilterTypeHandler, error) {
	if _, ok := FilterTypeHandlerM[filterType]; !ok {
		return nil, fmt.Errorf("invalid filtertype : %v", filterType)
	}
	return FilterTypeHandlerM[filterType].New(), nil
}

func NewFilterTypeHandler(filterType constant.FilterType, c []byte) (FilterTypeHandler, error) {
	handler, err := getFilterTypeHandler(filterType)
	if err != nil {
		return nil, err
	}
	if err := handler.Parse(c); err != nil {
		return nil, err
	}
	return handler, nil
}

type FilterTypeHandler interface {
	Parse(c []byte) error
	GetFilter() (any, error)
	GetType() constant.OperatorType
	New() FilterTypeHandler
}
