package readability

import (
	"encoding/binary"
	"encoding/json"
	"pagemail/server/models"
	"testing"
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

func getConfig() ReaderConfig {
	return ReaderConfig{
		NodeScript: "main.js",
		PythonScript: "test.js",
		ContextDir: "/Users/ellis/Git/pagemail/readability/",
	}
}

func TestPipeline(t *testing.T) {
	// url := "https://developer.mozilla.org/en-US/docs/Web/HTML/Element/audio"
	// url := "https://www.allthingsdistributed.com/2023/07/building-and-operating-a-pretty-big-storage-system.html"
	url := "https://www.example.com"
	data := getExample()
	config := getConfig()
	out, err := doReaderTask(config, url, data)
	if err != nil {
		t.Errorf("Failed test: %s", err)
	}
	taskout := new(models.ReadabilityResponse)
	err = json.Unmarshal(out, taskout)
	if err != nil || taskout == nil {
		t.Errorf("Failed to marshall output: %s", err)
	}
	t.Log(string(out))
	t.Log(taskout)
}

func TestCheck(t *testing.T) {
	// url := "https://developer.mozilla.org/en-US/docs/Web/HTML/Element/audio"
	// url := "https://www.allthingsdistributed.com/2023/07/building-and-operating-a-pretty-big-storage-system.html"
	url := "https://www.google.com"
	data := getExample()
	cfg := getConfig()
	is_readable := CheckIsReadable(cfg, url, data)
	t.Logf("Completed with result %t", is_readable)
}

func TestHeaderAdd(t *testing.T) {
	data := []byte{0x00, 0xFF, 0x1c}
	out := insertHeader(&data)
	if len(*out) != len(data)+4 {
		t.FailNow()
	}
	headerVal := binary.BigEndian.Uint32((*out)[:4])
	if int(headerVal) != len(data) {
		t.FailNow()
	}
	t.Log(out)
}
