package agtwo

import (
	"fmt"
	"log"
	"testing"
)

type User struct {
	Name string `json:"name" ag:"select:name"`
}

func (u User) GetSqlField(k string) string {
	return ""
}
func (u User) GetSelectField() []string {
	return []string{}
}
func TestAgGrid_BuildSelect(t *testing.T) {
	ag := NewAgGHandler(User{}, nil)
	ag.parse()
	fmt.Println(ag)
}
func TestAgGrid(t *testing.T) {
	var sortModels []SortModel
	sortModels = append(sortModels, SortModel{
		Sort:  "asc",
		ColId: "name",
	})
	sortModels = append(sortModels, SortModel{
		Sort:  "asc",
		ColId: "timestamp",
	})
	param := &Param{RowGroupCols: []RowGroupCol{{Field: "name", DisplayName: "name"}}, SortModels: sortModels}
	ag := NewAgGHandler(User{}, param)
	selectSql, err := ag.BuildSelectSql()
	if err != nil {
		t.Fatalf(err.Error())
	}
	sortSql, err := ag.BuildSortSql()
	sb := &SqlBuilder{}
	sb.SetSelectSql(selectSql).SetSortSql(sortSql).SetFromSql("FROM user")
	sql := sb.BuildNoLimitSql().ToSqlString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(sql)

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
