package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type CheckRule func(value any, ruleArgs ...string) (bool, error)

type ValidatorType struct {
	name          string
	supportedType []reflect.Kind
	checkFunc     CheckRule
}

func (v ValidatorType) isSupported(kind reflect.Kind) bool {
	return slices.Contains(v.supportedType, kind)
}

type ValidatorTag struct {
	name string
	args []string
}

func (v ValidatorTag) String() string {
	return v.name + ":" + strings.Join(v.args, ",")
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	sb.WriteString("Validation errors: \n")
	for _, err := range v {
		sb.WriteString(fmt.Sprintf("field %s: %s\n", err.Field, err.Err.Error()))
	}
	return sb.String()
}

var validators = map[string]ValidatorType{
	"len": {
		name:          "len",
		supportedType: []reflect.Kind{reflect.String},
		checkFunc: CheckRule(func(value any, args ...string) (bool, error) {
			if len(args) != 1 {
				return false, fmt.Errorf("len validator: must be 1 int arg")
			}
			strVal, ok := value.(string)
			if !ok {
				return false, fmt.Errorf("len validator: value must be string")
			}
			arg, err := strconv.Atoi(args[0])
			if err != nil {
				return false, fmt.Errorf("len validator: arg must be int")
			}
			return len(strVal) == arg, nil
		}),
	},
	"regexp": {
		name:          "regexp",
		supportedType: []reflect.Kind{reflect.String},
		checkFunc: CheckRule(func(value any, args ...string) (bool, error) {
			if len(args) != 1 {
				return false, fmt.Errorf("regexp validator: must be 1 int arg")
			}
			strVal, ok := value.(string)
			if !ok {
				return false, fmt.Errorf("regexp validator: value must be string")
			}
			pattern := args[0]
			if !ok {
				return false, fmt.Errorf("regexp validator: arg must be string")
			}
			return regexp.MatchString(pattern, strVal)
		}),
	},
	"in": {
		name:          "in",
		supportedType: []reflect.Kind{reflect.String, reflect.Int},
		checkFunc: CheckRule(func(value any, args ...string) (bool, error) {
			reflectValue := reflect.ValueOf(value)
			switch reflectValue.Kind() {
			case reflect.Int:
				return slices.Contains(args, strconv.Itoa(int(reflectValue.Int()))), nil
			case reflect.String:
				return slices.Contains(args, reflectValue.String()), nil
			default:
				return false, fmt.Errorf("in validator: value type not supported")
			}
		}),
	},
	"min": {
		name:          "min",
		supportedType: []reflect.Kind{reflect.Int},
		checkFunc: CheckRule(func(value any, args ...string) (bool, error) {
			if len(args) != 1 {
				return false, fmt.Errorf("regexp validator: must be 1 int arg")
			}
			intVal, ok := value.(int)
			if !ok {
				return false, fmt.Errorf("regexp validator: value must be int")
			}
			minVal, err := strconv.Atoi(args[0])
			if err != nil {
				return false, fmt.Errorf("regexp validator: arg must be int")
			}
			return intVal >= minVal, nil
		}),
	},
	"max": {
		name:          "max",
		supportedType: []reflect.Kind{reflect.Int},
		checkFunc: CheckRule(func(value any, args ...string) (bool, error) {
			if len(args) != 1 {
				return false, fmt.Errorf("regexp validator: must be 1 int arg")
			}
			intVal, ok := value.(int)
			if !ok {
				return false, fmt.Errorf("regexp validator: value must be int")
			}
			maxVal, err := strconv.Atoi(args[0])
			if err != nil {
				return false, fmt.Errorf("regexp validator: arg must be int")
			}
			return intVal <= maxVal, nil
		}),
	},
}

func parseTag(tagValue string) []ValidatorTag {
	result := []ValidatorTag{}
	for _, rule := range strings.Split(tagValue, "|") {
		if len(rule) > 0 {
			parts := strings.Split(rule, ":")
			if len(parts) > 0 {
				if parts[0] == "regexp" {
					result = append(result, ValidatorTag{name: parts[0], args: []string{parts[1]}})
				} else {
					result = append(result, ValidatorTag{name: parts[0], args: strings.Split(parts[1], ",")})
				}
			}
		}
	}
	return result
}

func createError(val any, tag ValidatorTag) error {
	return fmt.Errorf(fmt.Sprintf("Value '%v' is invalid (by validator '%v')", val, tag))
}

func validate(fieldName string, val reflect.Value, validatorTags []ValidatorTag) error {
	errors := ValidationErrors{}
	switch kind := val.Kind(); kind {
	case reflect.Int, reflect.String:
		for _, tag := range validatorTags {
			validator := validators[tag.name]
			if validator.isSupported(kind) {
				isValid, err := validator.checkFunc(val.Interface(), tag.args...)
				if err != nil {
					return err
				}
				if !isValid {
					errors = append(errors, ValidationError{
						Field: fieldName,
						Err:   createError(val.Interface(), tag),
					})
				}
			}
		}
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			if err := validate(fieldName, val.Index(i), validatorTags); err != nil {
				switch e := err.(type) {
				case ValidationErrors:
					errors = append(errors, e...)
				default:
					return err
				}
			}
		}
	default:
		return nil
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func Validate(v interface{}) error {
	if v == nil {
		return fmt.Errorf("input argument error: nil value")
	}
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("input argument error: can validate only struct type")
	}

	errors := make(ValidationErrors, 0, 10)

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := val.Type().Field(i)
		fieldName := val.Type().Field(i).Name
		if fieldTyp.IsExported() {
			tag := fieldTyp.Tag
			if validateTag := tag.Get("validate"); validateTag != "" {
				validatorTags := parseTag(validateTag)
				if err := validate(fieldName, fieldVal, validatorTags); err != nil {
					switch e := err.(type) {
					case ValidationErrors:
						errors = append(errors, e...)
					default:
						return err
					}
				}
			}
		}
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}
