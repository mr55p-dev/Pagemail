package readability

import (
	"encoding/json"
	"io"
	"os/exec"
	"pagemail/server/models"
)

type SynthesisTask struct {
	Engine            string `json:"engine"`
	TaskId            string `json:"taskId"`
	TaskStatus        string `json:"taskStatus"`
	OutputUri         string `json:"outputUri"`
	CreationTime      string `json:"creationTime"`
	RequestCharacters int    `json:"requestCharacters"`
	OutputFormat      string `json:"outputFormat"`
	TextType          string `json:"textType"`
	VoiceId           string `json:"voiceId"`
	LanguageCode      string `json:"languageCode"`
}

func StartReaderTask(record *models.PageRecord) (*SynthesisTask, error) {
	// Get the URL and invoke the pipeline
	url := record.Url

	task_data := new(SynthesisTask)
	raw_out, err := doReaderTask(url)
	if err != nil {
		return task_data, err
	}

	err = json.Unmarshal(raw_out, task_data)
	if err != nil {
		return task_data, err
	}

	return task_data, nil
}

func doReaderTask(url string) ([]byte, error) {
	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()

	document_tsk := exec.Command("node", "main.js", url)
	document_tsk.Dir = "/Users/ellis/Git/pagemail/readability/dist"
	document_tsk.Stdout = w
	err := document_tsk.Start()
	if err != nil {
		return []byte{}, err
	}

	parser_tsk := exec.Command("venv/bin/python3", "test.py")
	parser_tsk.Dir = "/Users/ellis/Git/pagemail/readability"
	parser_tsk.Stdin = r
	raw_output, err := parser_tsk.Output()
	if err != nil {
		return []byte{}, err
	}
	return raw_output, nil
}

// func ReaderTaskStatus(task_id string) (*Reader, error)
