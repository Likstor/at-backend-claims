package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

func Load(cfg any) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("waiting pointer to struct")
	}
	v = v.Elem()
	t := v.Type()

	var err error

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if !value.CanSet() {
			continue
		}

		envKey := field.Tag.Get("env")
		if envKey == "" {
			continue
		}

		envStr, exists := os.LookupEnv(envKey)

		if !exists {
			err = errors.Join(err, fmt.Errorf("%s is missing or empty", envKey))
			continue
		}

		if err := setValue(value, envStr); err != nil {
			err = errors.Join(err, fmt.Errorf("%s error: %v", envKey, err.Error()))
		}
	}

	return err
}

func setValue(v reflect.Value, s string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := time.ParseDuration(s)
			if err != nil {
				return err
			}
			v.SetInt(int64(d))
			return nil
		}

		i, err := strconv.ParseInt(s, 10, v.Type().Bits())
		if err != nil {
			return err
		}

		v.SetInt(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(s, 10, v.Type().Bits())
		if err != nil {
			return err
		}

		v.SetUint(i)

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, v.Type().Bits())
		if err != nil {
			return err
		}

		v.SetFloat(f)

	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}

		v.SetBool(b)

	default:
		return fmt.Errorf("unknown type: %v", v.Kind())
	}

	return nil
}
