package util

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Struct2Map(obj interface{}) gin.H {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	data := make(gin.H)
	for i := 0; i < t.NumField(); i++ {
		value := v.Field(i).Interface()
		if value == nil {

		}
		//如果struct里还有struct
		if _, ok := value.(gorm.Model); ok {
			tv := reflect.TypeOf(value)
			vv := reflect.ValueOf(value)
			for j := 0; j < tv.NumField(); j++ {
				data[strings.ToLower(tv.Field(i).Name)] = vv.Field(i).Interface()
			}
		} else {
			data[strings.ToLower(t.Field(i).Name)] = value
		}
	}
	return data
}
