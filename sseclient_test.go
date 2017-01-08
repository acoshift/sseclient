package sseclient_test

import (
	"log"
	"net/http"
	"testing"
	"time"

	. "github.com/acoshift/sseclient"
)

func ExampleNotify() {
	req, err := http.NewRequest(http.MethodGet, "https://dinosaur-facts.firebaseio.com/dinosaurs.json", nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	eventChan := make(chan *Event)
	cancel, err := Notify(client, req, eventChan)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		// cancel after 10 seconds passed
		time.Sleep(10 * time.Second)
		cancel()
	}()
	for event := range eventChan {
		log.Printf("%s: %s", event.Type, string(event.Data))
	}
}

func TestNotify(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://dinosaur-facts.firebaseio.com/dinosaurs.json", nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	eventChan := make(chan *Event)
	cancel, err := Notify(client, req, eventChan)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		time.Sleep(4 * time.Second)
		cancel()
		if _, ok := <-eventChan; ok {
			t.Errorf("channel not close after cancel")
		}
	}()
	for event := range eventChan {
		log.Printf("%s: %s", event.Type, string(event.Data))
	}
}

func TestNotify_empty(t *testing.T) {
	_, err := Notify(nil, nil, nil)
	if err != ErrClientNil {
		t.Errorf("expected error to be \"%v\", got \"%v\"", ErrClientNil, err)
	}

	_, err = Notify(&http.Client{}, nil, nil)
	if err != ErrReqNil {
		t.Errorf("expected error to be \"%v\", got \"%v\"", ErrReqNil, err)
	}

	req := &http.Request{}
	_, err = Notify(&http.Client{}, req, nil)
	if err != ErrChanNil {
		t.Errorf("expected error to be \"%v\", got \"%v\"", ErrChanNil, err)
	}
}
