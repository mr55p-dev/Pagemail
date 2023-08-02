package readability

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"
	"pagemail/server/models"
	"pagemail/server/net"

	"github.com/pocketbase/pocketbase"
)

func StartReaderTask(app *pocketbase.PocketBase, record *models.Page) (*models.SynthesisTask, error) {
	// Get the URL and invoke the pipeline
	url := record.Url
	body, err := net.FetchUrlContents(url)

	task_data := new(models.SynthesisTask)
	raw_out, err := doReaderTask(url, body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw_out, task_data)
	if err != nil {
		return nil, err
	}

	return task_data, nil
}

func doReaderTask(url string, contents *[]byte) ([]byte, error) {
	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()

	document_tsk := exec.Command("node", "main.js", "--url", url)
	document_tsk.Dir = "/Users/ellis/Git/pagemail/readability/dist"
	document_tsk.Stdout = w
	document_tsk.Stdin = bytes.NewReader(*contents)
	err := document_tsk.Start()
	defer document_tsk.Wait()
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

func CheckIsReadable(url string, contents *[]byte) bool {
	check_tsk := exec.Command("node", "main.js", "--check", "--url", url)
	check_tsk.Dir = "/Users/ellis/Git/pagemail/readability/dist"
	check_tsk.Stdin = bytes.NewReader(*contents)
	check_tsk.Start()

	return check_tsk.Wait() == nil
}
