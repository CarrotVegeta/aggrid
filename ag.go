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
	GetSqlField(k string) string
	GetSelectField() []string
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
	Param   *Param
	Handler AgGHandler
	db      *gorm.DB
	qf      *QueryFilter
	sortStr string
}

func NewAgGHandler(model AgGHandler, param *Param) *AgGrid {
	ag := &AgGrid{
		Param:   &Param{},
		qf:      &QueryFilter{},
		Handler: model,
	}
	if param != nil {
		ag.Param = param
	}
	return ag
}

type SqlBuilder struct {
	SelectSql string
	GroupSql  string
	QuerySql  string
	SortSql   string
	FromSql   string
	Args      []any
	sqlStr    strings.Builder
	Err       error
}

func (sb *SqlBuilder) WriteSqlStr(s string) {
	sb.sqlStr.WriteString(s + " ")
}
func (sb *SqlBuilder) BuildSelectSql() *SqlBuilder {
	sb.WriteSqlStr(sb.SelectSql)
	if sb.SelectSql == "" {
		sb.WriteSqlStr("SELECT * ")
	}
	return sb
}
func (sb *SqlBuilder) BuildFromSql() *SqlBuilder {
	if sb.FromSql == "" {
		sb.Err = InvalidFromSql
	}
	sb.WriteSqlStr(sb.FromSql + " ")
	return sb
}
func (sb *SqlBuilder) BuildQuerySql() *SqlBuilder {
	sb.WriteSqlStr(sb.QuerySql)
	return sb
}
func (sb *SqlBuilder) BuildGroupSql() *SqlBuilder {
	sb.WriteSqlStr(sb.GroupSql)
	return sb
}
func (sb *SqlBuilder) BuildSortSql() *SqlBuilder {
	sb.WriteSqlStr(sb.SortSql)
	return sb
}
func (sb *SqlBuilder) SqlString() string {
	return sb.sqlStr.String()
}
func (sb *SqlBuilder) BuildNoLimitSql() *SqlBuilder {
	sb.BuildSelectSql().BuildFromSql().BuildQuerySql().BuildGroupSql().BuildSortSql()
	return sb

}
func (sb *SqlBuilder) ToSqlString() (string, error) {
	if sb.Err != nil {
		return "", sb.Err
	}
	sqlStr := sb.SqlString()
	sb.sqlStr.Reset()
	return sqlStr, nil
}
func (sb *SqlBuilder) BuildAndLimitSql(offset, pageSize int) *SqlBuilder {
	sb.sqlStr.Reset()
	sb.BuildNoLimitSql()
	sb.BuildSelectSql().BuildFromSql().BuildQuerySql().BuildGroupSql().BuildSortSql().BuildLimitSql(offset, pageSize)
	return sb
}
func (sb *SqlBuilder) BuildLimitSql(offset, pageSize int) *SqlBuilder {
	sb.WriteSqlStr(BuildLimitSql(offset, pageSize))
	return sb
}
func (sb *SqlBuilder) BuildCountSql() *SqlBuilder {
	sqlStr := BuildCountSql(sb.SqlString())
	sb.sqlStr.Reset()
	sb.WriteSqlStr(sqlStr)
	return sb
}
func (a *AgGrid) Use(db *gorm.DB) *AgGrid {
	a.db = db
	return a
}
func (a *AgGrid) ExecSql(sb *SqlBuilder) (data []map[string]any, count int64, err error) {
	var (
		sqlStr, sqlCountStr, sqlLimitSqlStr string
	)
	sqlStr, err = sb.BuildNoLimitSql().ToSqlString()
	sqlCountStr, err = sb.BuildCountSql().ToSqlString()
	if a.Param.EndRow-a.Param.StartRow != 0 {
		sqlLimitSqlStr, err = sb.BuildAndLimitSql(a.Param.StartRow, a.Param.EndRow-a.Param.StartRow).ToSqlString()
	}
	if err != nil {
		return nil, 0, err
	}
	if err := a.db.Raw(sqlCountStr, sb.Args...).Scan(&count).Error; err != nil {
		return nil, 0, err
	}
	if sqlLimitSqlStr != "" {
		sqlStr = sqlLimitSqlStr
	}
	db := a.db.Raw(sqlStr, sb.Args...)
	err = db.Find(&data).Error
	if err != nil {
		return nil, 0, fmt.Errorf("%s:%v", RawSqlError, err)
	}
	return
}
func (a *AgGrid) buildGroupSelect() (string, error) {
	gn, err := a.getGroupName(a.Param.RowGroupCols, a.Param.GroupKeys)
	return fmt.Sprintf("SELECT %s,COUNT(*) AS count", gn), err
}
func (a *AgGrid) BuildSelect() string {
	var selectSql string
	for _, v := range a.Handler.GetSelectField() {
		if selectSql == "" {
			selectSql = v
			continue
		}
		selectSql += "," + v
	}
	if selectSql == "" {
		selectSql = "*"
	}
	return "SELECT " + selectSql
}
func (a *AgGrid) GetSelectSql() (string, error) {
	if len(a.Param.RowGroupCols) > 0 && len(a.Param.RowGroupCols) != len(a.Param.GroupKeys) {
		return a.buildGroupSelect()
	}
	return a.BuildSelect(), nil
}

// BuildGroupSql 如果分组参数大于0 并且 分组参数不等于key值，则拼接groupbysql
func (a *AgGrid) BuildGroupSql() (string, error) {
	if len(a.Param.RowGroupCols) > 0 && len(a.Param.RowGroupCols) != len(a.Param.GroupKeys) {
		gn, err := a.getGroupName(a.Param.RowGroupCols, a.Param.GroupKeys)
		if err != nil {
			return "", nil
		}
		if gn == "" {
			return "", nil
		}
		groupBySql := "GROUP BY " + gn
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
func (a *AgGrid) getGroupName(cols []RowGroupCol, keys []string) (string, error) {
	if len(cols) == 0 {
		return "", nil
	}
	groupName := cols[0].Field
	if len(keys) > 0 {
		if len(cols) == len(keys) {
			groupName = cols[len(cols)-1].Field
		} else {
			groupName = cols[len(keys)].Field
		}
	}
	field := a.Handler.GetSqlField(groupName)
	if field == "" {
		return "", fmt.Errorf("%s:%v", InvalidGroupField, groupName)
	}
	return field, nil
}

// BuildSortSql 生成排序sql
func (a *AgGrid) BuildSortSql(sortModels []SortModel) (string, error) {
	//groupName, err := a.getGroupName(a.Param.RowGroupCols, a.Param.GroupKeys)
	//if err != nil {
	//	return "", nil
	//}
	var sortStr string
	if len(sortModels) > 0 {
		if len(a.Param.RowGroupCols) > 0 && len(a.Param.RowGroupCols) != len(a.Param.GroupKeys) {

			var sm SortModel
			if len(a.Param.GroupKeys) == 0 {
				sm = sortModels[0]
			} else {
				if len(a.Param.GroupKeys) != len(sortModels) {
					sm = sortModels[len(a.Param.GroupKeys)]
				}
			}
			gn, err := a.getGroupName(a.Param.RowGroupCols, a.Param.GroupKeys)
			if err != nil {
				return "", err
			}
			if gn != sm.ColId {
				return "", err
			}
			ss, err := a.buildSortStr(sm.ColId, sm.Sort)
			if err != nil {
				return "", err
			}
			sortStr = ss
		} else {
			for _, v := range sortModels {
				ss, err := a.buildSortStr(v.ColId, v.Sort)
				if err != nil {
					return "", err
				}
				if sortStr == "" {
					sortStr = ss
					continue
				}
				sortStr += "," + ss
			}
		}
	}
	if sortStr == "" {
		return sortStr, nil
	}
	sortStr = "ORDER BY " + sortStr
	return sortStr, nil
}
func (a *AgGrid) buildSortStr(colId, sort string) (string, error) {
	k := a.Handler.GetSqlField(colId)
	if k == "" {
		return "", fmt.Errorf("%s:%v", InvalidSqlField, colId)
	}
	if !validSort(sort) {
		return "", fmt.Errorf("%s : %v", InvalidSortField, k)
	}
	return fmt.Sprintf("%s %s", k, sort), nil
}

// buildGroupQuery 生成组合查询sql
func (a *AgGrid) buildGroupQuery(cols []RowGroupCol, keys []string) error {
	if len(cols) > 0 && len(keys) > 0 {
		for i, v := range keys {
			field := a.Handler.GetSqlField(cols[i].Field)
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
		field := a.Handler.GetSqlField(k)
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
