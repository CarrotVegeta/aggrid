package agtwo

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

// AgGHandler
// GetSqlField获取前端传来的字段所对应的sql 字段，如果没有则无效
// GetSelectField 获取需要查询的字段
type AgGHandler interface {
}

type Param struct {
	StartRow     int            `json:"startRow" form:"startRow"`
	EndRow       int            `json:"endRow" form:"endRow"`
	FilterModel  map[string]any `json:"filterModel" form:"filterModel"`
	RowGroupCols []RowGroupCol  `json:"rowGroupCols" form:"rowGroupCols[]"`
	GroupKeys    []string       `json:"groupKeys" form:"groupKeys[]"`
	SortModels   []SortModel    `json:"sortModel" form:"sortModel[]"`
}
type SortModel struct {
	Sort  string `json:"sort" form:"sort"`
	ColId string `json:"colId" form:"colId"`
}
type RowGroupCol struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName" `
	Field       string `json:"field"`
}
type AgGrid struct {
	Param          *Param
	Handler        AgGHandler
	selectField    map[string]string
	groupField     map[string]string
	orderField     map[string]string
	db             *gorm.DB
	qf             *QueryFilter
	sqlStr         string
	sqlCountStr    string
	sqlLimitSqlStr string
	Error          error
}

// AddError add error to ag
func (a *AgGrid) AddError(err error) {
	if a.Error == nil {
		a.Error = err
	}
	a.Error = fmt.Errorf("%v; %w", a.Error, err)
}

func NewAgGHandler(model AgGHandler, param *Param) *AgGrid {
	ag := &AgGrid{
		Param:       &Param{},
		qf:          &QueryFilter{},
		Handler:     model,
		selectField: make(map[string]string),
		groupField:  make(map[string]string),
		orderField:  make(map[string]string),
	}
	if param != nil {
		ag.Param = param
	}
	ag.parse()
	return ag
}

func (a *AgGrid) Use(db *gorm.DB) *AgGrid {
	a.db = db
	return a
}
func (a *AgGrid) Raw(sb *SqlBuilder) *AgGrid {
	a.sqlStr = sb.BuildNoLimitSql().ToSqlString()
	a.sqlCountStr = sb.BuildCountSql().ToSqlString()
	if a.Param.EndRow-a.Param.StartRow != 0 {
		a.sqlLimitSqlStr = sb.BuildAndLimitSql(a.Param.StartRow, a.Param.EndRow-a.Param.StartRow).ToSqlString()
	}
	return a
}

// todo 使用链式操作
func (a *AgGrid) Count(sb *SqlBuilder, count int64) *AgGrid {
	err := a.db.Raw(a.sqlCountStr, sb.Args...).Count(&count).Error
	if err != nil {
		a.AddError(fmt.Errorf("ag grid count err:%v", err))
	}
	return a
}
func (a *AgGrid) Find(sb *SqlBuilder, data any) *AgGrid {
	if a.sqlLimitSqlStr != "" {
		a.sqlStr = a.sqlLimitSqlStr
	}
	db := a.db.Raw(a.sqlStr, sb.Args...)
	err := db.Find(&data).Error
	if err != nil {
		a.AddError(fmt.Errorf("%s:%v", RawSqlError, err))
	}
	return a
}

func (a *AgGrid) parse() {
	fn := GetStructTagField(a.Handler, "ag")
	a.parseSelectField(fn)
	a.parseGroupField(fn)
}
func (a *AgGrid) parseSelectField(fn StructTag) {
	for k, v := range fn {
		s := a.getAgTagValue(k, "select")
		if s != "" {
			a.selectField[v.Get("json")] = s
		}
	}
}
func (a *AgGrid) parseGroupField(fn StructTag) {
	for k, v := range fn {
		s := a.getAgTagValue(k, "group")
		if s == "" {
			s = a.getAgTagValue(k, "select")
		}
		if s != "" {
			a.groupField[v.Get("json")] = s
		}
	}
}

func (a *AgGrid) getAgTagValue(agTag, tag string) string {
	tags := strings.Split(agTag, ";")
	for _, v := range tags {
		ts := strings.Split(v, ":")
		if strings.Trim(ts[0], " ") == tag {
			if len(ts) > 1 {
				return ts[1]
			}
		}
	}
	return ""
}
func (a *AgGrid) getSelectField(k string) string {
	return a.selectField[k]
}
func (a *AgGrid) getOrderField(k string) string {
	return a.orderField[k]
}
func (a *AgGrid) getSelectKeyFields() []*KeyField {
	fields := make([]*KeyField, 0, len(a.selectField))
	for k, v := range a.selectField {
		fields = append(fields, &KeyField{
			Key:         k,
			SelectField: v,
		})
	}
	return fields
}
func (a *AgGrid) getGroupField(k string) string {
	return a.groupField[k]
}
func (a *AgGrid) buildGroupSelect() (string, error) {
	gn, err := a.getGroupName(a.Param.RowGroupCols, a.Param.GroupKeys)
	if gn == nil {
		return "", nil
	}
	a.setOrderField(gn)
	return fmt.Sprintf("SELECT %s,COUNT(*) AS count", gn.SelectField), err
}
func (a *AgGrid) setOrderField(kf *KeyField) {
	if kf != nil {
		a.orderField[kf.Key] = kf.SelectField
	}
}
func (a *AgGrid) buildSelect() string {
	var selectSql string
	for _, v := range a.getSelectKeyFields() {
		a.setOrderField(v)
		if selectSql == "" {
			selectSql = v.SelectField
			continue
		}
		selectSql += "," + v.SelectField
	}
	if selectSql == "" {
		selectSql = "*"
	}
	return "SELECT " + selectSql
}
func (a *AgGrid) BuildSelectSql() (string, error) {
	if len(a.Param.RowGroupCols) > 0 && len(a.Param.RowGroupCols) != len(a.Param.GroupKeys) {
		return a.buildGroupSelect()
	}
	return a.buildSelect(), nil
}

// BuildGroupSql 如果分组参数大于0 并且 分组参数不等于key值，则拼接groupBySql
func (a *AgGrid) BuildGroupSql() (string, error) {
	if len(a.Param.RowGroupCols) > 0 && len(a.Param.RowGroupCols) != len(a.Param.GroupKeys) {
		gn, err := a.getGroupName(a.Param.RowGroupCols, a.Param.GroupKeys)
		if err != nil {
			return "", nil
		}
		if gn.SelectField == "" {
			return "", nil
		}
		groupBySql := "GROUP BY " + gn.SelectField
		return groupBySql, nil
	}
	return "", nil
}
func (a *AgGrid) BuildQuerySql() (query string, args []any, err error) {
	err = a.buildGroupQuery(a.Param.RowGroupCols, a.Param.GroupKeys)
	if err != nil {
		return "", nil, nil
	}
	err = a.parseFilterModel()
	if err != nil {
		return "", nil, err
	}
	if a.qf.Query == "" {
		return "", nil, nil
	}
	query = "WHERE " + a.qf.Query
	return query, a.qf.Args, nil
}

// getGroupName 获取分组条件
func (a *AgGrid) getGroupName(cols []RowGroupCol, keys []string) (*KeyField, error) {
	g := &KeyField{}
	if len(cols) == 0 {
		return nil, nil
	}
	groupName := cols[0].Field
	if len(keys) > 0 {
		if len(cols) == len(keys) {
			groupName = cols[len(cols)-1].Field
		} else {
			groupName = cols[len(keys)].Field
		}
	}
	g.Key = groupName
	g.SelectField = a.getGroupField(groupName)
	if g.SelectField == "" {
		return nil, fmt.Errorf("%s:%v", InvalidGroupField, groupName)
	}
	return g, nil
}

// BuildSortSql 生成排序sql
func (a *AgGrid) BuildSortSql() (string, error) {
	var sortStr string
	for _, v := range a.Param.SortModels {
		orderField := a.getOrderField(v.ColId)
		if orderField == "" {
			continue
		}
		ss := fmt.Sprintf("%s %s", orderField, v.Sort)
		if sortStr == "" {
			sortStr = ss
			continue
		}
		sortStr += "," + ss
	}
	if sortStr == "" {
		return sortStr, nil
	}
	sortStr = "ORDER BY " + sortStr
	return sortStr, nil
}

// buildGroupQuery 生成组合查询sql
func (a *AgGrid) buildGroupQuery(cols []RowGroupCol, keys []string) error {
	if len(cols) > 0 && len(keys) > 0 {
		for i, v := range keys {
			field := a.getSelectField(cols[i].Field)
			if field == "" {
				return fmt.Errorf("%s:%v", InvalidSqlField, field)
			}
			query := fmt.Sprintf("%s = ? ", field)
			a.qf.And(query, v)
		}
	}
	return nil
}

// ParseFilterModel 解析查询参数 并生成对应sql
func (a *AgGrid) parseFilterModel() error {
	for k, v := range a.Param.FilterModel {
		field := a.getSelectField(k)
		if field == "" {
			return fmt.Errorf("%s : %v", InvalidSqlField, k)
		}
		bs, _ := json.Marshal(v)
		f := &Filter{}
		if err := json.Unmarshal(bs, f); err != nil {
			return err
		}
		err := f.H().Parse(field, bs)
		if err != nil {
			return err
		}
		q, err := f.Handler.BuildQuery()
		if err != nil {
			return err
		}
		a.qf.And(q.Query, q.Args...)
	}
	return nil
}
