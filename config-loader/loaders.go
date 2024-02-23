package configLoader

import (
	"os"
	"reflect"
	"strconv"
)

type ApplyFunc func(fieldType reflect.StructField, fieldValue reflect.Value, tag tagData) error

func nilApplyFunc(fieldType reflect.StructField, fieldValue reflect.Value, tag tagData) error {
	return nil
}

func MapLoader(configFile map[string]any) ApplyFunc {
	return func(fieldType reflect.StructField, fieldValue reflect.Value, tag tagData) error {
		// Set the value
		switch fieldType.Type.Kind() {
		case reflect.String:
			value, err := traverseMap[string](configFile, tag.key, tag.path...)
			if err != nil {
				if _, isPresent := err.(*KeyNotPresent); isPresent {
					return err
				} else {
					return err
				}
			}
			fieldValue.SetString(value)
			err = nil
		case reflect.Int:
			value, err := traverseMap[int](configFile, tag.key, tag.path...)
			if err != nil {
				if _, isPresent := err.(*KeyNotPresent); isPresent {
					return err
				} else {
					return err
				}
			}
			fieldValue.SetInt(int64(value))
			err = nil
		}
		return nil
	}
}

func FileLoader(configFile string, ignoreMissing bool) ApplyFunc {
	file, err := loadYamlFile(configFile)
	if err != nil {
		if ignoreMissing {
			return nilApplyFunc
		} else {
			panic(err)
		}
	}
	return MapLoader(file)

}

func EnvironmentLoader(envPrefix string) ApplyFunc {
	return func(fieldType reflect.StructField, fieldValue reflect.Value, tag tagData) error {
		// Read the environment variables
		envName := getEnvName(tag.config, envPrefix)
		envVal, ok := os.LookupEnv(envName)
		if !ok {
			return nil
		}
		switch fieldType.Type.Kind() {
		case reflect.String:
			fieldValue.SetString(envVal)
		case reflect.Int:
			envValInt, err := strconv.Atoi(envVal)
			if err != nil {
				return err
			}
			fieldValue.SetInt(int64(envValInt))
		}
		return nil
	}
}
