## IManagerConfig
Зберігає та надає доступ до деяких налаштувань.

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