package agtwo

import (
	"fmt"
	"testing"
)

type User struct {
	Name string `json:"name" ag:"select:sdfkd;group:sdlsdfkj"`
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
