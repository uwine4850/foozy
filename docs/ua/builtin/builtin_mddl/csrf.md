## Package builtin_mddl
Цей пакет містить готові реалізації проміжного ПО.

__GenerateAndSetCsrf__
```
GenerateAndSetCsrf(maxAge int, onError onError) func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)
```
Дана функція повертає стандартну реалізацію middleware.<br>
maxAge - час життя cookie.
onError - функція, яка буду виконана під час помилки.
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