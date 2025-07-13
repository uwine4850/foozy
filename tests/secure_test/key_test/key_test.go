package key_test

import (
	"testing"

	"github.com/uwine4850/foozy/pkg/secure"
)

func TestGenerateBytesKeys(t *testing.T) {
	key := secure.NewKey()
	key.GenerateBytesKeys(10)
	if len(key.StaticKey()) != 10 {
		t.Error("StaticKey is not valid")
	}
	if len(key.BlockKey()) != 10 {
		t.Error("BlockKey is not valid")
	}
	if len(key.HashKey()) != 10 {
		t.Error("HashKey is not valid")
	}
}

func TestGenerate32BytesKeys(t *testing.T) {
	key := secure.NewKey()
	key.Generate32BytesKeys()
	if len(key.StaticKey()) != 32 {
		t.Error("StaticKey is not valid")
	}
	if len(key.BlockKey()) != 32 {
		t.Error("BlockKey is not valid")
	}
	if len(key.HashKey()) != 32 {
		t.Error("HashKey is not valid")
	}
}

func TestOldKeys(t *testing.T) {
	key := secure.NewKey()
	key.Generate32BytesKeys()
	blockKey := key.BlockKey()
	hashKey := key.HashKey()
	key.Generate32BytesKeys()
	if blockKey != key.OldBlockKey() {
		t.Error("old BlockKey does not match with the expectation")
	}
	if hashKey != key.OldHashKey() {
		t.Error("old HashKey does not match with the expectation")
	}
}
