package filtertype

import (
	"fmt"
	"github.com/CarrotVegeta/aggrid/clickhouse"
	"github.com/CarrotVegeta/aggrid/constant"
	"github.com/CarrotVegeta/aggrid/interfaces"
	"github.com/CarrotVegeta/aggrid/mysql"
	"github.com/CarrotVegeta/aggrid/utils"
)

type FilterHandler interface {
	Parse(k string, c []byte) error
	BuildQuery(service *FilterTypeSqlService) (*utils.QueryFilter, error)
}

type FilterTypeSqlHandlerM map[constant.FilterType]interfaces.FilterTypeSqlHandler

var DefaultFilterTypeMysqlM = make(FilterTypeSqlHandlerM)
var DefaultFilterTypeClickhouseM = make(FilterTypeSqlHandlerM)

var defaultFilterTypeM = make(map[string]FilterTypeSqlHandlerM)

func initMysql() {
	registerFilterTypeMysqlHandler(constant.Text, &mysql.FilterText{})
	registerFilterTypeMysqlHandler(constant.Number, &mysql.FilterNumber{})
	registerFilterTypeMysqlHandler(constant.Date, &mysql.FilterDate{})
	registerFilterTypeMysqlHandler(constant.Array, &mysql.FilterArray{})
	registerFilterTypeMysqlHandler(constant.Set, &mysql.FilterSet{})
}
func initClickhouse() {
	registerFilterTypeClickhouseHandler(constant.Text, &clickhouse.FilterText{})
	registerFilterTypeClickhouseHandler(constant.Number, &clickhouse.FilterNumber{})
	registerFilterTypeClickhouseHandler(constant.Date, &clickhouse.FilterDate{})
	registerFilterTypeClickhouseHandler(constant.Array, &clickhouse.FilterArray{})
	registerFilterTypeClickhouseHandler(constant.Set, &clickhouse.FilterSet{})
}
func initDefaultFilter() {
	defaultFilterTypeM = make(map[string]FilterTypeSqlHandlerM)
	defaultFilterTypeM[constant.MySql] = DefaultFilterTypeMysqlM
	defaultFilterTypeM[constant.Clickhouse] = DefaultFilterTypeClickhouseM
}
func init() {
	initMysql()
	initClickhouse()
	initDefaultFilter()
}

func registerFilterTypeMysqlHandler(k constant.FilterType, handler interfaces.FilterTypeSqlHandler) {
	DefaultFilterTypeMysqlM[k] = handler
}
func registerFilterTypeClickhouseHandler(k constant.FilterType, handler interfaces.FilterTypeSqlHandler) {
	DefaultFilterTypeClickhouseM[k] = handler
}

type FilterTypeSqlService struct {
	FilterTypeSqlHandlerM map[constant.FilterType]interfaces.FilterTypeSqlHandler
	StorageType           string
}

func NewFilterTypeSqlService(storageType string) *FilterTypeSqlService {
	f := &FilterTypeSqlService{
		FilterTypeSqlHandlerM: make(map[constant.FilterType]interfaces.FilterTypeSqlHandler),
		StorageType:           storageType,
	}
	f.initTypeSqlHandlerM()
	return f
}

func (f *FilterTypeSqlService) initTypeSqlHandlerM() {
	if f.StorageType == "" {
		f.StorageType = constant.DefaultStorage
	}
	for k, v := range defaultFilterTypeM[f.StorageType] {
		f.FilterTypeSqlHandlerM[k] = v
	}
}
func (f *FilterTypeSqlService) registerFilterTypeSqlHandler(k constant.FilterType, handler interfaces.FilterTypeSqlHandler) {
	f.FilterTypeSqlHandlerM[k] = handler
}

func (f *FilterTypeSqlService) getFilterTypeSqlHandler(filterType constant.FilterType) (interfaces.FilterTypeSqlHandler, error) {
	if _, ok := f.FilterTypeSqlHandlerM[filterType]; !ok {
		return nil, fmt.Errorf("invalid filter-type:%v", filterType)
	}
	return f.FilterTypeSqlHandlerM[filterType].New(), nil
}
func (f *FilterTypeSqlService) NewFilterSqlHandler(filterType constant.FilterType) (interfaces.FilterTypeSqlHandler, error) {
	h, err := f.getFilterTypeSqlHandler(filterType)
	if err != nil {
		return nil, err
	}
	return h, nil
}
