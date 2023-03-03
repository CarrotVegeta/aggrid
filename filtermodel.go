package agtwo

import (
	"encoding/json"
	"fmt"
)

// FilterModel 单个查询条件
type FilterModel struct {
	FilterType string `json:"filterType" form:"filterType"`
	Key        string //查询的字段名
	Type       string `json:"type"`
	Param      []byte //查询的参数
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
func (f *FilterModel) BuildQuery() (*QueryFilter, error) {
	handler, err := NewFilterTypeHandler(f.FilterType, f.Param)
	if err != nil {
		return nil, err
	}
	filter, err := handler.GetFilter()
	if err != nil {
		return nil, err
	}
	h, err := NewFilterSqlHandler(f.FilterType)
	if err != nil {
		return nil, err
	}
	return h.BuildSql(f.Key, filter, handler.GetType())
}

var FilterTypeHandlerM = map[string]FilterTypeHandler{
	Text:   &FilterTextModel{},
	Number: &FilterNumberModel{},
	Date:   &FilterDateModel{},
	Array:  &FilterArrayModel{},
}

func getFilterTypeHandler(filterType string) (FilterTypeHandler, error) {
	if _, ok := FilterTypeHandlerM[filterType]; !ok {
		return nil, fmt.Errorf("invalid filtertype : %v", filterType)
	}
	return FilterTypeHandlerM[filterType], nil
}

type FilterType struct {
}

func NewFilterTypeHandler(filterType string, c []byte) (FilterTypeHandler, error) {
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
	GetType() string
}
