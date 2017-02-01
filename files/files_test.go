package files_test

import (
	"reflect"
	"testing"

	"github.com/joaodias/go-codebase/files"
	filesmocks "github.com/joaodias/go-codebase/files/mocks"
)

func TestHTTPDownload(t *testing.T) {
	type args struct {
		webClient *files.Client
		url       string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"Successfully get web content.", args{filesmocks.GetTestClient(), "somethingGood"}, []byte(` { "fake" : "data" } `), false},
		{"Cannot get web content.", args{filesmocks.GetTestClient(), "somethingBad"}, nil, true},
	}
	for _, tt := range tests {
		got, err := files.HTTPDownload(*tt.args.webClient, tt.args.url)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. HTTPDownload() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. HTTPDownload() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
