package helpers

import "reflect"

func IsFieldSet(data *reflect.Value, fieldName string) bool {
	field := data.FieldByName(fieldName)
	return field.IsValid() && !field.IsZero()
}
