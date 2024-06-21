## IManagerOneTimeData
Даний менеджер відповідає за передачу та збереження даних окремого запиту.
Тобто, він зберігає дані тільки в межах одного запиту, він не може передавати 
дані до інших запитів. Також цей менеджер може передавати дані із Middleware 
до запиту.

__SetUserContext__
```
SetUserContext(key string, value interface{})
```
Встановлює користувацький контекст, який доступний у межах коректного записту.

__GetUserContext__
```
GetUserContext(key string) (any, bool)
```
Повертає користувацький контекст. Важливо зазначити, що тут можуть бути 
повідомлення із Middleware, наприклад, значення помилки.

__DelUserContext__
```
DelUserContext(key string)
```
Видаляє користувацький контекст по ключу.

__SetSlugParams__
```
SetSlugParams(params map[string]string)
```
Встановлює значення slug параментра. У стандартній реалізвції використовується 
у роутері.

__GetSlugParams__
```
GetSlugParams(key string) (string, bool)
```
Повертає значення slug параментра.