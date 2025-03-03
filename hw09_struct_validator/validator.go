package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string

	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errStrings := make([]string, 0, len(v))
	for _, err := range v {
		errStrings = append(errStrings, fmt.Sprintf("%s: %s", err.Field, err.Err))
	}
	return strings.Join(errStrings, "; ")
}

type ProgrammError struct {
	Message string
}

func (e ProgrammError) Error() string {
	return e.Message
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.New("input is not a struct")
	}

	var validationErrors ValidationErrors

	for i := range val.NumField() {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)
		validateTag := field.Tag.Get("validate")

		if validateTag == "" {
			continue
		}

		validators := strings.Split(validateTag, "|")

		for _, validator := range validators {
			err := applyValidator(fieldValue, validator)
			if err != nil {
				var progErr ProgrammError
				if errors.As(err, &progErr) {
					return err // Прекращаем валидацию
				}
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			}
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

// applyValidator применяет конкретный валидатор к значению поля.

func applyValidator(value reflect.Value, validator string) error {
	parts := strings.SplitN(validator, ":", 2)

	if len(parts) != 2 {
		return ProgrammError{Message: fmt.Sprintf("invalid validator format: %s", validator)}
	}

	rule := parts[0]
	param := parts[1]

	switch value.Kind() {
	case reflect.String:
		return validateString(value.String(), rule, param)
	case reflect.Int:
		return validateInt(int(value.Int()), rule, param)
	case reflect.Slice:
		return validateSlice(value, rule, param)
	default:
		return ProgrammError{Message: fmt.Sprintf("unsupported type for validation: %s", value.Kind())}
	}
}

// validateString проверяет строку на соответствие правилам.

func validateString(value, rule, param string) error {
	switch rule {
	case "len":
		expectedLen, err := strconv.Atoi(param)
		if err != nil {
			return ProgrammError{Message: fmt.Sprintf("invalid length parameter: %s", param)}
		}

		if len([]rune(value)) != expectedLen {
			return fmt.Errorf("length must be %d", expectedLen)
		}
	case "regexp":
		re, err := regexp.Compile(param)
		if err != nil {
			return ProgrammError{Message: fmt.Sprintf("invalid regexp: %s", param)}
		}

		if !re.MatchString(value) {
			return fmt.Errorf("must match regexp %s", param)
		}
	case "in":
		allowedValues := strings.Split(param, ",")
		if slices.Contains(allowedValues, value) {
			return nil
		}
		return fmt.Errorf("must be one of %v", allowedValues)
	default:
		return ProgrammError{Message: fmt.Sprintf("unknown validator for string: %s", rule)}
	}
	return nil
}

// validateInt проверяет целое число на соответствие правилам.

func validateInt(value int, rule, param string) error {
	switch rule {
	case "min":
		minVal, err := strconv.Atoi(param)
		if err != nil {
			return ProgrammError{Message: fmt.Sprintf("invalid min parameter: %s", param)}
		}

		if value < minVal {
			return fmt.Errorf("must be at least %d", minVal)
		}

	case "max":
		maxVal, err := strconv.Atoi(param)
		if err != nil {
			return ProgrammError{Message: fmt.Sprintf("invalid max parameter: %s", param)}
		}

		if value > maxVal {
			return fmt.Errorf("must be at most %d", maxVal)
		}

	case "in":
		allowedValues := strings.Split(param, ",")

		for _, allowed := range allowedValues {
			allowedInt, err := strconv.Atoi(allowed)
			if err != nil {
				return ProgrammError{Message: fmt.Sprintf("invalid int value in 'in' validator: %s", allowed)}
			}

			if value == allowedInt {
				return nil
			}
		}

		return fmt.Errorf("must be one of %v", allowedValues)
	default:
		return ProgrammError{Message: fmt.Sprintf("unknown validator for int: %s", rule)}
	}

	return nil
}

func validateSlice(value reflect.Value, rule, param string) error {
	elemKind := value.Type().Elem().Kind()

	for i := 0; i < value.Len(); i++ {
		switch elemKind {
		case reflect.String:

			err := validateString(value.Index(i).String(), rule, param)
			if err != nil {
				return fmt.Errorf("element %d: %w", i, err)
			}

		case reflect.Int:

			err := validateInt(int(value.Index(i).Int()), rule, param)
			if err != nil {
				return fmt.Errorf("element %d: %w", i, err)
			}

		default:

			return fmt.Errorf("unsupported slice element type: %s", elemKind)
		}
	}

	return nil
}
