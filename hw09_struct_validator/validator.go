package hw09structvalidator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInputIsNotStruct = errors.New("input value is not a struct")
	ErrInvalidLength    = errors.New("length is invalid")
	ErrNotMatchRegexp   = errors.New("string is not matched by regexp")
	ErrNotIncludedInSet = errors.New("not included in validation set")
	ErrLessThanMin      = errors.New("less than the minimum")
	ErrMaxMoreMax       = errors.New("more than maximum")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	for _, e := range v {
		builder.WriteString("field: ")
		builder.WriteString(e.Field)
		builder.WriteString(" - ")
		builder.WriteString(e.Err.Error())
	}
	return builder.String()
}

type validator struct {
	errs ValidationErrors
}

func Validate(v interface{}) error {
	val := &validator{}
	return val.validate(v)
}

func (v *validator) validate(vv interface{}) error {
	value := reflect.ValueOf(vv)

	if value.Kind() != reflect.Struct {
		v.errs = append(v.errs, ValidationError{
			Field: "",
			Err:   ErrInputIsNotStruct,
		})
		return v.errs
	}

	valueType := value.Type()
	for i := 0; i < valueType.NumField(); i++ {
		fieldTypeValue := valueType.Field(i)
		structTags := fieldTypeValue.Tag.Get("validate")
		if len(structTags) == 0 {
			continue
		}
		fieldValue := value.Field(i)

		if !fieldValue.CanInterface() {
			continue
		}
		tags := strings.Split(structTags, "|")

		switch fieldValue.Kind() {
		case reflect.String:
			v.validateString(tags, fieldTypeValue.Name, fieldValue.String())

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.validateInteger(tags, fieldTypeValue.Name, fieldValue.Int())

		case reflect.Slice:
			switch fieldValue.Type().String() {
			case "[]string":
				for i := 0; i < fieldValue.Len(); i++ {
					v.validateString(tags, fieldTypeValue.Name, fieldValue.Index(i).String())
				}
			case "[]int", "[]int8", "[]int16", "[]in32", "[]int64":

				for i := 0; i < fieldValue.Len(); i++ {
					v.validateInteger(tags, fieldTypeValue.Name, fieldValue.Index(i).Int())
				}
			}
		}
	}
	if len(v.errs) > 0 {
		return v.errs
	}
	return nil
}

func (v *validator) validateString(tags []string, field string, value string) {
	for _, tag := range tags {
		t := strings.Split(tag, ":")
		switch t[0] {
		case "len":
			cond, err := strconv.Atoi(t[1])
			if err != nil {
				log.Printf("len value is not int")
				continue
			}
			if cond != len(value) {
				v.errs = append(v.errs, ValidationError{
					Field: field,
					Err:   fmt.Errorf("%w: need %d, now %d", ErrInvalidLength, cond, len(value)),
				})
			}
		case "regexp":
			rg, err := regexp.Compile(t[1])
			if err != nil {
				log.Printf("regexp invalid: %s", err)
				continue
			}
			matched := rg.MatchString(value)
			if !matched {
				v.errs = append(v.errs, ValidationError{
					Field: field,
					Err:   fmt.Errorf("%w: %s", ErrNotMatchRegexp, t[1]),
				})
			}
		case "in":
			set := strings.Split(t[1], ",")
			inSet := false
			for _, s := range set {
				if s == value {
					inSet = true
					break
				}
			}
			if !inSet {
				v.errs = append(v.errs, ValidationError{
					Field: field,
					Err:   fmt.Errorf("%w: value %s, set %v", ErrNotIncludedInSet, value, set),
				})
			}
		default:
			log.Printf("unknown validator's name %s", t[0])
		}
	}
}

func (v *validator) validateInteger(tags []string, field string, value int64) {
	for _, tag := range tags {
		t := strings.Split(tag, ":")
		switch t[0] {
		case "min":
			cond, err := strconv.Atoi(t[1])
			if err != nil {
				log.Printf("min value is not int")
			}
			if int(value) < cond {
				v.errs = append(v.errs, ValidationError{
					Field: field,
					Err:   fmt.Errorf("%w: value %d, condition %d", ErrLessThanMin, value, cond),
				})
			}
		case "max":
			cond, err := strconv.Atoi(t[1])
			if err != nil {
				log.Printf("min value is not int")
			}
			if int(value) > cond {
				v.errs = append(v.errs, ValidationError{
					Field: field,
					Err:   fmt.Errorf("%w: value %d, condition %d", ErrMaxMoreMax, value, cond),
				})
			}
		case "in":
			set := strings.Split(t[1], ",")
			inSet := false
			for _, s := range set {
				intVal, err := strconv.Atoi(s)
				if err != nil {
					log.Printf("in value is not int")
				}
				if intVal == int(value) {
					inSet = true
					break
				}
			}
			if !inSet {
				v.errs = append(v.errs, ValidationError{
					Field: field,
					Err:   fmt.Errorf("%w: value %d, set %v", ErrNotIncludedInSet, value, set),
				})
			}
		default:
			log.Printf("unknown validator's name %s", t[0])
		}
	}
}
