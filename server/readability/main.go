package readability

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"pagemail/server/models"
	"pagemail/server/net"
	"path/filepath"

	"github.com/pocketbase/pocketbase"
)

type ReaderConfig struct {
	NodeScript   string
	PythonScript string
	ContextDir   string
}

func (r *ReaderConfig) GetContextDir() string {
	ctxPath, err := filepath.Abs(r.ContextDir)
	if err != nil {
		log.Panic(err)
	}
	return ctxPath
}

func StartReaderTask(app *pocketbase.PocketBase, record *models.Page, cfg ReaderConfig) (*models.SynthesisTask, error) {
	// Get the URL and invoke the pipeline
	url := record.Url
	log.Print("Fetching url")
	buf, err := net.FetchUrlContents(url)

	log.Print("Starting reader task")
	task_data := new(models.SynthesisTask)
	raw_out, err := doReaderTask(cfg, url, buf)
	if err != nil {
		return nil, err
	}

	log.Print("Unmarshalling json")
	err = json.Unmarshal(raw_out, task_data)
	if err != nil {
		return nil, err
	}

	return task_data, nil
}

func CheckIsReadable(cfg ReaderConfig, url string, contents []byte) bool {
	ctxPath := cfg.GetContextDir()
	check_tsk := exec.Command("node", cfg.NodeScript, "--check", "--url", url)
	check_tsk.Dir = filepath.Join(ctxPath, "/dist")
	check_tsk.Stdin = bytes.NewReader(contents)
	check_tsk.Start()

	return check_tsk.Wait() == nil
}

func insertHeader(data *[]byte) *[]byte {
	contentSize := len(*data)
	newRef := make([]byte, contentSize+4)
	binary.BigEndian.PutUint32(newRef[0:4], uint32(contentSize))
	copy(newRef[4:], *data)
	return &newRef
}

func doReaderTask(cfg ReaderConfig, url string, contents []byte) ([]byte, error) {
	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()
	out := new(bytes.Buffer)

	ctxPath := cfg.GetContextDir()
	document_tsk := exec.Command("node", "main.js", "--url", url)
	document_tsk.Dir = filepath.Join(ctxPath, "dist")
	document_tsk.Stdout = w
	document_tsk.Stdin = bytes.NewReader(contents)

	parser_tsk := exec.Command("venv/bin/python3", "test.py")
	parser_tsk.Dir = ctxPath
	parser_tsk.Stdin = r
	parser_tsk.Stdout = out

	err := document_tsk.Start()
	if err != nil {
		return nil, fmt.Errorf("Failed to start node task with error: %s", err)
	}

	err = parser_tsk.Start()
	if err != nil {
		return nil, fmt.Errorf("Failed to start python task with error: %s", err)
	}

	err = document_tsk.Wait()
	if err != nil {
		return nil, fmt.Errorf("Node task exited with error: %s", err)
	}

	w.Close()

	err = parser_tsk.Wait()
	if err != nil {
		return nil, fmt.Errorf("Node task exited with error: %s", err)
	}
	log.Printf("Completed reader task for %s", url)
	return out.Bytes(), nil
}
