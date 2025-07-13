package hmac_test

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/uwine4850/foozy/pkg/secure"
)

var (
	hashKey        = []byte("1234567890abcdef1234567890abcdef") // 32 bytes
	invalidHashKey = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdef") // 32 bytes
	blockKey       = []byte("abcdefghijklmnopqrstuvwx12345678") // 32 bytes
)

func TestGenerateHmac(t *testing.T) {
	hmac, err := secure.GenerateHMAC(hashKey, []byte("OK"))
	if err != nil {
		t.Error(err)
	}
	if hex.EncodeToString(hmac) != "b4e4ddcff94c0cf7f4e1043dbb4e5ebb92d1d8592913eb756f123d15e8abdcb7" {
		t.Error("the generated hmac does not match the expected hmac")
	}
}

func TestVerifyHMAC(t *testing.T) {
	hmac, err := secure.GenerateHMAC(hashKey, []byte("OK"))
	if err != nil {
		t.Error(err)
	}
	ok, err := secure.VerifyHMAC(hashKey, []byte("OK"), hmac)
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Error("hmac failed verification")
	}
}

func TestVerifyHMACError(t *testing.T) {
	hmac, err := secure.GenerateHMAC(hashKey, []byte("OK"))
	if err != nil {
		t.Error(err)
	}
	ok, err := secure.VerifyHMAC(invalidHashKey, []byte("OK"), hmac)
	if err != nil {
		t.Error(err)
	}
	if ok {
		t.Error("hmac doesn't have to be verified")
	}
}

func TestEncrypt(t *testing.T) {
	_, err := secure.Encrypt(blockKey, []byte("OK"))
	if err != nil {
		t.Error(err)
	}
}

func TestDecrypt(t *testing.T) {
	enc, err := secure.Encrypt(blockKey, []byte("OK"))
	if err != nil {
		t.Error(err)
	}
	dec, err := secure.Decrypt(blockKey, enc)
	if err != nil {
		t.Error(err)
	}
	if string(dec) != "OK" {
		t.Error("incorrect data decryption")
	}
}

func TestCreateSecureData(t *testing.T) {
	d := "OK"
	_, err := secure.CreateSecureData(hashKey, blockKey, &d)
	if err != nil {
		t.Error(err)
	}
}

func TestReadSecureData(t *testing.T) {
	d := "OK"
	secureData, err := secure.CreateSecureData(hashKey, blockKey, &d)
	if err != nil {
		t.Error(err)
	}
	var read string
	if err := secure.ReadSecureData(hashKey, blockKey, secureData, &read); err != nil {
		t.Error(err)
	}
	if read != "OK" {
		t.Error("the data read does not match the data expected")
	}
}

func TestReadSecureDataError(t *testing.T) {
	d := "OK"
	secureData, err := secure.CreateSecureData(hashKey, blockKey, &d)
	if err != nil {
		t.Error(err)
	}
	var read string
	err = secure.ReadSecureData(invalidHashKey, blockKey, secureData, &read)
	if err == nil {
		t.Error("error not found")
	} else {
		if !errors.Is(err, secure.ErrInvalidHMAC{}) {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}
