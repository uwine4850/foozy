package liverealod

import (
	"errors"
	"github.com/uwine4850/foozy/internal/server"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type IWiretap interface {
	OnStart(fn func())
	GetOnStartFunc() func()
	OnTrigger(fn func(filePath string))
	SetUserParams(key string, value interface{})
	GetUserParams(key string) (interface{}, bool)
	SetDirs(dirs []string)
	Start() error
	Stop()
}

type WiretapFiles struct {
	dirs       []string
	files      []string
	onTrigger  func(filePath string)
	StartFunc  func()
	UserParams sync.Map
	Server     *server.Server
	wg         sync.WaitGroup
	done       chan bool
}

func NewWiretap() *WiretapFiles {
	return &WiretapFiles{onTrigger: func(filePath string) {}, StartFunc: func() {}}
}

// OnStart set the function that will be executed once during the start of the listening session.
func (f *WiretapFiles) OnStart(fn func()) {
	f.StartFunc = fn
}

func (f *WiretapFiles) GetOnStartFunc() func() {
	return f.StartFunc
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
	f.checkPlatform()
	err := f.readDirs()
	if err != nil {
		return err
	}
	f.wg.Add(1)
	go f.StartFunc()
	for _, filePath := range f.files {
		f.wg.Add(1)
		filePath := filePath
		go func() {
			err := f.watchFile(filePath, &f.wg)
			if err != nil {
				panic(err)
			}
		}()
	}
	<-f.done
	f.wg.Wait()
	return nil
}

func (f *WiretapFiles) Stop() {
	f.done <- true
}

// watchFile listening to an individual file.
func (f *WiretapFiles) watchFile(filePath string, wg *sync.WaitGroup) error {
	defer wg.Done()

	for {
		cmd := exec.Command("inotifywait", "-e", "modify", filePath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return errors.New(string(output))
		}

		f.onTrigger(filePath)

		time.Sleep(1 * time.Second)
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

func (f *WiretapFiles) checkPlatform() {
	if runtime.GOOS != "linux" {
		panic("At the moment, wiretap only supports the Linux platform.")
	}
}
