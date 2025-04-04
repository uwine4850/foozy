package interfaces

type IWiretap interface {
	OnStart(fn func())
	GetOnStartFunc() func()
	OnTrigger(fn func(filePath string))
	SetDirs(dirs []string)
	SetExcludeDirs(dirs []string)
	SetFiles(files []string)
	Start() error
}
