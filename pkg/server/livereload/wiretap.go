package livereload

import (
	"errors"
	"github.com/uwine4850/foozy/pkg/utils"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type WiretapFiles struct {
	dirs                []string
	excludeDirs         []string
	files               []string
	onTrigger           func(filePath string)
	onStart             func()
	UserParams          sync.Map
	wg                  sync.WaitGroup
	excludeDeletedFiles []string
	parts               int
}

func NewWiretap(dirs []string, excludeDirs []string) *WiretapFiles {
	return &WiretapFiles{dirs: dirs, excludeDirs: excludeDirs, parts: 10}
}

// OnStart set the function that will be executed once during the runServer of the listening session.
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

func (f *WiretapFiles) SetExcludeDirs(dirs []string) {
	f.excludeDirs = dirs
}

func (f *WiretapFiles) SetNumberOfFileParts(fileParts int) {
	f.parts = fileParts
	if f.parts < 2 {
		panic("parts of file must be greater than 2")
	}
}

// Start starting a wiretap.
func (f *WiretapFiles) Start() error {
	err := f.readDirs()
	if err != nil {
		return err
	}

	f.wg.Add(1)
	go f.onStart()

	var slices [][]string
	lenFiles := len(f.files)
	if lenFiles < f.parts {
		slices = append(slices, f.files)
	} else {
		sliceSize := len(f.files) / f.parts
		for i := 0; i < f.parts; i++ {
			start := i * sliceSize
			end := start + sliceSize
			if i == f.parts-1 {
				end = len(f.files)
			}
			slices = append(slices, f.files[start:end])
		}
	}

	for i := 0; i < len(slices); i++ {
		f.wg.Add(1)
		go func(i int) {
			// Initializing the last modification time map.
			filesModTime := map[string]time.Time{}
			for j := 0; j < len(slices[i]); j++ {
				filePath := slices[i][j]
				filesModTime[filePath] = time.Time{}
			}
			// Run an infinite loop to check file modifications.
			for {
				for j := 0; j < len(slices[i]); j++ {
					filePath := slices[i][j]
					// Skip iteration if file is deleted.
					if utils.SliceContains(f.excludeDeletedFiles, filePath) {
						continue
					}
					err := f.watchFile(filePath, &filesModTime)
					if err != nil {
						if errors.Is(err, fs.ErrNotExist) {
							f.excludeDeletedFiles = append(f.excludeDeletedFiles, filePath)
							return
						}
						panic(err)
					}
				}
			}
		}(i)
	}
	f.wg.Wait()
	return nil
}

func (f *WiretapFiles) watchFile(filePath string, filesModTime *map[string]time.Time) error {
	var init bool
	if (*filesModTime)[filePath].String() == "0001-01-01 00:00:00 +0000 UTC" {
		init = true
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	// When the lastModTime variable is only initialized, it should be set to the correct time to modify the file.
	// Otherwise, lastModTime and file modification time will not coincide and the event will be filled.
	if init {
		(*filesModTime)[filePath] = fileInfo.ModTime()
		init = false
		return nil
	}
	// If the lastModTime and the current time of file modification do not coincide,
	// it means that the file has been modified.
	// The diff variable is responsible for the time difference between the last file save and the current value.
	// If the time is greater than 1.5 seconds and the above conditions are met, the trigger will be fired.
	// This is to prevent the trigger from being triggered several times, because some editors save a file several times.
	diff := fileInfo.ModTime().Sub((*filesModTime)[filePath])
	if fileInfo.ModTime() != (*filesModTime)[filePath] && diff > 1000*time.Millisecond {
		(*filesModTime)[filePath] = fileInfo.ModTime()
		f.onTrigger(filePath)
	}
	return nil
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
			} else {
				if utils.SliceContains(f.excludeDirs, path) {
					return filepath.SkipDir
				}
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
