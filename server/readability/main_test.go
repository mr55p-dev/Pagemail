package readability

import (
	"encoding/json"
	"io"
	"net/http"
	"pagemail/server/models"
	"testing"
)

func TestPipeline(t *testing.T) {
	// url := "https://developer.mozilla.org/en-US/docs/Web/HTML/Element/audio"
	url := "https://www.allthingsdistributed.com/2023/07/building-and-operating-a-pretty-big-storage-system.html"
	contents, err := http.Get(url)
	if err != nil || contents.StatusCode != 200 {
		t.Errorf("Fetching test url failed, %d %s", contents.StatusCode, err)
	}
	buf, err := io.ReadAll(contents.Body)

	out, err := doReaderTask(url, &buf)
	if err != nil {
		t.Errorf("Failed test: %s", err)
	}
	taskout := new(models.SynthesisTask)
	err = json.Unmarshal(out, taskout)
	if err != nil {
	t.Errorf("Failed to marshall output: %s", err)
	}
	t.Log(taskout.TaskId)
}

func TestCheck(t *testing.T) {
	url := "https://developer.mozilla.org/en-US/docs/Web/HTML/Element/audio"
	// url := "https://www.allthingsdistributed.com/2023/07/building-and-operating-a-pretty-big-storage-system.html"
	contents, err := http.Get(url)
	if err != nil || contents.StatusCode != 200 {
		t.Errorf("Fetching test url failed, %d, %s", contents.StatusCode, err)
	}

	buf, err := io.ReadAll(contents.Body)

	is_readable := CheckIsReadable(url, &buf)
	t.Logf("Completed with result %t", is_readable)
}
