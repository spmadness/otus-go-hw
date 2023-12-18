package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
		Meta   Meta            `validate:"nested"`
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Meta struct {
		Gender     string `validate:"in:male,female"`
		CardNumber string `validate:"regexp:^\\d+$|len:16"`
		Status     int    `validate:"in:1,2,3,4,5"`
		Bio        string
	}

	BadStruct struct {
		BadValidatorValueString string `validate:"in:"`
	}

	BadStructTwo struct {
		UnsupportedValidatorField string `validate:"trim:rty"`
	}

	BadToken struct {
		BadTypeHeader []byte `validate:"in:1,2"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{"not a struct", ErrInputValueNotStruct},

		{BadStruct{BadValidatorValueString: "string"}, ErrEmptyValidatorValue},
		{BadToken{BadTypeHeader: []byte{66, 67, 68}}, ErrTypeNotSupported},
		{BadStructTwo{UnsupportedValidatorField: "qwerty"}, ErrValidatorNotSupported},

		{User{
			ID:     "123456789012345678901234567890123456",
			Age:    45,
			Email:  "test@test.com",
			Role:   "admin",
			Phones: []string{"+7999554400", "+7999554411", "+7999554422"},
			Meta: Meta{
				Gender:     "male",
				CardNumber: "1234567890123456",
				Status:     3,
			},
		}, nil},
		{User{}, ValidationErrors{
			ValidationError{Field: "ID", Err: ErrStringWrongLength},
			ValidationError{Field: "Age", Err: ErrIntMinError},
			ValidationError{Field: "Email", Err: ErrStringNoRegExpMatch},
			ValidationError{Field: "Role", Err: ErrStringNotInSet},
			ValidationError{Field: "Phones", Err: ErrEmptySlice},
			ValidationError{Field: "Gender", Err: ErrStringNotInSet},
			ValidationError{Field: "CardNumber", Err: ErrStringNoRegExpMatch},
			ValidationError{Field: "CardNumber", Err: ErrStringWrongLength},
			ValidationError{Field: "Status", Err: ErrIntNotInSet},
		}},
		{User{
			ID:     "",
			Email:  "test@test.com",
			Role:   "admin",
			Phones: []string{"+7999554400", "+7999", "+79995544225345"},
			Meta: Meta{
				Gender:     "male",
				CardNumber: "xxx",
				Status:     6,
			},
		}, ValidationErrors{
			ValidationError{Field: "ID", Err: ErrStringWrongLength},
			ValidationError{Field: "Age", Err: ErrIntMinError},
			ValidationError{Field: "Phones", Err: ErrStringWrongLength},
			ValidationError{Field: "CardNumber", Err: ErrStringNoRegExpMatch},
			ValidationError{Field: "CardNumber", Err: ErrStringWrongLength},
			ValidationError{Field: "Status", Err: ErrIntNotInSet},
		}},

		{Response{Code: 404, Body: "empty string"}, nil},
		{Response{Code: 405}, ValidationErrors{
			ValidationError{Field: "Code", Err: ErrIntNotInSet},
		}},

		{App{Version: "1.0.3"}, nil},
		{
			App{Version: "1.0.33"}, ValidationErrors{
				ValidationError{Field: "Version", Err: ErrStringWrongLength},
			},
		},
		// no validate tags in Token struct fields
		{Token{Header: []byte{66, 67, 68}, Payload: []byte{66, 67, 68}, Signature: []byte{66, 67, 68}}, nil},
	}

	for i, tCase := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tCase := tCase
			t.Parallel()

			err := Validate(tCase.in)

			if err == nil || tCase.expectedErr == nil {
				require.Equalf(t, tCase.expectedErr, err,
					"expected and actual errors must be nil, expected: %s, actual: %s", tCase.expectedErr, err)
				return
			}
			var veActual ValidationErrors
			var veExpected ValidationErrors

			if !errors.As(err, &veActual) {
				require.Equalf(t, tCase.expectedErr, err,
					"non-validation errors are not equal, expected : %s, actual : %s", tCase.expectedErr, err)
				return
			}

			if !errors.As(tCase.expectedErr, &veExpected) {
				t.Error("expected error is not of ValidationErrors type")
			}

			require.Truef(t, len(veActual) == len(veExpected), "actual and expected ValidationError count are not equal")

			for i := 0; i < len(veActual); i++ {
				require.Truef(t, veExpected[i].Field == veActual[i].Field,
					"field names are not equal, expected : %s, actual : %s", veExpected[i].Field, veActual[i].Field)
				require.Truef(t, errors.Is(veActual[i].Err, veExpected[i].Err),
					"error types are not equal, expected : %s, actual : %s", veActual[i].Err, veExpected[i].Err)
			}
		})
	}
}
