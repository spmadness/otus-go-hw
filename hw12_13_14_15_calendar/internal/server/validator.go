package server

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
)

func ValidateCreateEvent(e storage.Event) error {
	validate := validator.New()

	return ProcessRequestData(e, validate.StructExcept, "ID")
}

func ValidateUpdateEvent(e storage.Event) error {
	validate := validator.New()

	return ProcessRequestData(e, validate.StructExcept, "")
}

func ValidateDeleteEvent(e storage.Event) error {
	validate := validator.New()

	return ProcessRequestData(e, validate.StructPartial, "ID")
}

func ValidateListEvent(lm storage.ListEventValidation) error {
	validate := validator.New()

	return ProcessRequestData(lm, validate.StructExcept, "")
}

func ProcessRequestData(data interface{}, f func(s interface{}, fields ...string) error, fields ...string) error {
	err := f(data, fields...)
	if err != nil {
		var str string

		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			for _, e := range validationErrs {
				str += e.Error() + "\n"
			}
			return errors.New(str)
		}

		return err
	}

	return nil
}
