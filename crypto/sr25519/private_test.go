package sr25519

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivateKeyFromBytes(t *testing.T) {
	type args struct {
		pk []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *PrivateKey
		wantErr bool
	}{
		{
			"success-sofia-seed",
			args{
				sofiaSeed,
			},
			&sofiaPrivateKey,
			false,
		},
		{
			"success-sofia-bytes",
			args{
				sofiaPrivateKeyBytes,
			},
			&sofiaPrivateKey,
			false,
		},
		{
			"success-charlotte-seed",
			args{
				charlotteSeed,
			},
			&charlottePrivateKey,
			false,
		},
		{
			"success-charlotte-bytes",
			args{
				charlottePrivateKeyBytes,
			},
			&charlottePrivateKey,
			false,
		},
		{
			"err-len",
			args{
				[]byte{57, 212, 201},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrivateKeyFromBytes(tt.args.pk)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrivateKeyFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey_Kind(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"success",
			"sr25519",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := PrivateKey{}
			if got := pk.Kind(); got != tt.want {
				t.Errorf("PrivateKey.Kind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey_Bytes(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		pk   PrivateKey
		want []byte
	}{
		{
			"sucess-sofia",
			sofiaPrivateKey,
			sofiaPrivateKeyBytes,
		},
		{
			"sucess-charllotte",
			charlottePrivateKey,
			charlottePrivateKeyBytes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pk.Bytes(); !assert.Equal(tt.want, got) {
				t.Errorf("PrivateKey.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey_PublicKey(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		pk   PrivateKey
		want PublicKey
	}{
		{
			"sofia-sr25519",
			sofiaPrivateKey,
			sofiaPublicKey,
		},
		{
			"charlotte-sr25519",
			charlottePrivateKey,
			charlottePublicKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pk.PublicKey(); !assert.Equal(tt.want, got) {
				t.Errorf("PrivateKey.PublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey_Sign(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		pk      PrivateKey
		msg     []byte
		want    []byte
		wantErr bool
	}{
		{
			"success-charlotte",
			charlottePrivateKey,
			[]byte("message"),
			[]byte{0x66, 0x4f, 0x87, 0xa8, 0xb9, 0x0, 0x84, 0xfd, 0xb6, 0x40, 0x7a, 0xe2, 0x1a, 0xd0, 0x7b, 0x94, 0x8d, 0xb8, 0x27, 0x62, 0xaa, 0xce, 0xfa, 0x35, 0x8b, 0x4c, 0x59, 0x6d, 0x61, 0xbe, 0x25, 0x3a, 0x3e, 0xac, 0x6, 0xb9, 0xa7, 0x75, 0x9e, 0xa3, 0xd2, 0x45, 0xf9, 0x4a, 0x9b, 0x61, 0x0, 0xad, 0x9e, 0x68, 0x77, 0x87, 0x99, 0xce, 0xa, 0x54, 0x9f, 0x44, 0x32, 0xe, 0x4, 0xef, 0x76, 0x87},
			false,
		},
		{
			"success-sofia",
			sofiaPrivateKey,
			[]byte("egassem"),
			[]byte{0x6c, 0x9a, 0x86, 0x4d, 0xbc, 0x4b, 0xd3, 0xf4, 0xf8, 0x31, 0xfd, 0x8e, 0x84, 0xf7, 0x83, 0x8f, 0x71, 0xf1, 0x1d, 0xd1, 0xc, 0xa, 0xbb, 0x14, 0xf7, 0xb5, 0xe3, 0x43, 0x6b, 0x35, 0x6, 0x5a, 0x3a, 0x5d, 0xba, 0x3a, 0x28, 0x9f, 0xfe, 0xd4, 0x34, 0x1c, 0xab, 0x7f, 0x18, 0xdc, 0x51, 0x45, 0xdb, 0x68, 0x5e, 0xf3, 0x67, 0xb6, 0x54, 0xc4, 0xe0, 0x1c, 0x8b, 0x2d, 0x22, 0xab, 0x1c, 0x85},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pk.Sign(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(len(tt.want), len(got)) {
				t.Errorf("Sign() = %v,\n want %v", got, tt.want)
			}
		})
	}
}
