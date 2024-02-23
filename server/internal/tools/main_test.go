package tools

import "testing"

func TestGenerateId(t *testing.T) {
	vals := []int{1, 5, 100}
	for _, v := range vals {
		id := GenerateNewId(v)
		if len(id) != v {
			t.Log("Requested", v, "got", id, len(id))
			t.FailNow()
		}
		t.Log("Requested", v, "chars:", id)
	}
}
