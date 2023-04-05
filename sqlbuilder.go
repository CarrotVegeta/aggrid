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
