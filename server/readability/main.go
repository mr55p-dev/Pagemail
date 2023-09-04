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

func getReadabilityData(cfg *models.ReaderConfig, url string, data io.Reader, buf io.Writer) error {
	log.Print("Fetching readability data for ", url)

	proc := exec.Command("node", cfg.NodeScript, "--url", url)
	proc.Dir = cfg.GetContextDir()
	proc.Stdout = buf
	proc.Stdin = data

	if err := proc.Run(); err != nil {
		err := err.(*exec.ExitError)

		stdoutData := new(bytes.Buffer)
		cnt, _ := io.Copy(buf, stdoutData)
		return fmt.Errorf(
			"Readability (node) exited with status %d: %s\nFollowing was written to stdout (%d bytes) %s",
			err.ExitCode(),
			err.Stderr,
			cnt,
			stdoutData.String(),
		)
	}
	log.Printf("Completed readability data job (process exited with status %d)", proc.ProcessState.ExitCode())
	return nil
}

func getSpeechSynthesisJob(cfg *models.ReaderConfig, data io.Reader, buf io.Writer) error {
	log.Print("Fetching synthesis job")

	proc := exec.Command("venv/bin/python3", cfg.PythonScript)
	proc.Dir = cfg.GetContextDir() 
	proc.Stdin = data
	proc.Stdout = buf

	if err := proc.Run(); err != nil {
		err := err.(*exec.ExitError)
		stdoutData := new(bytes.Buffer)
		cnt, _ := io.Copy(buf, stdoutData)
		return fmt.Errorf(
			"Readability (node) exited with status %d: %s\nFollowing was written to stdout (%d bytes) %s",
			err.ExitCode(),
			err.Stderr,
			cnt,
			stdoutData.String(),
		)
	}

	log.Printf("Completed readability synthesis job (process exited with status %d)", proc.ProcessState.ExitCode())
	return nil
}

func doReaderTask(cfg *models.PMContext, url string, contents []byte) ([]byte, error) {
	log.Printf("Starting reader task for %s", url)
	conf := cfg.Readability

	inputData := insertHeader(contents)
	inputStream := bytes.NewReader(inputData)
	outputStream := new(bytes.Buffer)
	log.Print("Prepared header bytes")

	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()

	if err := getReadabilityData(conf, url, inputStream, w); err != nil {
		log.Print(err)
		return nil, err
	}

	if err := getSpeechSynthesisJob(conf, r, outputStream); err != nil {
		log.Print(err)
		return nil, err
	}

	log.Printf("Completed reader tasks successfully for %s", url)
	return outputStream.Bytes(), nil
}
