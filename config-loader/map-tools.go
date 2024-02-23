package configLoader

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func traverseMap[T any](target map[string]any, key string, segments ...string) (T, error) {
	// Traverse the config file
	zero := *new(T)
	configFileKey := target
	for _, segment := range segments {
		var ok bool
		configFileKey, ok = configFileKey[segment].(map[string]any)
		if !ok {
			return zero, errKeyNotPresent(key)
		}
	}
	value, ok := configFileKey[key]
	if !ok {
		return zero, errKeyNotPresent(key)
	}
	castedValue, ok := value.(T)
	if !ok {
		return zero, fmt.Errorf("Invalid type for variable %v", value)
	}
	return castedValue, nil
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

