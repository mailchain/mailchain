// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mailchain/mailchain/crypto/cipher"
	"github.com/mailchain/mailchain/crypto/cipher/ciphertest"
	"github.com/mailchain/mailchain/encoding/encodingtest"
	"github.com/mailchain/mailchain/internal/clients/etherscan"
	"github.com/mailchain/mailchain/internal/keystore"
	"github.com/mailchain/mailchain/internal/keystore/kdf/multi"
	"github.com/mailchain/mailchain/internal/keystore/keystoretest"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/mailbox"
	"github.com/mailchain/mailchain/internal/mailbox/mailboxtest"
	"github.com/mailchain/mailchain/stores"
	"github.com/mailchain/mailchain/stores/statemock"
	"github.com/mailchain/mailchain/stores/storestest"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_FetchMessages(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		inbox            stores.State
		cache            stores.Cache
		receivers        map[string]mailbox.Receiver
		ks               keystore.Store
		deriveKeyOptions multi.OptionsBuilders
	}
	tests := []struct {
		name       string
		args       args
		req        *http.Request
		wantStatus int
	}{
		{
			"422-empty-address",
			args{},
			httptest.NewRequest("GET", "/?address=&network=mainnet&protocol=ethereum", nil),
			http.StatusUnprocessableEntity,
		},
		{
			"422-receiver-not-supported",
			args{},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusUnprocessableEntity,
		},
		{
			"422-receiver-no-configured",
			args{
				receivers: map[string]mailbox.Receiver{
					"ethereum/mainnet": nil,
				},
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusUnprocessableEntity,
		},
		{
			"406-no-private-key-found",
			args{
				receivers: map[string]mailbox.Receiver{
					"ethereum/mainnet": etherscan.APIClient{},
				},
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet").Return(false).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil).WithContext(context.Background()),
			http.StatusNotAcceptable,
		},
		{
			"406-receiver-network-error",
			args{
				receivers: func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							receiver := mailboxtest.NewMockReceiver(mockCtrl)
							receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}).
								Return(nil, errors.New("network not supported")).Times(1)
							return receiver
						}(),
					}
				}(),
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet").Return(true).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusNotAcceptable,
		},
		{
			"500-receiver-internal-error",
			args{
				receivers: func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							receiver := mailboxtest.NewMockReceiver(mockCtrl)
							receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}).
								Return(nil, errors.New("internal error")).Times(1)
							return receiver
						}(),
					}
				}(),
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet").Return(true).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusInternalServerError,
		},
		{
			"500-decrypter-error",
			args{
				receivers: func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							receiver := mailboxtest.NewMockReceiver(mockCtrl)
							receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}).
								Return([]mailbox.Transaction{
									{
										Data:    encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a221620d3c47ef741473ebf42773d25687b7540a3d96429aec07dd1ce66c0d4fd16ea13"),
										BlockID: []byte("YS1ibG9jay1udW1iZXI="),
										Hash:    []byte("YS1oYXNo"),
									},
								}, nil).Times(1)
							return receiver
						}(),
					}
				}(),
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet").Return(true).Times(1)
					store.EXPECT().GetDecrypter([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet", byte(0x73), multi.OptionsBuilders{}).Return(nil, errors.New("not found"))
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusInternalServerError,
		},
		{
			"200-message",
			args{
				cache: func() stores.Cache {
					m := storestest.NewMockCache(mockCtrl)
					m.EXPECT().GetMessage("test://TestReadMessage/success-2204f3d89e5a").Return([]byte{}, errors.Errorf("message not found"))
					m.EXPECT().SetMessage("test://TestReadMessage/success-2204f3d89e5a", gomock.Any()).Return(nil)
					return m
				}(),
				inbox: func() stores.State {
					inbox := statemock.NewMockState(mockCtrl)
					inbox.EXPECT().GetReadStatus(mail.ID{71, 236, 160, 17, 227, 43, 82, 199, 16, 5, 173, 138, 143, 117, 225, 180, 76, 146, 201, 159, 209, 46, 67, 188, 207, 229, 113, 227, 194, 209, 61, 46, 154, 130, 106, 85, 15, 95, 246, 59, 36, 122, 244, 113}).Return(false, nil).Times(1)
					inbox.EXPECT().PutMessage("ethereum", "mainnet", "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761", gomock.Any()).Return(nil).Times(1)
					return inbox
				}(),
				receivers: func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							receiver := mailboxtest.NewMockReceiver(mockCtrl)
							receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}).
								Return([]mailbox.Transaction{
									{
										Data:    encodingtest.MustDecodeHex("500801120f7365637265742d6c6f636174696f6e1a221620d3c47ef741473ebf42773d25687b7540a3d96429aec07dd1ce66c0d4fd16ea13"),
										BlockID: []byte("YS1ibG9jay1udW1iZXI="),
										Hash:    []byte("YS1oYXNo"),
									},
								}, nil).Times(1)
							return receiver
						}(),
					}
				}(),
				ks: func() keystore.Store {
					decrypted, _ := ioutil.ReadFile("./testdata/simple.golden.eml")
					decrypter := ciphertest.NewMockDecrypter(mockCtrl)
					gomock.InOrder(
						decrypter.EXPECT().Decrypt(cipher.EncryptedContent(encodingtest.MustDecodeHex("7365637265742d6c6f636174696f6e"))).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil),
						decrypter.EXPECT().Decrypt(cipher.EncryptedContent(encodingtest.MustDecodeHex("7365637265742d6c6f636174696f6e"))).Return([]byte("test://TestReadMessage/success-2204f3d89e5a"), nil),
						decrypter.EXPECT().Decrypt(cipher.EncryptedContent([]byte{0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x61, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65})).Return(decrypted, nil),
					)

					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet").Return(true).Times(1)
					store.EXPECT().GetDecrypter([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet", byte(0x73), multi.OptionsBuilders{}).Return(decrypter, nil)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusOK,
		},
		{
			"200-zero-messages",
			args{
				inbox: func() stores.State {
					return statemock.NewMockState(mockCtrl)
				}(),
				receivers: func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							receiver := mailboxtest.NewMockReceiver(mockCtrl)
							receiver.EXPECT().Receive(context.Background(), "ethereum", "mainnet", []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}).
								Return([]mailbox.Transaction{}, nil).Times(1)
							return receiver
						}(),
					}
				}(),
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet").Return(true).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusOK,
		},
	}

	testName := t.Name()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(FetchMessages(tt.args.inbox, tt.args.cache, tt.args.receivers, tt.args.ks, tt.args.deriveKeyOptions))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, tt.req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tt.wantStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}

			goldenResponse, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/response-%s.json", testName, tt.name))
			if err != nil {
				assert.FailNow(t, err.Error())
			}

			if !assert.JSONEq(t, string(goldenResponse), rr.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), goldenResponse)
			}

		})
	}
}

func Test_GetMessages(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type args struct {
		inbox     stores.State
		receivers map[string]mailbox.Receiver
		ks        keystore.Store
	}
	tests := []struct {
		name       string
		args       args
		req        *http.Request
		wantStatus int
	}{
		{
			"200-message",
			args{
				inbox: func() stores.State {
					stateMock := statemock.NewMockState(mockCtrl)
					timeValue, err := time.Parse(time.RFC3339, "2019-03-12T20:23:13Z")
					if err != nil {
						t.Fatal(err)
					}
					stateMock.EXPECT().GetMessages("ethereum", "mainnet", "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761").Times(1).Return([]stores.Message{
						{
							Headers: stores.Header{
								Date:        timeValue,
								From:        "\u003c5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum\u003e",
								To:          "\u003c4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum\u003e",
								MessageID:   "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471",
								ContentType: "text/plain; charset=\"UTF-8\"",
							},
							Body:                    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac. Sed vestibulum dui lectus, quis eleifend urna mollis eu. Integer dictum metus ut sem rutrum aliquet. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Phasellus eget euismod nibh. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer bibendum venenatis sem sed auctor. Ut aliquam eu diam nec fermentum. Sed turpis nulla, viverra ac efficitur ac, fermentum vel sapien. Curabitur vehicula risus id odio congue tempor. Mauris tincidunt feugiat risus, eget auctor magna blandit sit amet. Curabitur consectetur, dolor eu imperdiet varius, dui neque mattis neque, vel fringilla magna tortor ut risus. Cras cursus sem et nisl interdum molestie. Aliquam auctor sodales blandit.\r\n",
							Subject:                 "Hello world",
							Status:                  "ok",
							StatusCode:              "",
							Read:                    false,
							BlockID:                 "YS1ibG9jay1udW1iZXI=",
							BlockIDEncoding:         "hex/0x-prefix",
							TransactionHash:         "YS1oYXNo",
							TransactionHashEncoding: "hex/0x-prefix",
						},
					}, nil)
					return stateMock
				}(),
				receivers: func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							return mailboxtest.NewMockReceiver(mockCtrl)
						}(),
					}
				}(),
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet").Return(true).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusOK,
		},
		{
			"200-zero-messages",
			args{
				inbox: func() stores.State {
					stateMock := statemock.NewMockState(mockCtrl)
					stateMock.EXPECT().GetMessages("ethereum", "mainnet", "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761").Times(1)
					return stateMock
				}(),
				receivers: func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							return mailboxtest.NewMockReceiver(mockCtrl)
						}(),
					}
				}(),
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet").Return(true).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusOK,
		},
		{
			"406-no-private-key-found",
			args{
				receivers: map[string]mailbox.Receiver{
					"ethereum/mainnet": etherscan.APIClient{},
				},
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet").Return(false).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil).WithContext(context.Background()),
			http.StatusNotAcceptable,
		},
		{
			"422-empty-address",
			args{},
			httptest.NewRequest("GET", "/?address=&network=mainnet&protocol=ethereum", nil),
			http.StatusUnprocessableEntity,
		},
		{
			"422-receiver-not-supported",
			args{},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusUnprocessableEntity,
		},
		{
			"422-receiver-no-configured",
			args{
				receivers: map[string]mailbox.Receiver{
					"ethereum/mainnet": nil,
				},
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusUnprocessableEntity,
		},
		{
			"500-receiver-internal-error",
			args{
				inbox: func() stores.State {
					stateMock := statemock.NewMockState(mockCtrl)
					stateMock.EXPECT().GetMessages("ethereum", "mainnet", "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761").Times(1).Return(nil, errors.New("internal error"))
					return stateMock
				}(),
				receivers: func() map[string]mailbox.Receiver {
					return map[string]mailbox.Receiver{
						"ethereum/mainnet": func() mailbox.Receiver {
							return mailboxtest.NewMockReceiver(mockCtrl)
						}(),
					}
				}(),
				ks: func() keystore.Store {
					store := keystoretest.NewMockStore(mockCtrl)
					store.EXPECT().HasAddress([]byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61}, "ethereum", "mainnet").Return(true).Times(1)
					return store
				}(),
			},
			httptest.NewRequest("GET", "/?address=0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761&network=mainnet&protocol=ethereum", nil),
			http.StatusInternalServerError,
		},
	}

	testName := t.Name()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetMessages(tt.args.inbox, tt.args.receivers, tt.args.ks))

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, tt.req)

			// Check the status code is what we expect.
			if !assert.Equal(t, tt.wantStatus, rr.Code) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tt.wantStatus)
			}

			goldenResponse, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/response-%s.json", testName, tt.name))
			if err != nil {
				assert.FailNow(t, err.Error())
			}

			if !assert.JSONEq(t, string(goldenResponse), rr.Body.String()) {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), goldenResponse)
			}
		})
	}
}

func Test_parseGetMessagesRequest(t *testing.T) {
	type args struct {
		queryParams map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    *GetMessagesRequest
		wantErr bool
	}{
		{
			"success",
			args{
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "mainnet",
					"protocol": "ethereum",
				},
			},
			&GetMessagesRequest{
				Address:      "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				Network:      "mainnet",
				Protocol:     "ethereum",
				addressBytes: []byte{0x56, 0x2, 0xea, 0x95, 0x54, 0xb, 0xee, 0x46, 0xd0, 0x3b, 0xa3, 0x35, 0xee, 0xd6, 0xf4, 0x9d, 0x11, 0x7e, 0xab, 0x95, 0xc8, 0xab, 0x8b, 0x71, 0xba, 0xe2, 0xcd, 0xd1, 0xe5, 0x64, 0xa7, 0x61},
			},
			false,
		},
		{
			"err-invalid-protocol",
			args{
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "mainnet",
					"protocol": "invalid",
				},
			},
			nil,
			true,
		},
		{
			"err-empty-address",
			args{
				map[string]string{
					"address":  "",
					"network":  "mainnet",
					"protocol": "ethereum",
				},
			},
			nil,
			true,
		},
		{
			"err-empty-protocol",
			args{
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "mainnet",
					"protocol": "",
				},
			},
			nil,
			true,
		},
		{
			"err-empty-network",
			args{
				map[string]string{
					"address":  "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
					"network":  "",
					"protocol": "ethereum",
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			q := req.URL.Query()
			for k, v := range tt.args.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			got, err := parseGetMessagesRequest(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGetMessagesRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("parseGetMessagesRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
