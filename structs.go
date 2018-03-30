package structs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func StructFromEnv(i interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if tagValue := fi.Tag.Get("json"); tagValue != "" {
			value := os.Getenv(strings.ToUpper(tagValue))
			if value == "" {
				return fmt.Errorf("missing env var value: %s", strings.ToUpper(tagValue))
			}
			v.FieldByName(fi.Name).SetString(value)
		}
	}
	return nil
}

func StructFromFile(i interface{}, envVar string) error {
	fPath := os.Getenv(envVar)
	if fPath != "" {
		raw, err := ioutil.ReadFile(fPath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(raw, i)
		if err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("%v env var is not set", envVar)
	}
}

func ToMap(in interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if tagValue := fi.Tag.Get("json"); tagValue != "" {
			if v.Field(i).Type() == reflect.TypeOf(true) {
				out[strings.ToUpper(tagValue)] = strconv.FormatBool(v.Field(i).Bool())
			} else {
				out[strings.ToUpper(tagValue)] = v.Field(i).String()
			}
		}
	}
	return out, nil
}

func Append(obj interface{}, config map[string]string) (map[string]string, error) {
	mMap, err := ToMap(obj)
	if err != nil {
		return nil, err
	}
	for key, value := range mMap {
		config[key] = value.(string)
	}
	return config, nil
}
