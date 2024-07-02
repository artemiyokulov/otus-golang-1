package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			// case 1: correct User object
			in: User{
				ID:     "f8a9df52-a90a-4127-b5aa-844a826cc0ef",
				Name:   "name1",
				Age:    18,
				Email:  "asd@asd.ru",
				Role:   "admin",
				Phones: []string{"111-111-111", "222-222-222"},
				meta:   nil,
			},
			expectedErr: nil,
		},
		// case 2: incorrect App object
		{
			in: App{
				Version: "1234",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   createError("1234", ValidatorTag{name: "len", args: []string{"5"}}),
				},
			},
		},
		// case 3: correct Token object (ignore without tags)
		{
			in:          Token{},
			expectedErr: nil,
		},
		// case 4: correct Response object
		{
			in: Response{
				Code: 200,
				Body: "ok",
			},
			expectedErr: nil,
		},
		// case 4: Incorrect Response object
		{
			in: Response{
				Code: 201,
				Body: "so-so",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   createError(201, ValidatorTag{name: "in", args: []string{"200", "404", "500"}}),
				},
			},
		},
		// case 5: incorrect User object
		{
			in: User{
				ID:     "a",
				Name:   "name1",
				Age:    17,
				Email:  "---",
				Role:   "any",
				Phones: []string{"333-333-333", "111-222"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   createError("a", ValidatorTag{name: "len", args: []string{"36"}}),
				},
				ValidationError{
					Field: "Age",
					Err:   createError(17, ValidatorTag{name: "min", args: []string{"18"}}),
				},
				ValidationError{
					Field: "Email",
					Err:   createError("---", ValidatorTag{name: "regexp", args: []string{"^\\w+@\\w+\\.\\w+$"}}),
				},
				ValidationError{
					Field: "Role",
					Err:   createError("any", ValidatorTag{name: "in", args: []string{"admin", "stuff"}}),
				},
				ValidationError{
					Field: "Phones",
					Err:   createError("111-222", ValidatorTag{name: "len", args: []string{"11"}}),
				},
			},
		},
		{
			// case 6: incorrect User object
			in: User{
				ID:     "f8a9df52-a90a-4127-b5aa-844a826cc0ef",
				Name:   "name1",
				Age:    56,
				Email:  "asd@asd.ru",
				Role:   "admin",
				Phones: []string{"111-111-111", "222-222-222"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   createError(56, ValidatorTag{name: "max", args: []string{"50"}}),
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			assert.Equal(t, tt.expectedErr, err)
			_ = tt
		})
	}
}
