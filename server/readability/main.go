package readability

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"pagemail/server/models"
	"pagemail/server/net"

	"github.com/pocketbase/pocketbase"
)

func StartReaderTask(app *pocketbase.PocketBase, record *models.Page) (*models.SynthesisTask, error) {
	// Get the URL and invoke the pipeline
	url := record.Url
	log.Print("Fetching url")
	buf, err := net.FetchUrlContents(url)

	log.Print("Starting reader task")
	task_data := new(models.SynthesisTask)
	raw_out, err := doReaderTask(url, buf)
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

func insertHeader(data *[]byte) *[]byte {
	contentSize := len(*data)
	newRef := make([]byte, contentSize + 4)
	binary.BigEndian.PutUint32(newRef[0:4], uint32(contentSize))
	copy(newRef[4:], *data)
	return &newRef
}

func doReaderTask(url string, contents *[]byte) ([]byte, error) {
	if len(*contents) == 0 {
		log.Print("Zero bytes passed")
		return nil, nil
	}
	r, w := io.Pipe()
	defer r.Close()
	defer w.Close()
	out := new(bytes.Buffer)
	os.WriteFile("out", *contents, fs.FileMode(os.O_RDWR))

	log.Printf("node main.js --url %s", url)
	document_tsk := exec.Command("node", "main.js", "--url", url)
	document_tsk.Dir = "/Users/ellis/Git/pagemail/readability/dist"
	document_tsk.Stdout = w
	document_tsk.Stdin = bytes.NewReader(*contents)

	parser_tsk := exec.Command("venv/bin/python3", "test.py")
	parser_tsk.Dir = "/Users/ellis/Git/pagemail/readability"
	parser_tsk.Stdin = r
	parser_tsk.Stdout = out

	err := document_tsk.Start()
	if err != nil {
		log.Printf("Task1 start %s", err)
	}

	err = parser_tsk.Start()
	if err != nil {
		log.Printf("Task2 start %s", err)
	}

	err = document_tsk.Wait()
	log.Print("Finished 1")
	if err != nil {
		log.Printf("Task1 end %s", err)
	}

	w.Close()

	err = parser_tsk.Wait()
	log.Print("Finished 2")
	if err != nil {
		log.Printf("Errored with %s", err)
		return out.Bytes(), err
	}
	log.Print("Completed without error")
	return out.Bytes(), nil
}

func CheckIsReadable(url string, contents *[]byte) bool {
	check_tsk := exec.Command("node", "main.js", "--check", "--url", url)
	check_tsk.Dir = "/Users/ellis/Git/pagemail/readability/dist"
	check_tsk.Stdin = bytes.NewReader(*contents)
	check_tsk.Start()

	return c]heck_tsk.Wait() == nil
}
