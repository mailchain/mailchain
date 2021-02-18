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

package addressing

import (
	"reflect"
	"testing"

	"github.com/mailchain/mailchain/encoding/encodingtest"
)

func TestDecodeByProtocol(t *testing.T) {
	type args struct {
		in       string
		protocol string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"ethereum",
			args{
				"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				"ethereum",
			},
			encodingtest.MustDecodeHex("5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"),
			false,
		},
		{
			"substrate",
			args{
				"5DJJhV3tVzsWG1jZfL157azn8iRyDC7HyNG1yh8v2nQYd994",
				"substrate",
			},
			encodingtest.MustDecodeBase32("5DJJhV3tVzsWG1jZfL157azn8iRyDC7HyNG1yh8v2nQYd994"),
			false,
		},
		{
			"algorand",
			args{
				"C7Z4NNMIMOGZW56JCILF6DVY4MBZJMHXUQ67W2WKVE6U5QJSIDPYUEAXQU",
				"algorand",
			},
			encodingtest.MustDecodeBase32("C7Z4NNMIMOGZW56JCILF6DVY4MBZJMHXUQ67W2WKVE6U5QJSIDPYUEAXQU"),
			false,
		},
		{
			"err",
			args{
				"0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761",
				"invalid",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeByProtocol(tt.args.in, tt.args.protocol)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeByProtocol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeByProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}
