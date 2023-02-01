package agtwo

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

// AgGHandler
// GetDB 获取db
// GetSqlField获取前端传来的字段所对应的sql 字段，如果没有则无效
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
}

func (sb *SqlBuilder) BuildSql() (string, error) {
	var sqlStr string
	sqlStr += sb.SelectSql + " "
	if sb.SelectSql == "" {
		sqlStr = "SELECT * "
	}
	if sb.FromSql == "" {
		return "", fmt.Errorf("from sql is null")
	}
	sqlStr += sb.FromSql + " "
	if sb.QuerySql != "" {
		sqlStr += sb.QuerySql + " "
	}
	if sb.GroupSql != "" {
		sqlStr += sb.GroupSql + " "
	}
	if sb.SortSql != "" {
		sqlStr += sb.SortSql + " "
	}
	return sqlStr, nil
}
func (ag *AgGrid) ExecSql(db *gorm.DB, sb *SqlBuilder) (data []map[string]any, count int64, err error) {
	sqlStr, err := sb.BuildSql()
	if err != nil {
		return nil, 0, err
	}
	sqlCountStr := fmt.Sprintf("SELECT COUNT(1) FROM (%s) AS a", sqlStr)
	if err := db.Raw(sqlCountStr, sb.Args...).Scan(&count).Error; err != nil {
		return nil, 0, err
	}
	if ag.Param.EndRow-ag.Param.StartRow != 0 {
		sqlStr += fmt.Sprintf("limit %d,%d", ag.Param.StartRow, ag.Param.EndRow-ag.Param.StartRow)
	}
	db = db.Raw(sqlStr, sb.Args...)
	err = db.Find(&data).Error
	if err != nil {
		return nil, 0, err
	}
	return
}
func (ag *AgGrid) buildGroupSelect() (string, error) {
	gn, err := ag.getGroupName(ag.Param.RowGroupCols, ag.Param.GroupKeys)
	return fmt.Sprintf("SELECT %s,COUNT(*) AS count", gn), err
}
func (ag *AgGrid) BuildSelect() string {
	var selectSql string
	for _, v := range ag.Handler.GetSelectField() {
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
func (ag *AgGrid) GetSelectSql() (string, error) {
	if len(ag.Param.RowGroupCols) > 0 && len(ag.Param.RowGroupCols) != len(ag.Param.GroupKeys) {
		return ag.buildGroupSelect()
	}
	return ag.BuildSelect(), nil
}

// BuildGroupSql 如果分组参数大于0 并且 分组参数不等于key值，则拼接groupbysql
func (ag *AgGrid) BuildGroupSql() (string, error) {
	if len(ag.Param.RowGroupCols) > 0 && len(ag.Param.RowGroupCols) != len(ag.Param.GroupKeys) {
		gn, err := ag.getGroupName(ag.Param.RowGroupCols, ag.Param.GroupKeys)
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
func (ag *AgGrid) BuildQuerySql() (query string, args []any, err error) {
	err = ag.buildGroupQuery(ag.Param.RowGroupCols, ag.Param.GroupKeys)
	if err != nil {
		return "", nil, nil
	}
	err = ag.parseFilterModel()
	if err != nil {
		return "", nil, err
	}
	if ag.qf.Query == "" {
		return "", nil, nil
	}
	query = "WHERE " + ag.qf.Query
	return query, ag.qf.Args, nil
}

// getGroupName 获取分组条件
func (ag *AgGrid) getGroupName(cols []RowGroupCol, keys []string) (string, error) {
	if len(cols) == 0 {
		return "", nil
	}
	groupName := cols[0].Field
	field := ag.Handler.GetSqlField(cols[0].Field)
	if field == "" {
		return "", fmt.Errorf("invalid group field :%v", groupName)
	}
	if len(keys) > 0 {
		if len(cols) == len(keys) {
			groupName = cols[len(cols)-1].Field
		} else {
			groupName = cols[len(keys)].Field
		}
		field = ag.Handler.GetSqlField(groupName)
		if field == "" {
			return "", fmt.Errorf("invalid group field2 :%v", groupName)
		}
	}
	return field, nil
}

// BuildSortSql 生成排序sql
func (ag *AgGrid) BuildSortSql(sortModels []SortModel) (string, error) {
	//groupName, err := ag.getGroupName(ag.Param.RowGroupCols, ag.Param.GroupKeys)
	//if err != nil {
	//	return "", nil
	//}
	var sortStr string
	if len(sortModels) > 0 {
		if len(ag.Param.RowGroupCols) > 0 && len(ag.Param.RowGroupCols) != len(ag.Param.GroupKeys) {

			var sm SortModel
			if len(ag.Param.GroupKeys) == 0 {
				sm = sortModels[0]
			} else {
				sm = sortModels[len(ag.Param.GroupKeys)]
			}
			gn, err := ag.getGroupName(ag.Param.RowGroupCols, ag.Param.GroupKeys)
			if err != nil {
				return "", err
			}
			if gn != sm.ColId {
				return "", err
			}
			ss, err := ag.buildsortstr(sm.ColId, sm.Sort)
			if err != nil {
				return "", err
			}
			sortStr = ss
		} else {
			for _, v := range sortModels {
				ss, err := ag.buildsortstr(v.ColId, v.Sort)
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
func (ag *AgGrid) buildsortstr(colid, sort string) (string, error) {
	k := ag.Handler.GetSqlField(colid)
	if k == "" {
		return "", fmt.Errorf("invalid colid %v", colid)
	}
	//if groupName != "" {
	//	if k != groupName {
	//		continue
	//	}
	//}
	if !validSort(sort) {
		return "", fmt.Errorf("invalid sort : %v", k)
	}
	return fmt.Sprintf("%s %s", k, sort), nil
}

// buildGroupQuery 生成组合查询sql
func (ag *AgGrid) buildGroupQuery(cols []RowGroupCol, keys []string) error {
	if len(cols) > 0 && len(keys) > 0 {
		for i, v := range keys {
			field := ag.Handler.GetSqlField(cols[i].Field)
			if field == "" {
				return fmt.Errorf("invalid field:%v", field)
			}
			query := fmt.Sprintf("%s = ? ", field)
			ag.qf.And(query, v)
		}
	}
	return nil
}

// ParseFilterModel 解析查询参数 并生成对应sql
func (ag *AgGrid) parseFilterModel() error {
	for k, v := range ag.Param.FilterModel {
		field := ag.Handler.GetSqlField(k)
		if field == "" {
			return fmt.Errorf("invalid model field : %v", k)
		}
		bs, _ := json.Marshal(v)
		f := &Filter{}
		if err := json.Unmarshal(bs, f); err != nil {
			return err
		}
		err := f.H().Parse(k, bs)
		if err != nil {
			return err
		}
		q, err := f.Handler.BuildQuery()
		if err != nil {
			return err
		}
		ag.qf.And(q.Query, q.Args...)
	}
	return nil
}
