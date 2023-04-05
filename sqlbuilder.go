package agtwo

import "strings"

type SqlBuilder struct {
	SelectSql string
	GroupSql  string
	QuerySql  string
	SortSql   string
	FromSql   string
	Args      []any
	sqlStr    strings.Builder
}

func NewSqlBuilder() *SqlBuilder {
	return &SqlBuilder{}
}
func (sb *SqlBuilder) SetSelectSql(selectSql string) *SqlBuilder {
	sb.SelectSql = selectSql
	return sb
}
func (sb *SqlBuilder) SetGroupSql(groupSql string) *SqlBuilder {
	sb.GroupSql = groupSql
	return sb
}
func (sb *SqlBuilder) SetQuerySql(querySql string) *SqlBuilder {
	sb.QuerySql = querySql
	return sb
}
func (sb *SqlBuilder) SetSortSql(sortSql string) *SqlBuilder {
	sb.SortSql = sortSql
	return sb
}
func (sb *SqlBuilder) SetFromSql(fromSql string) *SqlBuilder {
	sb.FromSql = fromSql
	return sb
}
func (sb *SqlBuilder) WriteSqlStr(s string) {
	if s == "" {
		return
	}
	sb.sqlStr.WriteString(s + " ")
}
func (sb *SqlBuilder) BuildSelectSql() *SqlBuilder {
	sb.WriteSqlStr(sb.SelectSql)
	if sb.SelectSql == "" {
		sb.WriteSqlStr("SELECT *")
	}
	return sb
}
func (sb *SqlBuilder) BuildFromSql() *SqlBuilder {
	sb.WriteSqlStr(sb.FromSql)
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
func (sb *SqlBuilder) ToSqlString() string {
	sqlStr := sb.SqlString()
	sb.sqlStr.Reset()
	return sqlStr
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
