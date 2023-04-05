package agtwo

import "encoding/json"

type FilterTextModel struct {
	FilterType string `json:"filterType" form:"filterType"`
	Type       string `json:"type" form:"type"`
	Filter     string `json:"filter" form:"filter"`
	Key        string `json:"key"`
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
func (ftm *FilterTextModel) GetType() string {
	return ftm.Type
}

type FilterDateModel struct {
	FilterType string `json:"FilterType" form:"filterType"`
	DateFrom   string `json:"dateFrom" form:"dateFrom"`
	DateTo     string `json:"dateTo" form:"dateTo"`
	Type       string `json:"type" form:"type"`
	Key        string `json:"key"`
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
func (fdm *FilterDateModel) GetType() string {
	return fdm.Type
}
func (fdm *FilterDateModel) ParseFilter() any {
	if fdm.Type == InRange {
		return []string{fdm.DateFrom, fdm.DateTo}
	}
	return fdm.DateFrom
}

type FilterNumberModel struct {
	FilterType string `json:"filterType" form:"filterType"`
	Type       string `json:"type" form:"type"`
	Filter     int64  `json:"filter" form:"filter"`
	FilterTo   int64  `json:"filterTo" form:"filterTo"`
	Key        string `json:"key"`
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
func (fnm *FilterNumberModel) GetType() string {
	return fnm.Type
}
func (fnm *FilterNumberModel) parseFilter() any {
	if fnm.Type == InRange {
		return []int64{fnm.Filter, fnm.FilterTo}
	}
	return fnm.Filter
}

type FilterArrayModel struct {
	FilterType string `json:"filterType" form:"filterType"`
	Type       string `json:"type" form:"type"`
	Filter     []any  `json:"filter" form:"filter"`
	Key        string `json:"key"`
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
func (fnm *FilterArrayModel) GetType() string {
	return fnm.Type
}
func (fnm *FilterArrayModel) parseFilter() any {
	return fnm.Filter
}
