package util

import "reflect"

// CheckStructEmptyFields checks if a struct has empty fields
func CheckStructEmptyFields(s any) []string {
	var emptyFields []string
	excludedFields := []string{"state", "sizeCache", "unknownFields"}
	val := reflect.ValueOf(s).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		// Check if the field should be excluded
		if contains(excludedFields, fieldName) {
			continue
		}

		// Check if the field is a zero value or an empty string
		if field.Interface() == reflect.Zero(field.Type()).Interface() || (field.Kind() == reflect.String && field.Len() == 0) {
			emptyFields = append(emptyFields, fieldName)
		}
	}
	return emptyFields
}

func contains(slice []string, value string) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}
	return false
}
