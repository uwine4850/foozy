package interfaces

type IWiretap interface {
	OnStart(fn func())
	GetOnStartFunc() func()
	OnTrigger(fn func(filePath string))
	SetUserParams(key string, value interface{})
	GetUserParams(key string) (interface{}, bool)
	SetDirs(dirs []string)
	SetExcludeDirs(dirs []string)
	Start() error
}
