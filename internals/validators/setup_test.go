package validators

import (
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func TestMain(m *testing.M) {
	validate = validator.New()
	validate.RegisterValidation("uppercase", UpperCase)
	validate.RegisterValidation("lowercase", LowerCase)
	validate.RegisterValidation("number", Number)
	os.Exit(m.Run())
}
