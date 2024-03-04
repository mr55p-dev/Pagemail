package configLoader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraverseMap(t *testing.T) {
	assert := assert.New(t)
	testMap := map[string]any{
		"key": "value",
		"nested_key_1": map[string]any{
			"nested_key_2": map[string]any{
				"nested_value": "nested_value",
				"nested_int":   2,
			},
		},
	}
	var val string
	err := traverseMap[string](
		&val,
		testMap,
		"nested_value",
		"nested_key_1",
		"nested_key_2",
	)
	assert.Nil(err, "there should be no error")
	assert.Equal("nested_value", val, "key should be nested_value")

	var intVal int
	err = traverseMap[int](
		&intVal,
		testMap,
		"nested_int",
		"nested_key_1",
		"nested_key_2",
	)
	assert.Nil(err)
	assert.Equal(intVal, 2)
}

func TestNoKeyTraverseMap(t *testing.T) {
	assert := assert.New(t)
	testMap := map[string]any{
		"key": "value",
		"nested_key_1": map[string]any{
			"nested_key_2": map[string]any{
				"nested_int": 2,
			},
		},
	}
	var val string
	err := traverseMap[string](
		&val,
		testMap,
		"nested_value",
		"nested_key_1",
		"nested_key_2",
	)
	assert.ErrorContains(err, "nested_value")
	assert.IsType(&KeyNotPresent{}, err)
	assert.Zero(val)
}

func TestGetEnvName(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]string{
		"key":             "KEY",
		"key-with-hyphen": "KEY_WITH_HYPHEN",
		"key_with-mixed":  "KEY_WITH_MIXED",
		"path.key-name":   "PATH_KEY_NAME",
	}
	for k, v := range tests {
		assert.Equal(v, getEnvName(k, ""))
	}
}

func TestGetEnvNameWithPrefix(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]string{
		"key":           "XYZ_KEY",
		"path.key-name": "XYZ_PATH_KEY_NAME",
	}
	for k, v := range tests {
		assert.Equal(v, getEnvName(k, "xyz"))
	}
}

func TestParseTag(t *testing.T) {
	assert := assert.New(t)
	config := "path.segment.key"
	tag := parseConfigTag(config)
	assert.Equal("key", tag.key)
	assert.Equal([]string{"path", "segment"}, tag.path)
	assert.Equal(config, tag.config)
	assert.False(tag.options.optional)
}

func TestParseTagOptional(t *testing.T) {
	assert := assert.New(t)
	config := "path.segment.key,optional"
	tag := parseConfigTag(config)
	assert.Equal("key", tag.key)
	assert.Equal([]string{"path", "segment"}, tag.path)
	assert.Equal(config, tag.config)
	assert.True(tag.options.optional)
}
