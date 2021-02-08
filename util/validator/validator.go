package validator

import (
	"Blog/util"
	"Blog/util/errmsg"
	"reflect"
	"sync"

	"github.com/go-playground/locales/zh_Hans_CN"
	uvTrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
)

var logger = util.GetLogger()

var vali *validator.Validate
var translator uvTrans.Translator
var once sync.Once

func init() {
	once.Do(func() {
		initValidator()
	})
}

func initValidator() {
	logger.Info("Init validator...")
	vali = validator.New()
	uvTranslator := uvTrans.New(zh_Hans_CN.New())
	translator, _ = uvTranslator.GetTranslator("zh_Hans_CN")
	err := zhTrans.RegisterDefaultTranslations(vali, translator)
	if err != nil {
		logger.Error("Fail to registe translator:", err.Error())
	}
	vali.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("label")
	})
}
func Validate(data interface{}) (string, errmsg.Code) {
	err := vali.Struct(data)
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			return v.Translate(translator), errmsg.ERROR
		}
	}
	return "", errmsg.SUCCESS
}
