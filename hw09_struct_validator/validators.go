package hw09structvalidator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidatorFunc func(reflect.Value, string) (bool, error)

var (
	LengthValidator  = "len"
	RegExpValidator  = "regexp"
	InRangeValidator = "in"
	MinValidator     = "min"
	MaxValidator     = "max"
	NestedValidator  = "nested"
)

var ValidatorMap = map[reflect.Kind]map[string]ValidatorFunc{
	reflect.String: {
		LengthValidator:  StringLengthValidator,
		RegExpValidator:  StringRegExpValidator,
		InRangeValidator: StringInValidator,
	},
	reflect.Int: {
		MinValidator:     IntegerMinValidator,
		MaxValidator:     IntegerMaxValidator,
		InRangeValidator: IntegerInValidator,
	},
}

func StringLengthValidator(fVal reflect.Value, data string) (bool, error) {
	tVal, err := strconv.Atoi(data)
	if err != nil {
		return false, err
	}

	str := fVal.String()

	if len(str) != tVal {
		return false, ErrStringWrongLength
	}
	return true, nil
}

func StringRegExpValidator(fVal reflect.Value, data string) (bool, error) {
	re, err := regexp.Compile(data)
	if err != nil {
		return false, err
	}
	match := re.FindStringSubmatch(fVal.String())
	if match == nil {
		return false, ErrStringNoRegExpMatch
	}

	return true, nil
}

func StringInValidator(fVal reflect.Value, data string) (bool, error) {
	words := strings.Split(data, ",")
	for _, word := range words {
		if word == fVal.String() {
			return true, nil
		}
	}

	return false, ErrStringNotInSet
}

func IntegerMinValidator(fVal reflect.Value, data string) (bool, error) {
	tVal, err := strconv.Atoi(data)
	if err != nil {
		return false, err
	}

	d := fVal.Int()
	if int(d) < tVal {
		return false, ErrIntMinError
	}

	return true, nil
}

func IntegerMaxValidator(fVal reflect.Value, data string) (bool, error) {
	tVal, err := strconv.Atoi(data)
	if err != nil {
		return false, err
	}

	d := fVal.Int()
	if int(d) > tVal {
		return false, ErrIntMaxError
	}

	return true, nil
}

func IntegerInValidator(fVal reflect.Value, data string) (bool, error) {
	intSet := strings.Split(data, ",")
	for _, i := range intSet {
		i, err := strconv.Atoi(i)
		if err != nil {
			return false, err
		}

		if int64(i) == fVal.Int() {
			return true, nil
		}
	}
	return false, ErrIntNotInSet
}
