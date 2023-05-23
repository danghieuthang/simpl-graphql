package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

// func (s *interface) MapArgsToStruct(m map[string]interface{}) error {
// 	for k, v := range m {
// 		err := SetField(s, k, v)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func Contains(s []interface{}, str string) bool {
	for _, v := range s {
		if strings.EqualFold(v.(string), str) {
			return true
		}
	}

	return false
}
