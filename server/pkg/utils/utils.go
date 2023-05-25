package utils

import (
	"encoding/json"
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

func Map[T any](s interface{}) *T {
	var target T
	ms, _ := json.Marshal(s)
	json.Unmarshal(ms, &target)
	return &target
}

func GetType(t interface{}) string {
	valueOf := reflect.ValueOf(t)

	if valueOf.Type().Kind() == reflect.Ptr {
		return reflect.Indirect(valueOf).Type().Name()
	} else {
		return valueOf.Type().Name()
	}
}

func GetFieldValue(v interface{}, field string) interface{} {
	valueOf := reflect.ValueOf(v)
	if valueOf.Type().Kind() == reflect.Ptr {
		return reflect.Indirect(valueOf).FieldByName(field)
	} else {
		return valueOf.FieldByName(field)
	}
}
