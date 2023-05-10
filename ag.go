package aggrid

import (
	"encoding/json"
	"fmt"
	"github.com/CarrotVegeta/aggrid/filter"
	"github.com/CarrotVegeta/aggrid/filtermodel"
	"github.com/CarrotVegeta/aggrid/filtertype"
	"github.com/CarrotVegeta/aggrid/utils"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

// AgGHandler aggrid查询
type AgGHandler interface {
	BuildFromSql() string
}

type Param struct {
	StartRow     int            `json:"startRow" form:"startRow"`
	EndRow       int            `json:"endRow" form:"endRow"`
	FilterModel  map[string]any `json:"filterModel" form:"filterModel"`
	RowGroupCols []RowGroupCol  `json:"rowGroupCols" form:"rowGroupCols[]"`
	GroupKeys    []string       `json:"groupKeys" form:"groupKeys[]"`
	SortModels   []SortModel    `json:"sortModel" form:"sortModel[]"`
}

func (p *Param) AppendSortModel(m *SortModel) *Param {
	p.SortModels = append(p.SortModels, *m)
	return p
}
func (p *Param) SetFilterModel(k string, filterModel filtermodel.FilterTypeHandler) {
	if p.FilterModel == nil {
		p.FilterModel = make(map[string]any)
	}
	p.FilterModel[k] = filterModel
}

type SortModel struct {
	Sort  string `json:"sort" form:"sort"`
	ColId string `json:"colId" form:"colId"`
}

func NewSortModel(sort string, colID string) *SortModel {
	return &SortModel{
		Sort:  sort,
		ColId: colID,
	}
}

type RowGroupCol struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName" `
	Field       string `json:"field"`
}
type AgGrid struct {
	Param            *Param
	Handler          AgGHandler
	selectField      map[string]string
	filterField      map[string]string
	havingField      map[string]string
	groupField       map[string]string
	orderField       map[string]string
	filterSqlService *filtertype.FilterTypeSqlService
	db               *gorm.DB
	qf               *utils.QueryFilter
	havingQf         *utils.QueryFilter
	sqlStr           string
	sqlLimitSqlStr   string
	customGroupSql   string
	sb               *SqlBuilder
	Error            error
}

// AddError add error to ag
func (a *AgGrid) AddError(err error) {
	if a.Error == nil {
		a.Error = err
	}
	a.Error = fmt.Errorf("%v; %w", a.Error, err)
}

func NewAgGHandler(model AgGHandler, param *Param, storageType ...string) *AgGrid {
	var st string
	if len(storageType) > 0 {
		st = storageType[0]
	}
	ag := &AgGrid{
		Param:            &Param{},
		qf:               &utils.QueryFilter{},
		havingQf:         &utils.QueryFilter{},
		Handler:          model,
		filterField:      make(map[string]string),
		havingField:      make(map[string]string),
		groupField:       make(map[string]string),
		orderField:       make(map[string]string),
		selectField:      make(map[string]string),
		filterSqlService: filtertype.NewFilterTypeSqlService(st),
		sb:               NewSqlBuilder(),
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
func (a *AgGrid) Group(groupSql string) *AgGrid {
	a.customGroupSql = groupSql
	return a
}
func (a *AgGrid) Count(count *int64) *AgGrid {
	sqlCountStr := a.SetSql().BuildCountSql().ToSqlString()
	err := a.db.Raw(sqlCountStr, a.sb.QueryFilter.Args...).Count(count).Error
	if err != nil {
		a.AddError(fmt.Errorf("ag grid count err:%v", err))
	}
	return a
}
func (a *AgGrid) SetSql() *SqlBuilder {
	a.sb.sqlStr.Reset()
	selectSql, err := a.BuildSelectSql()
	if err != nil {
		a.AddError(err)
		return nil
	}
	fromSql := a.BuildFromSql()
	groupSql, err := a.BuildGroupSql()
	if a.customGroupSql != "" {
		groupSql = a.customGroupSql
	}
	if err != nil {
		a.AddError(err)
		return nil
	}
	sortSql, err := a.BuildSortSql()
	a.sb.SetSelectSql(selectSql).SetFromSql(fromSql).SetGroupSql(groupSql).SetSortSql(sortSql)
	return a.sb
}

func (a *AgGrid) Where(qf *utils.QueryFilter) *AgGrid {
	a.sb.SetQueryFilter(qf)
	return a
}
func (a *AgGrid) Having(qf *utils.QueryFilter) *AgGrid {
	a.sb.SetHavingFilter(qf)
	return a
}
func (a *AgGrid) Find(data any) *AgGrid {
	a.sb = a.SetSql().BuildNoLimitSql()
	if a.Param.EndRow-a.Param.StartRow != 0 {
		a.sb.BuildAndLimitSql(a.Param.StartRow, a.Param.EndRow-a.Param.StartRow)
	}
	sqlStr := a.sb.ToSqlString()
	var m []map[string]any
	err := a.db.Raw(sqlStr, a.sb.QueryFilter.Args...).Find(&m).Error
	if err != nil {
		a.AddError(fmt.Errorf("%s:%v", RawSqlError, err))
		return a
	}
	//if IsMap(data) {
	//	data = m
	//	return a
	//}
	bs, _ := json.Marshal(m)
	if err := json.Unmarshal(bs, &data); err != nil {
		a.AddError(fmt.Errorf("ag unmarshal result to data err :%v", err))
		return a
	}
	return a
}

//	func IsMap(v any) bool {
//		rt := reflect.TypeOf(v)
//		if rt.Elem().Kind() == reflect.Slice || rt.Elem() == reflect.Array {
//			rt.
//		}
//	}
func IsStruct(v any) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Slice {
		if rv.Len() > 0 {
			elem := rv.Index(0)
			if elem.Kind() == reflect.Struct {
				return true
			}
		}
	}
	return false
}

func (a *AgGrid) parse() {
	fn := utils.GetStructTagField(a.Handler, "ag")
	a.parseFilterField(fn)
	a.parseGroupField(fn)
	a.parseSelectField(fn)
	a.parseHavingField(fn)
}
func (a *AgGrid) parseSelectField(fn utils.StructTag) {
	for k, v := range fn {
		s := a.getAgTagValue(k, "select")
		if s != "" {
			a.selectField[v.Get("json")] = s
		}
	}
}
func (a *AgGrid) parseFilterField(fn utils.StructTag) {
	for k, v := range fn {
		s := a.getAgTagValue(k, "filter")
		if s != "" {
			a.filterField[v.Get("json")] = s
		}
	}
}
func (a *AgGrid) parseHavingField(fn utils.StructTag) {
	for k, v := range fn {
		s := a.getAgTagValue(k, "having")
		if s != "" {
			a.havingField[v.Get("json")] = s
		}
	}
}
func (a *AgGrid) parseGroupField(fn utils.StructTag) {
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
func (a *AgGrid) getFilterField(k string) string {
	return a.filterField[k]
}
func (a *AgGrid) getHavingField(k string) string {
	return a.havingField[k]
}
func (a *AgGrid) getOrderField(k string) string {
	return a.orderField[k]
}
func (a *AgGrid) getSelectKeyFields() []*KeyField {
	fields := make([]*KeyField, 0, len(a.filterField))
	for k, v := range a.selectField {
		fields = append(fields, &KeyField{
			Key:   k,
			Field: v,
		})
	}
	return fields
}
func (a *AgGrid) getSelectField(k string) string {
	return a.selectField[k]
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
	return fmt.Sprintf("SELECT %s,COUNT(*) AS count", gn.Field), err
}
func (a *AgGrid) setOrderField(kf *KeyField) {
	if kf != nil {
		a.orderField[kf.Key] = kf.Field
	}
}
func (a *AgGrid) buildSelect() string {
	var selectSql string
	for _, v := range a.getSelectKeyFields() {
		a.setOrderField(v)
		if selectSql == "" {
			selectSql = v.Field
			continue
		}
		selectSql += "," + v.Field
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
		if gn.Field == "" {
			return "", nil
		}
		groupBySql := "GROUP BY " + gn.Field
		return groupBySql, nil
	}
	return "", nil
}
func (a *AgGrid) BuildQuerySql() (qf *utils.QueryFilter, err error) {
	err = a.buildGroupQuery(a.Param.RowGroupCols, a.Param.GroupKeys)
	if err != nil {
		return
	}
	err = a.parseFilterModel()
	if err != nil {
		return
	}
	if a.qf.Query == "" {
		return
	}
	a.qf.Query = "WHERE " + a.qf.Query
	return a.qf, nil
}
func (a *AgGrid) BuildHavingSql() (qf *utils.QueryFilter, err error) {
	err = a.buildGroupHavingQuery(a.Param.RowGroupCols, a.Param.GroupKeys)
	if err != nil {
		return
	}
	err = a.parseFilterModelToHavingSql()
	if err != nil {
		return
	}
	if a.havingQf.Query == "" {
		return
	}
	a.havingQf.Query = "HAVING " + a.havingQf.Query
	return a.havingQf, nil
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
	g.Field = a.getGroupField(groupName)
	if g.Field == "" {
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
func (a *AgGrid) BuildFromSql() string {
	return a.Handler.BuildFromSql()
}

// buildGroupQuery 生成组合查询where sql
func (a *AgGrid) buildGroupQuery(cols []RowGroupCol, keys []string) error {
	if len(cols) > 0 && len(keys) > 0 {
		for i, v := range keys {
			field := a.getFilterField(cols[i].Field)
			if field == "" {
				return fmt.Errorf("%s:%v", InvalidSqlField, cols[i].Field)
			}
			query := fmt.Sprintf("%s = ? ", field)
			a.qf.And(query, v)
		}
	}
	return nil
}

// buildGroupQuery 生成组合查询条件having sql
func (a *AgGrid) buildGroupHavingQuery(cols []RowGroupCol, keys []string) error {
	if len(cols) > 0 && len(keys) > 0 {
		for i, v := range keys {
			field := a.getHavingField(cols[i].Field)
			if field == "" {
				return fmt.Errorf("%s:%v", InvalidSqlField, field)
			}
			query := fmt.Sprintf("%s = ? ", field)
			a.havingQf.And(query, v)
		}
	}
	return nil
}

// ParseFilterModel 解析查询参数 并生成对应sql
func (a *AgGrid) parseFilterModel() error {
	for k, v := range a.Param.FilterModel {
		field := a.getFilterField(k)
		if field == "" {
			continue
		}
		bs, _ := json.Marshal(v)
		f := &filter.Filter{}
		if err := json.Unmarshal(bs, f); err != nil {
			return err
		}
		err := f.H().Parse(field, bs)
		if err != nil {
			return err
		}
		q, err := f.Handler.BuildQuery(a.filterSqlService)
		if err != nil {
			return err
		}
		a.qf.And(q.Query, q.Args...)
	}
	return nil
}

// ParseFilterModel 解析查询参数 并生成对应sql
func (a *AgGrid) parseFilterModelToHavingSql() error {
	for k, v := range a.Param.FilterModel {
		field := a.getHavingField(k)
		if field == "" {
			continue
		}
		bs, _ := json.Marshal(v)
		f := &filter.Filter{}
		if err := json.Unmarshal(bs, f); err != nil {
			return err
		}
		err := f.H().Parse(field, bs)
		if err != nil {
			return err
		}
		q, err := f.Handler.BuildQuery(a.filterSqlService)
		if err != nil {
			return err
		}
		a.havingQf.And(q.Query, q.Args...)
	}
	return nil
}
