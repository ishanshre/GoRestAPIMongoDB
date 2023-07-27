package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func UpperCase(f1 validator.FieldLevel) bool {
	field := f1.Field().String()
	exp, err := regexp.Compile("([A-Z])")
	if err != nil {
		return false
	}
	u := exp.FindAllString(field, 1)
	return len(u) != 0
}

func LowerCase(f1 validator.FieldLevel) bool {
	field := f1.Field().String()
	exp, err := regexp.Compile("([a-z])")
	if err != nil {
		return false
	}
	u := exp.FindAllString(field, 1)
	return len(u) != 0
}

func Number(f1 validator.FieldLevel) bool {
	field := f1.Field().String()
	exp, err := regexp.Compile("([0-9])")
	if err != nil {
		return false
	}
	u := exp.FindAllString(field, 1)
	return len(u) != 0
}
