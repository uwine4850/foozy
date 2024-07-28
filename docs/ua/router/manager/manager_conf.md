## IManagerConfig
Зберігає та надає доступ до деяких налаштувань.

__DebugConfig()__
```
DebugConfig() IManagerDebugConfig
```
Надає доступ до конфігурації дабагу.

__PrintLog__
```
PrintLog(enable bool)
```
Вмикає консольний лог. Тобто, деякі дані будуть виводитись у консоль.

__IsPrintLog__
```
IsPrintLog() bool
```
Повертає значення PrintLog.

__Key__
```
Key() IManagerDebugConfig
```
Надає доступ до конфігурації об'єкту ключа.

## IManagerDebugConfig

__Debug__
```
Debug(enable bool)
```
Вмикає або вимикає дебаг мод. Якщо ця опція уввімкнена, веб-сторінка буде 
відображати деякі помилки.

__IsDebug__
```
IsDebug() bool
```
Повертає значення Debug.

__ErrorLogging__
```
ErrorLogging(enable bool)
```
Вмикає логування вибраних помилок у .log файл. Шлях до файлу потрібно встановити 
у методі __ErrorLoggingFile__.

__IsErrorLogging__
```
IsErrorLogging() bool
```
Повертає значення ErrorLogging.

__ErrorLoggingFile__
```
ErrorLoggingFile(path string)
```
Встановлює шлях до файлу логів. Якщо файл не існує він буде створений.

__GetErrorLoggingFile__
```
ErrorLoggingFile(path string)
```
Повертає шлях до файлу логів.

__SkipLoggingLevel__
```
SkipLoggingLevel(skip int)
```
Встановлення рівня логування. За замовчуванням дорівнює 3.
Даний параметр показує місце виклику функції логування у проекті. 
Наприклад, номер 0 буде показувати місце викливу функції логування у реалізації фреймворку.

__LoggingLevel__
```
LoggingLevel() int
```
Повертає значення рівня логування.

## IKey

__HashKey__
```
HashKey() string
```
Генерує hash key.

__OldHashKey__
```
OldHashKey() string
```
Повертає старий ключ логування.
Цей метод використовується після того як hash key згенерований повторно.
Важливо зазначити, що даний метод повертає минулий ключ, починаючи із активного.

__BlockKey__
```
BlockKey() string
```
Генерує block key.

__OldBlockKey__
```
OldBlockKey() string
```
Повертає старий ключ логування.
Цей метод використовується після того як block key згенерований повторно.
Важливо зазначити, що даний метод повертає минулий ключ, починаючи із активного.

__StaticKey()__
```
StaticKey() string
```
Створює static key. Цей ключ генерується один раз під час запуску серверу, та не змінюється.

__Date() time.Time__
```
Date() time.Time
```
Повертає час останнього генерування ключів.

__GenerateBytesKeys__
```
GenerateBytesKeys(length int)
```
Генерує block key та hash key

__Generate32BytesKeys__
```
Generate32BytesKeys()
```
Генерує ключі розміром 32 байта.

__Get32BytesKey__
```
Get32BytesKey() IKey
```
Повертає структуру IKey яка відповідає за операції із ключами.