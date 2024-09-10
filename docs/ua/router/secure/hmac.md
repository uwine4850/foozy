## package secure
Даний розділ безпекового пакету відповідає за операції із HMAC.

__GenerateHMAC__
```
GenerateHMAC(hashKey []byte, data []byte) ([]byte, error)
```
Генерує дані зашифровані за допомогою HMAC. Для генерації використовується хеш-ключ 
розміром 32 байта та алгоритм SHA-256.

__VerifyHMAC__
VerifyHMAC(hashKey []byte, data []byte, hmacCode []byte) (bool, error)
```
Проводить верифікаю HMAC. Для цього потрібно передати хеш-ключ за допомогою якого шифрувалися дані, дані які очікуються та зашифровані раніше дані.
```

__Encrypt__
```
Encrypt(blockKey []byte, data []byte) (string, error)
```
Зашифровує дані за допомогою AES алгоритму та GCM моді.

__Decrypt__
```
Decrypt(blockKey []byte, enc string) ([]byte, error)
```
Розшифровує дані, які були зашифрованні за допомогою метода secure.Encrypt.

__CreateSecureData__
```
CreateSecureData(hashKey []byte, blockKey []byte, writeData interface{}) (string, error)
```
Створює безпечні дані за допомогою HMAC та звичайного шифрувааня.

__ReadSecureData__
```
ReadSecureData(hashKey []byte, blockKey []byte, secureData string, readData interface{}) error
```
Читає дані, які зашифрованні з допомогою secure.CreateSecureData.