package readability

import (
	"context"
	"testing"
)

func TestFetchJob(t *testing.T) {
	res, err := FetchJobStatus(context.Background(), "d75be692-5f58-4534-b7fe-4d6e51c53a51")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Log(res)
}
