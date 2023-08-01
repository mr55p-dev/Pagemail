package readability

import (
	"encoding/json"
	"testing"
)

func TestPipeline(t *testing.T) {
	out, err := doReaderTask("https://developer.mozilla.org/en-US/docs/Web/HTML/Element/audio")
	if err != nil {
		t.Errorf("Failed test: %s", err)
	}
	taskout := new(SynthesisTask)
	err = json.Unmarshal(out, taskout)
	if err != nil {
		t.Errorf("Failed to marshall output: %s", err)
	}
	t.Log(taskout.TaskId)
}
