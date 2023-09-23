package request

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Request interface {
	Validate(interface{}) interface{}
}

type Validator struct {
	Errors  map[string]string
	Context *gin.Context
}

func (v *Validator) Validate(request interface{}) interface{} {
	err := v.Context.ShouldBindJSON(request)
	if err != nil {
		v.Errors["json"] = err.Error()
		return nil
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			v.Errors[err.Field()] = err.Tag()
		}
		return nil
	}
	return request
}

func V(c *gin.Context) *Validator {
	return &Validator{
		Errors:  make(map[string]string),
		Context: c,
	}
}
