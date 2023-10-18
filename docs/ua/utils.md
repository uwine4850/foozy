## Package utils
У цьому пакеті знаходяться функції які не належать конкретному функціоналу, а є глобальними.

__PathExist__
```
PathExist(path string) bool
```
Перевіряє чи існує шлях у файловій системі.

__SplitUrl__
```
SplitUrl(url string) []string
```
Розділяє URL по знаку ``/``.

__SplitUrlFromFirstSlug__
```
SplitUrlFromFirstSlug(url string) string
```
Розділяє URL по першому slug параметру і повертає ліву сторону.

__SliceContains__
```
SliceContains[T comparable](slice []T, item T) bool
```
Перевіряє чи є значення в зрізі.

__GenerateCsrfToken__
```
GenerateCsrfToken() string
```
Генерує CSRF токен.

__MergeMap__
```
MergeMap[T1 comparable, T2 any](map1 *map[T1]T2, map2 map[T1]T2)
```
Об'єднує дві карти в одну(map1).
