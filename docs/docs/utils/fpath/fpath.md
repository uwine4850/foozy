## fpath

#### CurrentFileDir
Returns the path to the directory in which this function is called.
```golang
func CurrentFileDir() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Dir(file)
}
```

#### PathExist
Checks to see if a path exists in the file directory.
```golang
func PathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
```