package common

import (
	"context"
	"encoding/json"
	"fmt"

	"os"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/datatypes"
)

const (
	ERR_PARSE_VALUE_ENV = "cannot parse value of %v env"
)

var (
	FormatErr = func(prefix string, err error) error {
		return ErrorWrapper(fmt.Sprintf(ERR_PARSE_VALUE_ENV, prefix), err)
	}
)

func ConvertStruct2Map(ctx context.Context, obj interface{}) map[string]interface{} {
	m := map[string]interface{}{}
	if obj != nil {
		values := reflect.ValueOf(obj).Elem()
		types := values.Type()

		for i := 0; i < values.NumField(); i++ {
			key := types.Field(i).Tag.Get("filter")
			fmt.Sprint(key)
			value := values.Field(i)

			if key != "" && !value.IsNil() {
				m[key] = value.Interface()
			}
		}
	}

	return m
}

var ConvertMap2StringSQL = func(cond map[string]interface{}) ([]string, []interface{}) {
	sqls := []string{}
	values := []interface{}{}

	for k, v := range cond {
		operator := "="
		if k != "" && v != nil {
			typeValue := fmt.Sprintf("%T", v)
			if strings.Contains(typeValue, "[]") {
				operator = "IN"
			}
			sqls = append(sqls, fmt.Sprintf("%s %s ?", k, operator))

			values = append(values, v)
		}
	}

	return sqls, values
}

type osENV struct {
	name  string
	value string
}

func (o *osENV) ParseInt() (value int64, err error) {
	v, err := strconv.ParseInt(o.value, 10, 64)
	if err != nil {
		return v, FormatErr(o.name, err)
	}
	return v, nil
}

func (o *osENV) ParseUInt() (value uint64, err error) {
	v, err := strconv.ParseUint(o.value, 10, 64)
	if err != nil {
		return v, FormatErr(o.name, err)
	}
	return v, nil
}

func (o *osENV) ParseString() (value string, err error) {
	return o.value, nil
}

func (o *osENV) ParseBool() (value bool, err error) {
	v, err := strconv.ParseBool(o.value)
	if err != nil {
		return v, FormatErr(o.name, err)
	}
	return v, nil
}

func (o *osENV) ParseFloat() (value float64, err error) {
	v, err := strconv.ParseFloat(o.value, 64)
	if err != nil {
		return v, FormatErr(o.name, err)
	}
	return v, nil
}

func GetOSEnv(envName string) *osENV {
	value := os.Getenv(envName)
	return &osENV{name: envName, value: value}
}

func GetOffset(page int, pageSize int) int {
	return pageSize * (page - 1)
}

func UnmarshalJSON(input string) (datatypes.JSON, error) {
	var result datatypes.JSON
	if input != "" {
		err := json.Unmarshal([]byte(input), &result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func Contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
