package models

import (
	"reflect"
	"time"
)

type User struct {
	Id    string
	Email string
	Name  string
}

type ReadabilityStatus string

const (
	ReadabilityUnknown    ReadabilityStatus = "UNKNOWN"
	ReadabilityFailed     ReadabilityStatus = "FAILED"
	ReadabilityProcessing ReadabilityStatus = "PROCESSING"
	ReadabilityComplete   ReadabilityStatus = "COMPLETE"
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

type ReadabilityResponseMetadata struct {
	RequestId      string
	HTTPStatusCode int
	RetryAttempts  int
}

type ReadabilityResponse struct {
	ResponseMetadata ReadabilityResponseMetadata
	SynthesisTask    SynthesisTask
}

type Page struct {
	Created             time.Time         `json:"created" mapkey:"created"`
	Id                  string            `json:"id" mapkey:"id"`
	Url                 string            `json:"url" mapkey:"url"`
	LastCrawled         time.Time         `json:"last_crawled" mapkey:"last_crawled"`
	Title               string            `json:"title,omitempty" mapkey:"title"`
	Description         string            `json:"description,omitempty" mapkey:"description"`
	ImageUrl            string            `json:"image_url,omitempty" mapkey:"image_url"`
	IsReadable          bool              `json:"is_readable" mapkey:"is_readable"`
	ReadabilityStatus   ReadabilityStatus `json:"readability_status" mapkey:"readability_status"`
	ReadabilityTaskData ReadabilityResponse
}

func (p Page) ToMap() map[string]any {
	out := make(map[string]any)

	val := reflect.ValueOf(p)
	t := val.Type()
	n := val.NumField()
	for i := 0; i < n; i++ {
		v_field := val.Field(i)
		if v_field.IsZero() {
			continue
		}
		t_field := t.Field(i)
		field_val := v_field.Interface()
		tag := t_field.Tag.Get("mapkey")
		out[tag] = field_val
	}

	return out
}
