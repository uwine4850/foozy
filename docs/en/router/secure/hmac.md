## package secure
This section of the security package is responsible for operations with HMAC.

__GenerateHMAC__
```
GenerateHMAC(hashKey []byte, data []byte) ([]byte, error)
```
Generates data encrypted using HMAC. A hash key is used for generation 
32 bytes in size and SHA-256 algorithm.

__VerifyHMAC__
VerifyHMAC(hashKey []byte, data []byte, hmacCode []byte) (bool, error)
```
Performs HMAC verification. To do this, you need to transfer the hash key with which the data was encrypted, the data that is expected and the previously encrypted data.
```

__Encrypt__
```
Encrypt(blockKey []byte, data []byte) (string, error)
```
Encrypts data using AES algorithm and GCM mode.

__Decrypt__
```
Decrypt(blockKey []byte, enc string) ([]byte, error)
```
Decrypts data that was encrypted using the secure.Encrypt method.

__CreateSecureData__
```
CreateSecureData(hashKey []byte, blockKey []byte, writeData interface{}) (string, error)
```
Creates secure data using HMAC and conventional encryption.

__ReadSecureData__
```
ReadSecureData(hashKey []byte, blockKey []byte, secureData string, readData interface{}) error
```
Reads data that is encrypted using secure.CreateSecureData.