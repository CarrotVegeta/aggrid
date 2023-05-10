package aggrid

import (
	"github.com/CarrotVegeta/aggrid/utils"
	"strings"
)

type SqlBuilder struct {
	SelectSql         string
	GroupSql          string
	QueryFilter       *utils.QueryFilter
	HavingQueryFilter *utils.QueryFilter
	SortSql           string
	FromSql           string
	sqlStr            strings.Builder
}

func NewSqlBuilder() *SqlBuilder {
	return &SqlBuilder{QueryFilter: &utils.QueryFilter{}, HavingQueryFilter: &utils.QueryFilter{}}
}
func (sb *SqlBuilder) SetSelectSql(selectSql string) *SqlBuilder {
	sb.SelectSql = selectSql
	return sb
}
func (sb *SqlBuilder) SetGroupSql(groupSql string) *SqlBuilder {
	sb.GroupSql = groupSql
	return sb
}
func (sb *SqlBuilder) SetQueryFilter(qf *utils.QueryFilter) *SqlBuilder {
	if qf == nil {
		return sb
	}
	sb.QueryFilter = qf
	return sb
}
func (sb *SqlBuilder) SetHavingFilter(qf *utils.QueryFilter) *SqlBuilder {
	if qf == nil {
		return sb
	}
	sb.HavingQueryFilter = qf
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
func (sb *SqlBuilder) buildSelectSql() *SqlBuilder {
	sb.WriteSqlStr(sb.SelectSql)
	if sb.SelectSql == "" {
		sb.WriteSqlStr("SELECT *")
	}
	return sb
}
func (sb *SqlBuilder) buildFromSql() *SqlBuilder {
	sb.WriteSqlStr(sb.FromSql)
	return sb
}
func (sb *SqlBuilder) buildQuerySql() *SqlBuilder {
	sb.WriteSqlStr(sb.QueryFilter.Query)
	return sb
}
func (sb *SqlBuilder) buildHavingSql() *SqlBuilder {
	sb.WriteSqlStr(sb.HavingQueryFilter.Query)
	return sb
}
func (sb *SqlBuilder) buildGroupSql() *SqlBuilder {
	sb.WriteSqlStr(sb.GroupSql)
	return sb
}
func (sb *SqlBuilder) buildSortSql() *SqlBuilder {
	sb.WriteSqlStr(sb.SortSql)
	return sb
}
func (sb *SqlBuilder) buildLimitSql(offset, pageSize int) *SqlBuilder {
	sb.WriteSqlStr(BuildLimitSql(offset, pageSize))
	return sb
}
func (sb *SqlBuilder) BuildCountSql() *SqlBuilder {
	sb.buildSelectSql().buildFromSql().buildQuerySql().buildGroupSql().buildHavingSql().buildSortSql()
	sb.QueryFilter.Args = append(sb.QueryFilter.Args, sb.HavingQueryFilter.Args...)
	sqlStr := BuildCountSql(sb.SqlString())
	sb.sqlStr.Reset()
	sb.WriteSqlStr(sqlStr)
	return sb
}

func (sb *SqlBuilder) SqlString() string {
	return sb.sqlStr.String()
}
func (sb *SqlBuilder) BuildNoLimitSql() *SqlBuilder {
	sb.buildSelectSql().buildFromSql().buildQuerySql().buildGroupSql().buildHavingSql().buildSortSql()
	sb.QueryFilter.Args = append(sb.QueryFilter.Args, sb.HavingQueryFilter.Args...)
	return sb

}
func (sb *SqlBuilder) ToSqlString() string {
	sqlStr := sb.SqlString()
	sb.sqlStr.Reset()
	return sqlStr
}
func (sb *SqlBuilder) BuildAndLimitSql(offset, pageSize int) *SqlBuilder {
	sb.sqlStr.Reset()
	sb.BuildNoLimitSql().buildLimitSql(offset, pageSize)
	return sb
}
