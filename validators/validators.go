package validators

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
	"time"
)

func RegisterCustomValidators() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("date", isValidDate)
		v.RegisterValidation("beforeToday", isBeforeToday)
		v.RegisterValidation("phoneNumber", isValidPhoneNumber)
	}
}

func isValidDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

func isBeforeToday(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}
	return !date.After(time.Now())
}

func isValidPhoneNumber(fl validator.FieldLevel) bool {
	phoneNumberRegex := regexp.MustCompile(`^\+?[\d\s-]{7,15}$`)
	phone := fl.Field().String()
	return phoneNumberRegex.MatchString(phone)
}
