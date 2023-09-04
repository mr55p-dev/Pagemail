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
	log.Print("Starting reader task")
	// Get the URL and invoke the pipeline
	url := record.Url
	buf, err := net.FetchUrlContents(url)
	log.Print("Fetched url contents")

	task_data := new(polly.StartSpeechSynthesisTaskOutput)
	log.Printf("Starting reader task for %s", url)
	rawOut, err := doReaderTask(cfg, url, buf)
	if err != nil {
		log.Print("Error executing reader task", err)
		return nil, err
	}

	err = json.Unmarshal(rawOut, task_data)
	if err != nil {
		log.Print("Error parsing task output", err)
		return nil, err
	}

	log.Printf("Synthesis task obtained with id %s", *task_data.SynthesisTask.TaskId)
	c, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	taskData, taskErr := AwaitJobCompletion(c, cfg, task_data.SynthesisTask.TaskId)
	go func() {
		select {
		case job := <-taskData:
			log.Print("Got some output", job)
			status := job.SynthesisTask.TaskStatus
			UpdateJobState(app, record.Id, models.ReadabilityFromPolly(&status), task_data)
		case err := <-taskErr:
			log.Print("Got an error", err)
			UpdateJobState(app, record.Id, models.ReadabilityFailed, task_data)
		}
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
	stdOut := new(bytes.Buffer)
	stdErr := new(bytes.Buffer)
	log.Print("Entered doReaderTask")
	in := insertHeader(contents)
	log.Print("Header inserted")

	rCfg := cfg.Readability
	ctxPath := rCfg.GetContextDir()
	log.Print(rCfg)
	document_tsk := exec.Command("node", rCfg.NodeScript, "--url", url)
	document_tsk.Dir = ctxPath
	document_tsk.Stdout = w
	document_tsk.Stdin = bytes.NewReader(in)
	document_tsk.Stderr = stdErr

	parser_tsk := exec.Command("venv/bin/python3", rCfg.PythonScript)
	parser_tsk.Dir = ctxPath
	parser_tsk.Stdin = r
	parser_tsk.Stdout = stdOut
	parser_tsk.Stderr = stdErr
	log.Print("Setup tasks")

	err := document_tsk.Start()
	if err != nil {
		log.Printf("Written to stderr: %s", stdErr.String())
		return nil, fmt.Errorf("Failed to start node task with error: %s", err)
	}
	log.Print("Started nodejs")

	err = parser_tsk.Start()
	if err != nil {
		log.Printf("Written to stderr: %s", stdErr.String())
		return nil, fmt.Errorf("Failed to start python task with error: %s", err)
	}
	log.Print("Started python")

	err = document_tsk.Wait()
	if err != nil {
		log.Printf("Written to stderr: %s", stdErr.String())
		return nil, fmt.Errorf("Node task exited with error: %s", err)
	}
	log.Print("Finished nodejs")

	w.Close()
	log.Print("Closed write end of pipe")

	err = parser_tsk.Wait()
	if err != nil {
		log.Printf("Written to stderr: %s", stdErr.String())
		return nil, fmt.Errorf("Python task exited with error: %s", err)
	}
	log.Printf("Completed reader tasks succesfully for %s", url)
	return stdOut.Bytes(), nil
}
