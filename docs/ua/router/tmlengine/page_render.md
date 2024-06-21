## IRender
Даний інтерфейс призначений для спрощення взаємодії із шаблонізатором. Також він 
містить власні методи для відображення даних на сторінці.

__New()__
```
New() (interface{}, error)
```
Імплементація інтерфейсу __INewInstance__. Використовується для автоматичного 
створення нового екземпляру. Якщо шаблонізатор не встановлено в ручну, буде 
використовуватися стандартна реалазіція шаблонізатора.

__SetContext__
```
SetContext(data map[string]interface{})
```
Встановлює контекст для шаблонізатора. У шаблоні можливо викликати встановлені 
дані по ключу.

__GetContext__
```
GetContext() map[string]interface{}
```
Повертає значення контексту шаблонізатора.

__SetTemplateEngine__
```
SetTemplateEngine(engine ITemplateEngine)
```
Встановлює шаблонізатор. Не потрібно викликати для використання стандартної 
реалізації.

__GetTemplateEngine__
```
GetTemplateEngine() ITemplateEngine
```
Повертає шаблонізатор, який використовується.

__RenderTemplate__
```
RenderTemplate(w http.ResponseWriter, r *http.Request) error
```
Налаштовує та запускає шаблонізатор.

__SetTemplatePath__
```
SetTemplatePath(templatePath string)
```
Встановлює шлях до HTML шаблона. Потрібно викликати перед __RenderTemplate__.

__RenderJson__
```
RenderJson(data interface{}, w http.ResponseWriter) error
```
Відображає на сторінці дані у форматі JSON.

## Інші функції

__CreateAndSetNewRenderInstance__
```
CreateAndSetNewRenderInstance(manager interfaces.IManager) error
```
Створює та встановлює у менеджер новий екзкмпляр рендера. Використовується у 
роутері.