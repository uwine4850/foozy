## Package livereload
Цей пакет призначений для спрощення та пришвидшення розробки, він повинен використовуватися лише на стадії розробки.<br>
Коли вносяться зміни у файлах проєкту, його потрібно постійно перезавантажувати. Цей пакет потрібен для того, щоб вирішити 
це завдання і робити це автоматично.<br>
Даний пакет поділяється на дві частини, а саме та частину із перезавантаженням і частиною із прослуховуванням 
збереження файлів.

### Приклад використання
```
package main

import "github.com/uwine4850/foozy/pkg/livereload"

func main() {
	reload := livereload.NewReload("project/cmd/main.go", livereload.NewWiretap([]string{"project", "pkg"},
		[]string{}))
	reload.Start()
}
```
У цьому прикладі перевантажується сервер який знаходиться в файлі ``project/cmd/main.go``. Перезавантаження відбувається 
коли зберігається будь-який файл в директорії ``project`` або ``pkg``.


## Прослуховування збереження файлів
За цей функціонал відповідає інтерфейс ``IWiretap``. Далі будуть описуватися методи які з ним зв'язані.

__SetDirs__
```
SetDirs(dirs []string)
```
Обов'язковий метод. З його допомогою встановлюються директорії в яких будуть прослуховуватися файли.

__OnStart__
```
OnStart(fn func())
```
Метод запускає функцію один раз ``fn`` під час старту прослуховування.

__GetOnStartFunc__
```
GetOnStartFunc() func()
```
Повертає функцію яка буда встановлена методом __OnStart__.

__OnTrigger__
```
OnTrigger(fn func(filePath string))
```
Метод встановлює функцію ``fn`` яка виконується кожен раз при збереженні файлу. Параметр ``filePath string`` це шлях до 
файлу який зберігся.

__SetUserParams__
```
SetUserParams(key string, value interface{})
```
Встановлює параметри які можуть передаватися між функціями ``OnStart(fn func())`` ``OnTrigger(fn func(filePath string))``.

__GetUserParams__
```
GetUserParams(key string) (interface{}, bool)
```
Повертає параметри користувача які встановлені методом __SetUserParams__.

__Start__
```
Start() error
```
Запускає прослуховування.

### Приклад використання
```
wiretap := livereload.NewWiretap3()
wiretap.SetDirs([]string{"project", "project_files"})
wiretap.OnStart(func() {
    fmt.Println("Start.")
})
wiretap.OnTrigger(func(filePath string) {
    fmt.Println("Trigger.")
})
err := wiretap.Start()
if err != nil {
    panic(err)
}
```

## Перезавантаження серверу
Для реалізації цього функціоналу використовується структура ``Reload``.

Конструктор __NewReload(pathToServerFile string, wiretap interfaces.IWiretap) *Reload__<br>
* pathToServerFile - шлях до файла який запускає сервер, наприклад, ``"project/cmd/main.go``.
* wiretap - екземпляр ``interfaces.IWiretap``.

__Start__
```
Start()
```
Запуск перезавантаження сервера.
