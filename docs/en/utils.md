## Package utils
This package contains functions that do not belong to a specific functionality, but are global.

__PathExist__
```
PathExist(path string) bool
```
Checks if the path exists in the file system.

__SplitUrl__
```
SplitUrl(url string) []string
```
Separates URLs by the sign ``/``.

__SplitUrlFromFirstSlug__
```
SplitUrlFromFirstSlug(url string) string
```
Splits the URL by the first slug parameter and returns the left side.

__SliceContains__
```
SliceContains[T comparable](slice []T, item T) bool
```
Checks if the value is in the slice.

__GenerateCsrfToken__
```
GenerateCsrfToken() string
```
Generates a CSRF token.

__MergeMap__
```
MergeMap[T1 comparable, T2 any](map1 *map[T1]T2, map2 map[T1]T2)
```
Combines two maps into one (map1).
