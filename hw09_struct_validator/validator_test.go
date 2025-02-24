package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
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
			in: User{
				ID:    "123456789012345678901234567890123456", // len:36
				Age:   25,
				Email: "test@example.com",
				Role:  "admin",
				Phones: []string{
					"12345678901", // len:11
					"09876543210", // len:11
				},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:    "123",           // len:3 (invalid)
				Age:   15,              // min:18 (invalid)
				Email: "invalid-email", // regexp (invalid)
				Role:  "user",          // in:admin,stuff (invalid)
				Phones: []string{
					"123", // len:3 (invalid)
				},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("length must be 36")},
				{Field: "Age", Err: fmt.Errorf("must be at least 18")},
				{Field: "Email", Err: fmt.Errorf("must match regexp ^\\w+@\\w+\\.\\w+$")},
				{Field: "Role", Err: fmt.Errorf("must be one of [admin stuff]")},
				{Field: "Phones", Err: fmt.Errorf("element 0: length must be 11")},
			},
		},
		{
			in: App{
				Version: "11.0.0", // len:5 (invalid)
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: fmt.Errorf("length must be 5")},
			},
		},
		{
			in: Response{
				Code: 400, // in:200,404,500 (invalid)
				Body: "OK",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("must be one of [200 404 500]")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if err == nil && tt.expectedErr == nil {
				return
			}

			if err == nil || tt.expectedErr == nil {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
				return
			}

			var expectedValidationErrors ValidationErrors
			if !errors.As(tt.expectedErr, &expectedValidationErrors) {
				t.Errorf("expected error type: ValidationErrors, got: %T", tt.expectedErr)
				return
			}

			var actualValidationErrors ValidationErrors
			if !errors.As(err, &actualValidationErrors) {
				t.Errorf("expected error type: ValidationErrors, got: %T", err)
				return
			}

			if len(expectedValidationErrors) != len(actualValidationErrors) {
				t.Errorf("expected %d errors, got %d", len(expectedValidationErrors), len(actualValidationErrors))
				return
			}

			for j := range expectedValidationErrors {
				if expectedValidationErrors[j].Field != actualValidationErrors[j].Field ||
					expectedValidationErrors[j].Err.Error() != actualValidationErrors[j].Err.Error() {
					t.Errorf("expected error: %v, got: %v", expectedValidationErrors[j], actualValidationErrors[j])
				}
			}
		})
	}
}
