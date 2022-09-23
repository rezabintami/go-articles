package helpers

import (
	"errors"
	"go-articles/server/config"
	"log"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

var (
	vd    *validator.Validate
	trans ut.Translator
)

// Validate checking data
func Validate(data interface{}) (map[string]string, bool, error) {
	err := vd.Struct(data)
	if err != nil {
		errorMap := map[string]string{}
		errorType := ""
		for _, e := range err.(validator.ValidationErrors) {
			if errorType == "" {
				errorType = e.Tag()
				errorMap[LcFirst(e.Field())] = formatMessage(e, trans)
			} else {
				if e.Tag() == errorType {
					errorMap[LcFirst(e.Field())] = formatMessage(e, trans)
				}
			}
		}
		return errorMap, true, nil
	}
	return map[string]string{}, false, nil
}

func setValidate(v *validator.Validate) (ut.Translator, error) {
	translator := en.New()
	uni := ut.New(translator, translator)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Println("[Error] Configuration.setValidate : translator not found")
		return nil, errors.New("translator not found")
	}

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Println("[Error] Configuration.setValidate : ", err)
		return nil, err
	}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("validName"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return trans, nil
}

// Validation for initial function validation
func InitValidation() error {
	vd = validator.New()
	vd.RegisterValidation("phone", ValidatePhone)
	translator, err := setValidate(vd)
	trans = translator
	if err != nil {
		log.Println("[Error] Validation.Validation : ", err)
		return err
	}
	return nil
}

func formatMessage(err validator.FieldError, trans ut.Translator) string {
	message := ""
	switch err.Tag() {
	case "phone":
		message = config.Message("validateMessage.phoneNumber", nil)
	default:
		message = err.Translate(trans)
	}
	return LcFirst(message)
}

func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// ValidatePhone custom validation for phone number
func ValidatePhone(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	regexPhone := regexp.MustCompile(`^([+][6][2][8])[0-9]*$`)
	return regexPhone.MatchString(fl.Field().String())
}

