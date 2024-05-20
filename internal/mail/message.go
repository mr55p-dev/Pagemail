package mail

import "io"

type Message struct {
	From    string
	To      string
	Subject string
	Body    io.Reader
	Tags    []struct {
		Name  string
		Value string
	}
}

type MessageOpt func(msg *Message)

func WithSender(from string) MessageOpt {
	return func(msg *Message) {
		msg.From = from
	}
}

func WithSubject(subject string) MessageOpt {
	return func(msg *Message) {
		msg.Subject = subject
	}
}

func WithBody(body io.Reader) MessageOpt {
	return func(msg *Message) {
		msg.Body = body
	}
}

func WithTags(tags map[string]string) MessageOpt {
	return func(msg *Message) {
		for k, v := range tags {
			msg.Tags = append(msg.Tags, struct {
				Name  string
				Value string
			}{Name: k, Value: v})
		}
	}
}

const (
	DEFAULT_SUBJECT = "No subject"
	DEFAULT_FROM    = "mail@pagemail.io"
)

func MakeMessage(to string, opts ...MessageOpt) Message {
	msg := Message{
		To:      to,
		From:    DEFAULT_FROM,
		Subject: DEFAULT_SUBJECT,
	}
	for _, opt := range opts {
		opt(&msg)
	}
	return msg
}
