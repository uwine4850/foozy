## package secure
Даний розділ безпекового пакету відповідає за операції із CSRF токеном.

__ValidateFormCsrfToken__
```
ValidateFormCsrfToken(r *http.Request, frm *form.Form) error
```
Метод який проводить валідацію CSRF token. Для цього форма повинна мати поле із назвою ``csrf_token``, крім того дані cookies
також повинні мати поле ``csrf_token``.<br>
Найпростіший спосіб для цього - додати вбудований middleware [csrf](https://github.com/uwine4850/foozy/blob/master/docs/ua/builtin/builtin_mddl/csrf.md) який автоматично буде додавати поле ``csrf_token`` в дані cookies.
Після цього просто потрібно додати в середину HTML форми змінну ``{{ csrf_token | safe }}`` та запустити даний метод.<br>
Підключення вбудованого middleware для створення токена буде відбуватись наступним чином:
```
mddl := middlewares.NewMiddleware()
mddl.AsyncHandlerMddl(builtin_mddl.GenerateAndSetCsrf)
...
newRouter.SetMiddleware(mddl)
```

__ValidateHeaderCSRFToken__
```
ValidateHeaderCSRFToken(r *http.Request, tokenName string)
```
Метод який проводить валідацію CSRF token. Для цього потрібно 
щоб токен передавався через header. Також потрібно щоб токен 
був у cookies перед використанням цього метода.