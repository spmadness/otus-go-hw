package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrInputValueNotStruct        = errors.New("input value type must be struct")
	ErrTypeNotSupported           = errors.New("not supported type used")
	ErrValidatorNotSupported      = errors.New("not supported validator used")
	ErrEmptyValidatorHandlerValue = errors.New("empty validator handler value")
	ErrEmptyValidatorValue        = errors.New("empty validator value")

	ErrStringWrongLength   = errors.New("string length is not of required value")
	ErrStringNoRegExpMatch = errors.New("string doesn't match regexp")
	ErrStringNotInSet      = errors.New("string not found in set")
	ErrIntMinError         = errors.New("integer is less than min value")
	ErrIntMaxError         = errors.New("integer is greater than max value")
	ErrIntNotInSet         = errors.New("integer not found in set")
	ErrEmptySlice          = errors.New("slice is empty")
)

func (v ValidationErrors) Error() string {
	messages := make([]string, 0)
	for _, ve := range v {
		msg := fmt.Sprintf("Field: %s, Error: %s", ve.Field, ve.Err)
		messages = append(messages, msg)
	}
	return strings.Join(messages, "| ")
}

func Validate(v interface{}) error {
	ve := ValidationErrors{}

	vType := reflect.TypeOf(v)
	vVal := reflect.ValueOf(v)

	if vVal.Kind() != reflect.Struct {
		return ErrInputValueNotStruct
	}

	for i := 0; i < vType.NumField(); i++ {
		fType := vType.Field(i)
		fVal := vVal.Field(i)

		if !fType.IsExported() || fType.Tag == "" {
			continue
		}

		if err := processField(fVal, fType, &ve); err != nil {
			return err
		}
	}

	if len(ve) > 0 {
		return ve
	}

	return nil
}

func processField(fVal reflect.Value, fType reflect.StructField, ve *ValidationErrors) error {
	validateStr, ok := fType.Tag.Lookup("validate")
	if !ok {
		return nil
	}

	validators := strings.Split(validateStr, "|")

	if fVal.Kind() == reflect.Struct {
		hasNested := false
		for _, vv := range validators {
			if vv == NestedValidator {
				hasNested = true
			}
		}
		if !hasNested {
			return nil
		}
		return processStruct(fVal, ve)
	}

	if fVal.Kind() == reflect.Slice {
		return processSlice(fVal, fType, validators, ve)
	}

	err := processValidation(fVal, fType, validators, ve)
	if err != nil {
		return err
	}

	return nil
}

func processStruct(fVal reflect.Value, ve *ValidationErrors) error {
	var veLocal ValidationErrors

	err := Validate(fVal.Interface())

	if errors.As(err, &veLocal) {
		*ve = append(*ve, veLocal...)
	} else if err != nil {
		return err
	}
	return nil
}

func processSlice(fVal reflect.Value, fType reflect.StructField, validators []string, ve *ValidationErrors) error {
	if fVal.Len() == 0 {
		*ve = append(*ve, ValidationError{
			fType.Name,
			ErrEmptySlice,
		})
		return nil
	}

	for i := 0; i < fVal.Len()-1; i++ {
		err := processValidation(fVal.Index(i), fType, validators, ve)
		if err != nil {
			return err
		}
	}

	return nil
}

func processValidation(fVal reflect.Value, fType reflect.StructField, validators []string, ve *ValidationErrors) error {
	for _, validator := range validators {
		validatorSlice := strings.Split(validator, ":")
		validatorName := validatorSlice[0]
		validatorVal := validatorSlice[1]
		if validatorVal == "" {
			return ErrEmptyValidatorValue
		}

		kind, ok := ValidatorMap[fVal.Kind()]
		if !ok {
			return ErrTypeNotSupported
		}

		f, ok := kind[validatorName]
		if !ok {
			return ErrValidatorNotSupported
		}

		if f == nil {
			return ErrEmptyValidatorHandlerValue
		}

		ok, err := f(fVal, validatorVal)
		if !ok {
			*ve = append(*ve, ValidationError{
				Field: fType.Name,
				Err:   fmt.Errorf("%w", err),
			})
		}
	}

	return nil
}
