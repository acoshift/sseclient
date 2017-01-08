// Package sseclient is the client of server side events
package sseclient

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
)

// Event is the event sent from server
type Event struct {
	Type string
	Data []byte
}

// CancelFunc is the function use for cancel notify
type CancelFunc func()

// Errors
var (
	ErrClientNil = errors.New("sse: client is nil")
	ErrReqNil    = errors.New("sse: req is nil")
	ErrChanNil   = errors.New("sse: channel is nil")
)

var emptyFunc = func() {}

const (
	headerEventStream = "text/event-stream"
)

var (
	eventEvent = []byte("event")
	eventData  = []byte("data")
	delim      = []byte{':', ' '}
)

// Notify does a request and notify to given channel when receive event
// CancelFunc use to close response body and channel
// Do not manually close channel
func Notify(client *http.Client, req *http.Request, ch chan *Event) (CancelFunc, error) {
	if client == nil {
		return emptyFunc, ErrClientNil
	}
	if req == nil {
		return emptyFunc, ErrReqNil
	}
	if ch == nil {
		return emptyFunc, ErrChanNil
	}
	if req.Header.Get("Accept") != headerEventStream {
		req.Header.Add("Accept", headerEventStream)
	}
	resp, err := client.Do(req)
	if err != nil {
		return emptyFunc, err
	}
	cancel := func() {
		resp.Body.Close()
		close(ch)
	}
	go func() {
		reader := bufio.NewReader(resp.Body)
		var event *Event
		for {
			bs, err := reader.ReadBytes('\n')
			if err != nil {
				return
			}
			if len(bs) < 2 {
				continue
			}
			rec := bytes.SplitN(bs, delim, 2)
			if len(rec) < 2 {
				continue
			}
			if bytes.Equal(rec[0], eventEvent) {
				event = &Event{}
				event.Type = string(bytes.TrimSpace(rec[1]))
			} else if bytes.Equal(rec[0], eventData) {
				if event == nil {
					continue
				}
				event.Data = bytes.TrimSpace(rec[1])
				ch <- event
				event = nil
			}
		}
	}()
	return cancel, nil
}
