package agtwo

import "fmt"

type SqlStr string

func BuildCountSql(sql string) string {
	return fmt.Sprintf("SELECT COUNT(1) FROM (%s) AS a", sql)
}
func BuildLimitSql(offset, pageSize int) string {
	return fmt.Sprintf("LIMIT %d,%d", offset, pageSize)
}
