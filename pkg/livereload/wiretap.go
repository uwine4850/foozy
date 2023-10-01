package livereload

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type WiretapFiles struct {
	dirs       []string
	files      []string
	onTrigger  func(filePath string)
	onStart    func()
	UserParams sync.Map
	wg         sync.WaitGroup
}

func NewWiretap3() *WiretapFiles {
	return &WiretapFiles{onTrigger: func(filePath string) {}, onStart: func() {}}
}

// OnStart set the function that will be executed once during the start of the listening session.
func (f *WiretapFiles) OnStart(fn func()) {
	f.onStart = fn
}

func (f *WiretapFiles) GetOnStartFunc() func() {
	return f.onStart
}

// OnTrigger a function that will be executed each time the trigger is executed.
func (f *WiretapFiles) OnTrigger(fn func(filePath string)) {
	f.onTrigger = fn
}

// SetUserParams set the parameter to be passed between methods.
func (f *WiretapFiles) SetUserParams(key string, value interface{}) {
	f.UserParams.Store(key, value)
}

func (f *WiretapFiles) GetUserParams(key string) (interface{}, bool) {
	res, ok := f.UserParams.Load(key)
	return res, ok
}

// SetDirs set the directories whose files will be listened to.
func (f *WiretapFiles) SetDirs(dirs []string) {
	f.dirs = dirs
}

// Start starting a wiretap.
func (f *WiretapFiles) Start() error {
	err := f.readDirs()
	if err != nil {
		return err
	}

	f.wg.Add(1)
	go f.onStart()

	for i := 0; i < len(f.files); i++ {
		f.wg.Add(1)
		filePath := f.files[i]
		go func() {
			err := f.watchFile(filePath, &f.wg)
			if err != nil {
				panic(err)
			}
		}()
	}
	f.wg.Wait()
	return nil
}

// watchFile listening to an individual file.
func (f *WiretapFiles) watchFile(filePath string, wg *sync.WaitGroup) error {
	defer wg.Done()

	lastModTime := time.Time{}
	var init bool
	if lastModTime.String() == "0001-01-01 00:00:00 +0000 UTC" {
		init = true
	}

	for {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return err
		}
		// When the lastModTime variable is only initialized, it should be set to the correct time to modify the file.
		// Otherwise, lastModTime and file modification time will not coincide and the event will be filled.
		if init {
			lastModTime = fileInfo.ModTime()
			init = false
			continue
		}
		// If the lastModTime and the current time of file modification do not coincide,
		// it means that the file has been modified.
		if fileInfo.ModTime() != lastModTime {
			lastModTime = fileInfo.ModTime()
			f.onTrigger(filePath)
		}
	}
}

// readDirs reads directories and writes the absolute path of each file to the slice.
func (f *WiretapFiles) readDirs() error {
	for i := 0; i < len(f.dirs); i++ {
		visit := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				abs, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				f.files = append(f.files, abs)
			}
			return nil
		}

		// Recursively traverse all files in the root folder and its subfolders.
		err := filepath.Walk(f.dirs[i], visit)
		if err != nil {
			return err
		}
	}
	return nil
}
