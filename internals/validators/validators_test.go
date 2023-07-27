package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpperCase(t *testing.T) {
	type TestStruct struct {
		Field1 string `validate:"uppercase"`
	}
	validStruct := TestStruct{
		Field1: "Abc123",
	}
	err := validate.Struct(validStruct)
	assert.NoError(t, err)

	invalidStruct := TestStruct{
		Field1: "asdasd",
	}
	err = validate.Struct(invalidStruct)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "uppercase")
}
func TestLowerCase(t *testing.T) {
	type TestStruct struct {
		Field1 string `validate:"lowercase"`
	}
	validStruct := TestStruct{
		Field1: "Abc123",
	}
	err := validate.Struct(validStruct)
	assert.NoError(t, err)

	invalidStruct := TestStruct{
		Field1: "ASSSS",
	}
	err = validate.Struct(invalidStruct)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lowercase")
}
func TestNumber(t *testing.T) {
	type TestStruct struct {
		Field1 string `validate:"number"`
	}
	validStruct := TestStruct{
		Field1: "Abc123",
	}
	err := validate.Struct(validStruct)
	assert.NoError(t, err)

	invalidStruct := TestStruct{
		Field1: "asdasd",
	}
	err = validate.Struct(invalidStruct)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "number")
}
