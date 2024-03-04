package configLoader

import (
	"os"
	"reflect"
	"strconv"
)

type Loader func(fieldType reflect.StructField, fieldValue reflect.Value, tag tagData) error

func nilLoaderFn(fieldType reflect.StructField, fieldValue reflect.Value, tag tagData) error {
	return nil
}

func MapLoader(configFile map[string]any) Loader {
	return func(fieldType reflect.StructField, fieldValue reflect.Value, tag tagData) error {
		// Set the value
		switch fieldType.Type.Kind() {
		case reflect.String:
			var value string
			err := traverseMap[string](&value, configFile, tag.key, tag.path...)
			if err != nil {
				if _, isPresent := err.(*KeyNotPresent); isPresent {
					return err
				} else {
					return err
				}
			}
			fieldValue.SetString(value)
		case reflect.Int:
			var value int
			err := traverseMap[int](&value, configFile, tag.key, tag.path...)
			if err != nil {
				if _, isPresent := err.(*KeyNotPresent); isPresent {
					return err
				} else {
					return err
				}
			}
			fieldValue.SetInt(int64(value))
		}
		return nil
	}
}

func FileLoader(configFile string, ignoreMissing bool) Loader {
	file, err := loadYamlFile(configFile)
	if err != nil {
		if ignoreMissing {
			return nilLoaderFn
		} else {
			panic(err)
		}
	}
	return MapLoader(file)

}

func EnvironmentLoader(envPrefix string) Loader {
	return func(fieldType reflect.StructField, fieldValue reflect.Value, tag tagData) error {
		// Read the environment variables
		envName := getEnvName(tag.key, envPrefix)
		envVal, ok := os.LookupEnv(envName)
		if !ok {
			return &KeyNotPresent{"Key expected in variable " + envName}
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
