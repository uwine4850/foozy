## HMAC
Implementation and working with the HMAC algorithm.

#### GenerateHMAC
Generates an HMAC signature using the SHA-256 algorithm and a key.
```golang
func GenerateHMAC(hashKey []byte, data []byte) ([]byte, error) {
	newHMAC := hmac.New(sha256.New, hashKey)
	_, err := newHMAC.Write(data)
	if err != nil {
		return nil, err
	}
	return newHMAC.Sum(nil), nil
}
```

#### VerifyHMAC
Compares the received signature with the correct one.<br>
data - expected content.
```golang
func VerifyHMAC(hashKey []byte, data []byte, hmacCode []byte) (bool, error) {
	// A valid signature has been generated and is expected.
	expectedHMAC, err := GenerateHMAC(hashKey, data)
	if err != nil {
		return false, err
	}
	return hmac.Equal(expectedHMAC, hmacCode), nil
}
```

#### Encrypt
The function is designed to encrypt data using the AES encryption algorithm in GCM mode (Galois/Counter Mode).
```golang
func Encrypt(blockKey []byte, data []byte) (string, error) {
	cipherBlock, err := aes.NewCipher(blockKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}
```

#### Decrypt
Function is designed to decrypt data that was encrypted using the Encrypt function.
It creates AES and GCM using the same key and uses them to decrypt the data.
```golang
func Decrypt(blockKey []byte, enc string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return nil, err
	}

	cipherBlock, err := aes.NewCipher(blockKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}
```

#### CreateSecureData
Creates encrypted data using hmac and the `Encrypt` function.
```golang
func CreateSecureData(hashKey []byte, blockKey []byte, writeData interface{}) (string, error) {
	if reflect.TypeOf(writeData).Kind() != reflect.Pointer {
		return "", errors.New("the writeData argument must be a pointer")
	}
	data, err := json.Marshal(writeData)
	if err != nil {
		return "", err
	}

	genereatedHMAC, err := GenerateHMAC(hashKey, data)
	if err != nil {
		return "", err
	}
	data = append(data, genereatedHMAC...)

	return Encrypt(blockKey, data)
}
```

#### ReadSecureData
Reads encrypted data created using `CreateSecureData`. All keys must match.
```golang
func ReadSecureData(hashKey []byte, blockKey []byte, secureData string, readData interface{}) error {
	if reflect.TypeOf(readData).Kind() != reflect.Pointer {
		return errors.New("the writeData argument must be a pointer")
	}
	data, err := Decrypt(blockKey, secureData)
	if err != nil {
		return err
	}

	sig := data[len(data)-sha256.Size:]
	data1 := data[:len(data)-sha256.Size]

	isValidHMAC, err := VerifyHMAC(hashKey, data1, sig)
	if err != nil {
		return err
	}
	if !isValidHMAC {
		return ErrInvalidHMAC{}
	}

	if err := json.Unmarshal(data1, readData); err != nil {
		return err
	}

	return nil
}
```