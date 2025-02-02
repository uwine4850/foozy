package securetest

import (
	"errors"
	"testing"

	"github.com/uwine4850/foozy/pkg/secure"
)

var (
	hashKey        = []byte("1234567890abcdef1234567890abcdef") // 32 bytes
	invalidHashKey = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdef") // 32 bytes
	blockKey       = []byte("abcdefghijklmnopqrstuvwx12345678") // 32 bytes
)

func TestGenerateAndVerifyHMAC(t *testing.T) {
	hmacData, err := secure.GenerateHMAC(hashKey, []byte("secret data"))
	if err != nil {
		t.Error(err)
	}
	isValid, err := secure.VerifyHMAC(hashKey, []byte("secret data"), hmacData)
	if err != nil {
		t.Error(err)
	}
	if !isValid {
		t.Errorf("The data has not been validated.")
	}
}

func TestHMACValidationFailed(t *testing.T) {
	hmacData, err := secure.GenerateHMAC(invalidHashKey, []byte("secret data"))
	if err != nil {
		t.Error(err)
	}
	isValid, err := secure.VerifyHMAC(hashKey, []byte("secret data"), hmacData)
	if err != nil {
		t.Error(err)
	}
	if isValid {
		t.Errorf("The data was validated, although it was not valid.")
	}
}

func TestEncryptAndDecryptData(t *testing.T) {
	enc, err := secure.Encrypt(blockKey, []byte("secret data"))
	if err != nil {
		t.Error(err)
	}
	dec, err := secure.Decrypt(blockKey, enc)
	if err != nil {
		t.Error(err)
	}
	if string(dec) != "secret data" {
		t.Errorf("The encrypt and decrypt data do not match.")
	}
}

func TestCreateReadSecureData(t *testing.T) {
	message := "secret data"
	secureData, err := secure.CreateSecureData(hashKey, blockKey, &message)
	if err != nil {
		t.Error(err)
	}
	var data string
	if err := secure.ReadSecureData(hashKey, blockKey, secureData, &data); err != nil {
		t.Error(err)
	}
	if data != "secret data" {
		t.Errorf("The create and read data do not match.")
	}
}

func TestErrInvalidHMAC(t *testing.T) {
	message := "secret data"
	secureData, err := secure.CreateSecureData(hashKey, blockKey, &message)
	if err != nil {
		t.Error(err)
	}
	var data string
	errOk := secure.ReadSecureData(invalidHashKey, blockKey, secureData, &data)
	errInvHMAC := secure.ErrInvalidHMAC{}
	if errOk == nil || !errors.Is(errOk, errInvHMAC) {
		t.Errorf("HMAC is not valid, but no error was raised.")
	}
}
