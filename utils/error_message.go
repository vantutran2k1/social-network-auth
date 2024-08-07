package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type ErrorMessage struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func BindAndValidate(c *gin.Context, obj any) []ErrorMessage {
	err := c.ShouldBindJSON(obj)
	if err == nil {
		return nil
	}

	var ve validator.ValidationErrors
	ok := errors.As(err, &ve)
	if !ok {
		return []ErrorMessage{
			{Message: err.Error()},
		}
	}

	bindErrs := make([]ErrorMessage, len(ve))
	for i, fe := range ve {
		bindErrs[i] = ErrorMessage{
			Field:   getFieldName(fe, obj),
			Message: getErrorMsg(fe),
		}
	}

	return bindErrs
}

func getFieldName(fe validator.FieldError, obj any) string {
	fieldName := fe.Field()

	field, _ := reflect.TypeOf(obj).Elem().FieldByName(fieldName)
	jsonName := field.Tag.Get("json")
	if jsonName != "" {
		fieldName = jsonName
	}
	return fieldName
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte", "max":
		return "Should be less than or equal to " + fe.Param()
	case "gte", "min":
		return "Should be greater than or equal to " + fe.Param()
	case "lt":
		return "Should be less than " + fe.Param()
	case "gt":
		return "Should be greater than " + fe.Param()
	case "eq":
		return "Should be equal to " + fe.Param()
	case "ne":
		return "Should be not equal to " + fe.Param()
	case "email":
		return "Should be a valid email address"
	case "oneof":
		return "Should be one of " + strings.Join(strings.Split(fe.Param(), " "), ", ")
	case "date":
		return "Should be a valid date with format YYYY-MM-DD"
	case "beforeToday":
		return "Should be a valid date before today"
	case "phoneNumber":
		return "Should be a valid phone number"
	}
	return "Unknown error"
}
