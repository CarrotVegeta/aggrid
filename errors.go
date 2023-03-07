package agtwo

import "errors"

var (
	InvalidGroupField    = errors.New("invalid group field")
	InvalidSortField     = errors.New("invalid sort field")
	InvalidSqlField      = errors.New("invalid sql field")
	InvalidFromSql       = errors.New("from sql is null")
	BuildNoLimitSqlError = errors.New("build no limit sql error")
	BuildLimitSqlError   = errors.New("build  limit sql error")
	RawSqlError          = errors.New("to raw sql error")
)
