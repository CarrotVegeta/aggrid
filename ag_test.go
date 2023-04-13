package agtwo

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"testing"
)

var db *gorm.DB

func InitDB() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(err.Error())
	}
	db = gormDB

}

// select 代表json对应的查询字段名,filter:代表where里面执行的字段名,group:代表group的字段名
type User struct {
	Name string `json:"name" ag:"select:name;filter:name;group:name"`
}

func (u User) BuildFromSql() string {
	return "FROM users"
}

func TestAgGrid_BuildSelect(t *testing.T) {
	ag := NewAgGHandler(User{}, nil)
	ag.parse()
	fmt.Println(ag)
}
func TestAgGrid(t *testing.T) {
	InitDB()
	var sortModels []SortModel
	sortModels = append(sortModels, SortModel{
		Sort:  "asc",
		ColId: "name",
	})
	sortModels = append(sortModels, SortModel{
		Sort:  "asc",
		ColId: "timestamp",
	})
	param := &Param{SortModels: sortModels}
	ag := NewAgGHandler(User{}, param)
	qf, err := ag.BuildQuerySql()
	if err != nil {
		t.Fatalf(err.Error())
	}
	var count int64
	var result []map[string]any
	err = ag.Use(db).Where(qf).Count(&count).Find(&result).Error
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(result)

}
func TestAgGridSort(t *testing.T) {

	var sortModels []SortModel
	sortModels = append(sortModels, SortModel{
		Sort:  "asc",
		ColId: "name",
	})
	param := &Param{SortModels: sortModels}
	ag := NewAgGHandler(User{}, param)
	selectSql, err := ag.BuildSelectSql()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(selectSql)
	sql, err := ag.BuildSortSql()
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(sql)
}
func TestAgGridGroup(t *testing.T) {

	var sortModels []SortModel
	sortModels = append(sortModels, SortModel{
		Sort:  "asc",
		ColId: "name",
	})
	param := &Param{RowGroupCols: []RowGroupCol{{Field: "name", DisplayName: "name"}}, SortModels: sortModels}
	ag := NewAgGHandler(User{}, param)
	selectSql, err := ag.buildGroupSelect()
	if err != nil {
		t.Fatalf(err.Error())
	}
	if selectSql == "" {
		t.Fatalf("invalid")
	}
	t.Log(selectSql)
}
