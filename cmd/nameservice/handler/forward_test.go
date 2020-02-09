package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/nameservice"
	"github.com/mailchain/mailchain/nameservice/nameservicetest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestForward(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("contents"))
		}),
	)
	defer server.Close()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		resolver nameservice.ForwardLookup
	}
	tests := []struct {
		name       string
		args       args
		req        *http.Request
		wantStatus int
		wantBody   string
	}{
		{
			"success",
			args{
				func() nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					m.EXPECT().ResolveName(gomock.Any(), "ethereum", "mainnet", "test.eth").Return(encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"), nil)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?domain-name=test.eth", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "ethereum",
					"network":  "mainnet",
				})
				return req
			}(),
			200,
			"{\"address\":\"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761\"}\n",
		},
		{
			"err-format",
			args{
				func() nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					m.EXPECT().
						ResolveName(gomock.Any(), "ethereum", "mainnet", "test.eth").
						Return(nil, nameservice.ErrFormat)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?domain-name=test.eth", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "ethereum",
					"network":  "mainnet",
				})
				return req
			}(),
			200,
			"{\"address\":\"\",\"status\":1}\n",
		},
		{
			"err-serv-fail",
			args{
				func() nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					m.EXPECT().
						ResolveName(gomock.Any(), "ethereum", "mainnet", "test.eth").
						Return(nil, nameservice.ErrServFail)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?domain-name=test.eth", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "ethereum",
					"network":  "mainnet",
				})
				return req
			}(),
			200,
			"{\"address\":\"\",\"status\":2}\n",
		},
		{
			"err-nx-domain",
			args{
				func() nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					m.EXPECT().
						ResolveName(gomock.Any(), "ethereum", "mainnet", "test.eth").
						Return(nil, nameservice.ErrNXDomain)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?domain-name=test.eth", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "ethereum",
					"network":  "mainnet",
				})
				return req
			}(),
			200,
			"{\"address\":\"\",\"status\":3}\n",
		},
		{
			"err-not-imp",
			args{
				func() nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					m.EXPECT().
						ResolveName(gomock.Any(), "ethereum", "mainnet", "test.eth").
						Return(nil, nameservice.ErrNotImp)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?domain-name=test.eth", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "ethereum",
					"network":  "mainnet",
				})
				return req
			}(),
			200,
			"{\"address\":\"\",\"status\":4}\n",
		},
		{
			"err-refused",
			args{
				func() nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					m.EXPECT().
						ResolveName(gomock.Any(), "ethereum", "mainnet", "test.eth").
						Return(nil, nameservice.ErrRefused)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?domain-name=test.eth", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "ethereum",
					"network":  "mainnet",
				})
				return req
			}(),
			200,
			"{\"address\":\"\",\"status\":5}\n",
		},
		{
			"err-invalid-protocol",
			args{
				func() nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					m.EXPECT().ResolveName(gomock.Any(), "invalid", "mainnet", "test.eth").Return(encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"), nil)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?domain-name=test.eth", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "invalid",
					"network":  "mainnet",
				})
				return req
			}(),
			500,
			"{\"code\":500,\"message\":\"failed to encode address: \\\"invalid\\\" unsupported protocol\"}\n",
		},
		{
			"err-unknown",
			args{
				func() nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					m.EXPECT().ResolveName(gomock.Any(), "ethereum", "mainnet", "test.eth").Return(nil, errors.Errorf("failed"))
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?domain-name=test.eth", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "ethereum",
					"network":  "mainnet",
				})
				return req
			}(),
			500,
			"{\"code\":500,\"message\":\"failed\"}\n",
		},
		{
			"err-missing-domain-name-query",
			args{
				func() nameservice.ForwardLookup {
					m := nameservicetest.NewMockForwardLookup(mockCtrl)
					return m
				}(),
			},
			func() *http.Request {
				req := httptest.NewRequest("GET", "/?", nil)
				req = mux.SetURLVars(req, map[string]string{
					"protocol": "ethereum",
					"network":  "mainnet",
				})
				return req
			}(),
			412,
			"{\"code\":412,\"message\":\"domain-name must be specified exactly once\"}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Forward(tt.args.resolver))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, tt.req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tt.wantStatus, rr.Code) {
				t.Errorf("Forward() returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}
			if !assert.Equal(t, tt.wantBody, rr.Body.String()) {
				t.Errorf("Forward() returned unexpected body: got %v want %v",
					rr.Body.String(), tt.wantBody)
			}
		})
	}
}
