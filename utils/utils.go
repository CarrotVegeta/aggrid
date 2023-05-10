package utils

import (
	"reflect"
	"strings"
)

const sortASC = "ASC"
const sortDESC = "DESC"

func validSort(sort string) bool {
	s := strings.ToUpper(sort)
	if s != sortASC {
		if s != sortDESC {
			return false
		}
	}
	return true
}

type FieldName map[string]struct{}

type StructTag map[string]reflect.StructTag

func GetStructTagField(model any, tag string) StructTag {
	refValue := reflect.TypeOf(model)
	fn := StructTag{}
	for i := 0; i < refValue.NumField(); i++ {
		k := refValue.Field(i).Tag.Get(tag)
		fn[k] = refValue.Field(i).Tag
	}
	return fn
}
