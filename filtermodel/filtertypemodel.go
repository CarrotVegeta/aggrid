package filtermodel

import (
	"encoding/json"
	"github.com/CarrotVegeta/aggrid/constant"
)

type FilterTextModel struct {
	FilterType constant.FilterType   `json:"filterType" form:"filterType"`
	Type       constant.OperatorType `json:"type" form:"type"`
	Filter     string                `json:"filter" form:"filter"`
	Key        string                `json:"key"`
}

func NewFilterTextModel(t constant.OperatorType, key, filter string) *FilterTextModel {
	return &FilterTextModel{
		FilterType: constant.Text,
		Type:       t,
		Filter:     filter,
		Key:        key,
	}
}
func (ftm *FilterTextModel) New() FilterTypeHandler {
	return &FilterTextModel{}
}
func (ftm *FilterTextModel) Parse(c []byte) error {
	if err := json.Unmarshal(c, ftm); err != nil {
		return err
	}
	return nil
}
func (ftm *FilterTextModel) GetFilter() (any, error) {
	return ftm.Filter, nil
}
func (ftm *FilterTextModel) GetType() constant.OperatorType {
	return ftm.Type
}

type FilterDateModel struct {
	FilterType constant.FilterType   `json:"FilterType" form:"filterType"`
	DateFrom   string                `json:"dateFrom" form:"dateFrom"`
	DateTo     string                `json:"dateTo" form:"dateTo"`
	Type       constant.OperatorType `json:"type" form:"type"`
	Key        string                `json:"key"`
}

func (fdm *FilterDateModel) New() FilterTypeHandler {
	return &FilterDateModel{}
}
func (fdm *FilterDateModel) Parse(c []byte) error {
	if err := json.Unmarshal(c, fdm); err != nil {
		return err
	}
	return nil
}
func (fdm *FilterDateModel) GetFilter() (any, error) {
	return fdm.ParseFilter(), nil
}
func (fdm *FilterDateModel) GetType() constant.OperatorType {
	return fdm.Type
}
func (fdm *FilterDateModel) ParseFilter() any {
	if fdm.Type == constant.InRange {
		return []string{fdm.DateFrom, fdm.DateTo}
	}
	return fdm.DateFrom
}

type FilterNumberModel struct {
	FilterType constant.FilterType   `json:"filterType" form:"filterType"`
	Type       constant.OperatorType `json:"type" form:"type"`
	Filter     int64                 `json:"filter" form:"filter"`
	FilterTo   int64                 `json:"filterTo" form:"filterTo"`
	Key        string                `json:"key"`
}

func NewFilterNumberModel(t constant.OperatorType, filter, filterTo int64, key string) *FilterNumberModel {
	return &FilterNumberModel{
		FilterType: constant.Number,
		Type:       t,
		Filter:     filter,
		FilterTo:   filterTo,
		Key:        key,
	}
}
func (fnm *FilterNumberModel) New() FilterTypeHandler {
	return &FilterNumberModel{}
}
func (fnm *FilterNumberModel) Parse(c []byte) error {
	if err := json.Unmarshal(c, fnm); err != nil {
		return err
	}
	return nil
}
func (fnm *FilterNumberModel) GetFilter() (any, error) {
	return fnm.parseFilter(), nil
}
func (fnm *FilterNumberModel) GetType() constant.OperatorType {
	return fnm.Type
}
func (fnm *FilterNumberModel) parseFilter() any {
	if fnm.Type == constant.InRange {
		return []int64{fnm.Filter, fnm.FilterTo}
	}
	return fnm.Filter
}

type FilterArrayModel struct {
	FilterType constant.FilterType   `json:"filterType" form:"filterType"`
	Type       constant.OperatorType `json:"type" form:"type"`
	Filter     []any                 `json:"filter" form:"filter"`
	Key        string                `json:"key"`
}

func NewFilterArrayModel(t constant.OperatorType, key string, filter []any) *FilterArrayModel {
	return &FilterArrayModel{
		FilterType: constant.Array,
		Type:       t,
		Filter:     filter,
		Key:        key,
	}
}
func (fnm *FilterArrayModel) New() FilterTypeHandler {
	return &FilterArrayModel{}
}
func (fnm *FilterArrayModel) Parse(c []byte) error {
	if err := json.Unmarshal(c, fnm); err != nil {
		return err
	}
	return nil
}
func (fnm *FilterArrayModel) GetFilter() (any, error) {
	return fnm.parseFilter(), nil
}
func (fnm *FilterArrayModel) GetType() constant.OperatorType {
	return fnm.Type
}
func (fnm *FilterArrayModel) parseFilter() any {
	return fnm.Filter
}

type FilterSetModel struct {
	FilterType constant.FilterType `json:"filterType" form:"filterType"`
	Values     []any               `json:"values" form:"values"`
}

func (fnm *FilterSetModel) New() FilterTypeHandler {
	return &FilterSetModel{}
}
func (fnm *FilterSetModel) Parse(c []byte) error {
	if err := json.Unmarshal(c, fnm); err != nil {
		return err
	}
	return nil
}
func (fnm *FilterSetModel) GetFilter() (any, error) {
	return fnm.parseFilter(), nil
}
func (fnm *FilterSetModel) GetType() constant.OperatorType {
	return constant.InRange
}
func (fnm *FilterSetModel) parseFilter() any {
	return fnm.Values
}
