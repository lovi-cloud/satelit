package testutils

import (
	"reflect"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// CompareStruct check Equal struct values.
// return []string is [unexpected Field, expected value, actual value]
func CompareStruct(expected, actual interface{}) (bool, []string) {
	ve := reflect.Indirect(reflect.ValueOf(expected))
	te := ve.Type()

	va := reflect.Indirect(reflect.ValueOf(actual))
	ta := va.Type()

	if !reflect.DeepEqual(te, ta) {
		return false, []string{"type", te.String(), ta.String()}
	}

	for i := 0; i < te.NumField(); i++ {
		fe := ve.Field(i)
		ie := fe.Interface()
		fa := va.Field(i)
		ia := fa.Interface()

		var value interface{}
		var ok bool
		if value, ok = ie.(uuid.UUID); ok {
			// this type is uuid.UUID
			if !reflect.DeepEqual(value, ia.(uuid.UUID)) {
				return false, []string{te.Field(i).Name, value.(uuid.UUID).String(), ia.(string)}
			}
		} else if value, ok = ie.(string); ok {
			// this type is string
			if !strings.EqualFold(value.(string), ia.(string)) {
				return false, []string{te.Field(i).Name, value.(string), ia.(string)}
			}
		}
	}

	return true, nil
}
