package interfaces

import (
	"github.com/CarrotVegeta/aggrid/constant"
	"github.com/CarrotVegeta/aggrid/utils"
)

type FilterTypeSqlHandler interface {
	New() FilterTypeSqlHandler
	BuildSql(k string, v any, t constant.OperatorType, f ...constant.F) (*utils.QueryFilter, error)
}
