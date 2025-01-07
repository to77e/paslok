package validator

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	Validator *validator.Validate
	once      sync.Once
)

func GetInstance() *validator.Validate {
	once.Do(func() {
		Validator = validator.New()
	})
	return Validator
}
