package readability

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"pagemail/server/models"
	"pagemail/server/net"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/pocketbase/pocketbase"
)

func StartReaderTask(app *pocketbase.PocketBase, cfg *models.PMContext, record *models.Page) (*polly.StartSpeechSynthesisTaskOutput, error) {
	// Get the URL and invoke the pipeline
	url := record.Url
	buf, err := net.FetchUrlContents(url)

	task_data := new(polly.StartSpeechSynthesisTaskOutput)
	raw_out, err := doReaderTask(cfg, url, buf)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	log.Print(string(raw_out))
	err = json.Unmarshal(raw_out, task_data)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	c, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	out := AwaitJobCompletion(c, task_data.SynthesisTask.TaskId)
	go func() {
		state := <-out
		status := state.status.SynthesisTask.TaskStatus
		UpdateJobState(app, record.Id, models.ReadabilityFromPolly(&status), task_data)
	}()

	return task_data, nil
}

func CheckIsReadable(ctx *models.PMContext, url string, contents []byte) bool {
	rCtx := ctx.Readability
	ctxPath := rCtx.GetContextDir()

	out_buf := new(bytes.Buffer)
	inp_with_header := insertHeader(contents)
	inp_buf := bytes.NewReader(inp_with_header)

	check_tsk := exec.Command("node", rCtx.NodeScript, "--check", "--url", url)
	check_tsk.Dir = ctxPath
	check_tsk.Stdin = inp_buf
	check_tsk.Stdout = out_buf

	check_tsk.Start()

	success := check_tsk.Wait() == nil
	return success
}

func insertHeader(data []byte) []byte {
	contentSize := len(data)
	newRef := make([]byte, contentSize+4)
	binary.BigEndian.PutUint32(newRef[0:4], uint32(contentSize))
	copy(newRef[4:], data)
	return newRef
}

func doReaderTask(cfg *models.PMContext, url string, contents []byte) ([]byte, error) {
	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()
	out := new(bytes.Buffer)
	input := insertHeader(contents)

	rCfg := cfg.Readability
	ctxPath := rCfg.GetContextDir()
	document_tsk := exec.Command("node", rCfg.NodeScript, "--url", url)
	document_tsk.Dir = ctxPath
	document_tsk.Stdout = w
	document_tsk.Stdin = bytes.NewReader(input)

	parser_tsk := exec.Command("venv/bin/python3", rCfg.PythonScript)
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
		return nil, fmt.Errorf("Python task exited with error: %s", err)
	}
	return out.Bytes(), nil
}
