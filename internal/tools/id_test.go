package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateId(t *testing.T) {
	assert := assert.New(t)
	vals := []int{1, 5, 8, 20, 25, 50, 75, 100}
	outs := make([]string, 0, len(vals))
	for _, v := range vals {
		id := GenerateNewId(v)
		outs = append(outs, id)
		assert.Len(id, v)
		t.Log(id)
	}

	for idx, id := range outs {
		for j := idx + 1; j < len(outs); j++ {
			assert.NotEmpty(id, outs[j], "Collision in ids")
		}
	}
}
