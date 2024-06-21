## Package builtin_mddl
Цей пакет містить готові реалізації проміжного ПО.

__GenerateAndSetCsrf__
```
GenerateAndSetCsrf(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)
```
Дана функція це стандартна реалізація middleware.<br>
З допомогою цієї функції можна згенерувати та встановити значення csrf_token у параметри cookies. Приклад застосування:
```
mddl := middlewares.NewMiddleware()
mddl.AsyncHandlerMddl(builtin_mddl.GenerateAndSetCsrf)
```

__GenerateCsrfToken__
```
GenerateCsrfToken()
```
Генерує CSRF токен.