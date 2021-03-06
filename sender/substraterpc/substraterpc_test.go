package substraterpc

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/types"
)

func TestNew(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("{\"result\":\"%v\"}", types.ExamplaryMetadataV11SubstrateString)))
	}))
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			"success",
			args{
				server.URL,
			},
			false,
			false,
		},
		{
			"failed",
			args{
				"host:23425",
			},
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got.client == nil) != tt.wantNil {
				t.Errorf("New() got = %v, wantNil %v", err, tt.wantErr)
				return
			}
		})
	}
}
