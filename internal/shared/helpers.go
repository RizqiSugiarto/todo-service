package shared

import (
	"reflect"
	"time"

	"github.com/digisata/todo-service/internal/entity"
)

func ConvertToJakartaTime(t time.Time) time.Time {
	return t.Add(7 * time.Hour)
}

func CreateUpdateValueMap[T entity.UpdateTaskRequest | entity.UpdateActivityRequest | entity.UpdateTextRequest](req T) map[string]interface{} {
	updateValue := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}

	val := reflect.ValueOf(req)
	typ := reflect.TypeOf(req)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		dbTag := fieldType.Tag.Get("db")

		if dbTag == "" || dbTag == "-" {
			continue
		}

		// Check for pointer fields
		if field.Kind() == reflect.Ptr {
			if !field.IsNil() {
				updateValue[dbTag] = field.Elem().Interface()
			}
		} else {
			// Handle bool type separately
			if field.Kind() == reflect.Bool {
				updateValue[dbTag] = field.Bool()
			} else if !field.IsZero() {
				updateValue[dbTag] = field.Interface()
			}
		}
	}

	return updateValue
}
