package readability

import (
	"encoding/binary"
	"encoding/json"
	"pagemail/server/models"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/polly"
)

func getExample() []byte {
	return []byte(`
<html>
<head>
	<title>Example article</title>
</head>
<body>
	<article>
		<h1>This is a test article</h1>
		<section>
			<p>This is the entire article!</p>
		</section>
	</article>
</body>

</html>
	`)
}

func getConfig() models.ReaderConfig {
	return models.ReaderConfig{
		NodeScript: "main.js",
		PythonScript: "test.py",
		ContextDir: "../../readability/dist",
	}
}

func TestPipeline(t *testing.T) {
	url := "https://developer.mozilla.org/en-US/docs/Web/HTML/Element/audio"
	data := getExample()
	config := getConfig()
	out, err := doReaderTask(config, url, data)
	if err != nil {
		t.Errorf("Failed test: %s", err)
		t.Error(out)
		t.FailNow()
	}
	taskout := new(polly.StartSpeechSynthesisTaskOutput)
	err = json.Unmarshal(out, taskout)
	if err != nil || taskout == nil {
		t.Errorf("Failed to marshall output: %s", err)
	}
	t.Log(string(out))
	t.Log(taskout)
}

func TestCheck(t *testing.T) {
	url := "https://developer.mozilla.org/en-US/docs/Web/HTML/Element/audio"
	data := getExample()
	data = insertHeader(data)
	cfg := getConfig()
	is_readable := CheckIsReadable(cfg, url, data)
	t.Logf("Completed with result %t", is_readable)
}

func TestHeaderAdd(t *testing.T) {
	data := []byte{0x00, 0xFF, 0x1c}
	out := insertHeader(data)
	if len(out) != len(data)+4 {
		t.FailNow()
	}
	headerVal := binary.BigEndian.Uint32((out)[:4])
	if int(headerVal) != len(data) {
		t.FailNow()
	}
	t.Log(out)
}

func TestCrawlAll(t * testing.T) {
	// CrawlAll()
}
