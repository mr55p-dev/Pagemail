package configLoader

import (
	"fmt"
	"reflect"
	"strings"
)

type tagOptions struct {
	optional bool
}

type tagData struct {
	config  string
	key     string
	path    []string
	options tagOptions
}

func parseTagOptions(opts []string) tagOptions {
	out := tagOptions{}
	for _, v := range opts {
		switch v {
		case "optional":
			out.optional = true
		}
	}
	return out
}

func parseConfigTag(config string) tagData {
	data := tagData{
		config: config,
	}
	v := strings.Split(config, ",")
	if len(v) > 1 {
		data.options = parseTagOptions(v[1:])
	}

	path := strings.Split(v[0], ".")
	data.key = path[len(path)-1]
	data.path = path[:len(path)-1]
	return data
}

func LoadConfig(dest any, loaders ...ApplyFunc) error {
	var err error
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("Panic generated while loading config: %v", v)
		}
	}()
	valueOf := reflect.ValueOf(dest).Elem()
	typeOf := valueOf.Type()
	for i := 0; i < typeOf.NumField(); i++ {
		fieldType := typeOf.Field(i)
		fieldValue := valueOf.Field(i)
		config := fieldType.Tag.Get("config")
		if config == "" {
			continue
		}
		tag := parseConfigTag(config)

		isLoaded := false
		for _, loader := range loaders {
			err = loader(fieldType, fieldValue, tag)
			if err != nil {
				if _, keyMissing := err.(*KeyNotPresent); keyMissing {
					continue
				} else {
					return fmt.Errorf("Error loading %s: %s", config, err.Error())
				}
			}
			isLoaded = true
		}

		if err != nil { // check any errors after loading
			_, keyMissing := err.(*KeyNotPresent)
			if keyMissing {
				if isLoaded { // we have already loaded the key
					continue
				} else if tag.options.optional { // field is marked as optional
					continue
				}
			}
			return err // all other cases should fall through and error
		}
	}
	return nil
}
