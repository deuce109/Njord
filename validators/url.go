package validators

import (
	"errors"
	"reflect"
	"regexp"
)

func IsPath(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return errors.New("isPath is only useable on strings")
	}
	match, err := regexp.MatchString(`((?:[^/]*/)*)(.+)`, st.String())
	if err != nil {
		// Untested due to only failing when the regex compile fails
		return err
	} else if !match {
		return errors.New("string does not match URI regex")
	}

	return nil
}
