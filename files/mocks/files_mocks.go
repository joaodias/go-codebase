package filesmocks

import (
	"errors"
	files "github.com/joaodias/go-codebase/files"
	"strings"
)

type FakeWebClient struct{}

func GetTestClient() *files.Client {
	client := new(files.Client)
	client.HTTP = new(FakeWebClient)
	return client
}

func (r *FakeWebClient) Get(url string) ([]byte, error) {
	if -1 != strings.Index(url, "somethingBad") {
		return nil, errors.New("Cannot get the desired web content.")
	}
	return []byte(jsonSample), nil
}

var jsonSample = ` { "fake" : "data" } `
