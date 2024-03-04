package configLoader

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func traverseMap[T string | int](dest *T, target map[string]any, key string, segments ...string) error {
	// Traverse the config file
	configFileKey := target
	for _, segment := range segments {
		var ok bool
		configFileKey, ok = configFileKey[segment].(map[string]any)
		if !ok {
			return errKeyNotPresent(key)
		}
	}
	value, ok := configFileKey[key]
	if !ok {
		return errKeyNotPresent(key)
	}
	castedValue, ok := value.(T)
	if !ok {
		return fmt.Errorf("Invalid type for variable %v", value)
	}
	(*dest) = castedValue
	return nil
}

func loadYamlFile(filename string) (map[string]any, error) {
	out := make(map[string]any)
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
