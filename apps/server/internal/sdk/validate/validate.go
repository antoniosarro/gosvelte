package validate

import (
	"errors"
	"reflect"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

var (
	Timestamp = validation.By(func(value interface{}) error {
		s, _ := value.(string)
		if _, err := time.Parse(time.RFC3339, s); err != nil {
			return errors.New("invalid timestamp")
		}
		return nil
	})
)

func SplitErrors(err error) []string {
	errStr := strings.TrimSuffix(err.Error(), ".")

	if reflect.TypeOf(err) == reflect.TypeOf(validation.Errors{}) {
		return strings.Split(errStr, ";")
	}
	return nil
}
