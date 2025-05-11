package livereload

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"syscall"
	"time"
	"unsafe"
)

type InotifyEvent struct {
	Wd     int32
	Mask   uint32
	Cookie uint32
	Len    uint32
	Name   [256]byte
}

// Pause time between the execution of an event for one observation.
// This delay applies to each observation separately.
var observedEventDelay = 1000 * time.Millisecond

func SetObserverEventDelay(delay int) {
	observedEventDelay = time.Duration(delay) * time.Millisecond
}

// ObservedElement the element being monitored.
type ObservedElement struct {
	Wd            int32
	Path          string
	LastEventTime time.Time
}

// Ready shows whether the delay between invents has passed
// and whether the file is ready before the re-invent.
func (of *ObservedElement) Ready() bool {
	return time.Since(of.LastEventTime) > observedEventDelay
}

type Wiretap struct {
	observedElements []*ObservedElement
	globalWD         int32
	initFD           int32
	onStartFn        func()
	onTriggerFn      func(filePath string)
	dirs             []string
	excludeDirs      []string
	files            []string
	trackedElemets   []string
}

func NewWiretap() *Wiretap {
	return &Wiretap{
		observedElements: []*ObservedElement{},
		globalWD:         0,
	}
}

// OnStart starts every time during the start of the wiretapping.
func (w *Wiretap) OnStart(fn func()) {
	w.onStartFn = fn
}

func (w *Wiretap) GetOnStartFunc() func() {
	return w.onStartFn
}

// OnTrigger a function that will be executed each time the trigger is executed.
func (w *Wiretap) OnTrigger(fn func(filePath string)) {
	w.onTriggerFn = fn
}

// SetDirs sets the directories to be listened to.
// It is important to specify that it is the directory and all files in it that is listened to.
// One directory and all files in it is one [ObservedElement].
// Subdirectories are already considered new [ObservedElement].
func (w *Wiretap) SetDirs(dirs []string) {
	w.dirs = dirs
}

// SetExcludeDirs excludes the directory and absolutely all subdirectories from listening.
func (w *Wiretap) SetExcludeDirs(dirs []string) {
	w.excludeDirs = dirs
}

// SetFiles adds individual files to the wiretap.
func (w *Wiretap) SetFiles(files []string) {
	w.files = files
}

func (w *Wiretap) initInotify() error {
	fd, _, err := syscall.Syscall(syscall.SYS_INOTIFY_INIT, 0, 0, 0)
	if err != 0 {
		return err
	}
	w.initFD = int32(fd)
	return nil
}

func (w *Wiretap) addObservedElement(path string, mask uint32) error {
	cPath := []byte(path)
	_, _, errno := syscall.Syscall6(syscall.SYS_INOTIFY_ADD_WATCH, uintptr(w.initFD), uintptr(unsafe.Pointer(&cPath[0])), uintptr(mask), 0, 0, 0)
	if errno != 0 {
		return errno
	}
	w.globalWD++
	w.observedElements = append(w.observedElements, &ObservedElement{
		Wd:   w.globalWD,
		Path: path,
	})
	return nil
}

// collectTrackedElements is gathering elements for a wiretap.
// Directories without files are ignored, but not their subdirectories.
func (w *Wiretap) collectTrackedElements() {
	for i := 0; i < len(w.dirs); i++ {
		filepath.WalkDir(w.dirs[i], func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				if slices.Contains(w.excludeDirs, path) {
					return fs.SkipDir
				}
				if ok, err := hasFiles(path); err != nil {
					return err
				} else if ok {
					w.trackedElemets = append(w.trackedElemets, path)
				}
			}
			return err
		})
	}
	for i := 0; i < len(w.files); i++ {
		w.trackedElemets = append(w.trackedElemets, w.files[i])
	}
}

func (w *Wiretap) processingTrackedElemetns() {
	for i := 0; i < len(w.trackedElemets); i++ {
		w.addObservedElement(w.trackedElemets[i], syscall.IN_MODIFY)
	}
}

// Start starts listening. Performs some initialization actions.
func (w *Wiretap) Start() error {
	if err := w.initInotify(); err != nil {
		return err
	}
	w.collectTrackedElements()
	w.processingTrackedElemetns()
	if w.GetOnStartFunc() != nil {
		w.GetOnStartFunc()()
	}

	var buffer []uint8
	for {
		buffer = make([]byte, syscall.SizeofInotifyEvent+4096)
		n, err := syscall.Read(int(w.initFD), buffer)
		if err != nil {
			return fmt.Errorf("error reading from inotify descriptor: %s", err.Error())
		}

		for i := 0; i < n; {
			event := (*InotifyEvent)(unsafe.Pointer(&buffer[i]))
			observedElement := w.observedElements[event.Wd-1]
			if !observedElement.Ready() {
				i += int(syscall.SizeofInotifyEvent) + int(event.Len)
				continue
			}
			if w.onTriggerFn != nil {
				w.onTriggerFn(observedElement.Path)
			}

			observedElement.LastEventTime = time.Now()
			i += int(syscall.SizeofInotifyEvent) + int(event.Len)
		}
	}
}

func hasFiles(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			return true, nil
		}
	}
	return false, nil
}
