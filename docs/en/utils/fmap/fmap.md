## Package fmap

__MergeMap__
```go
MergeMap[T1 comparable, T2 any](map1 *map[T1]T2, map2 map[T1]T2)
```
Combines two maps into one (map1).

__MergeMapSync__
```go
MergeMapSync[T1 comparable, T2 any](mu *sync.Mutex, map1 *map[T1]T2, map2 map[T1]T2)
```
Combines two maps into one (map1). This operation is performed safely in asynchronous mode.

__Compare__
```go
Compare[T1 comparable, T2 comparable](map1 *map[T1]T2, map2 *map[T1]T2, exclude []T1) bool
```
Compares each element of two cards against each other. 

__YamlMapToStruct__
```go
YamlMapToStruct(targetMap *map[string]interface{}, targetStruct typeopr.IPtr) error
```
Writes a yaml map to the structure.
IMPOrTANT: the field of the structure to be written must have the yaml tag:"<field_name>". This 
name must correspond to the name of the field in the `targetMap` structure. 
Works in depth, you can make as many attachments as you want.