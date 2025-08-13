## wiretap
Starts listening to selected files using the `InotifyEvent` system event. When at least one file save is detected, it calls the trigger function.

#### SetObserverEventDelay
Pause time between the execution of an event for one observation.<br>
This delay applies to each observation separately.
```golang
func SetObserverEventDelay(delay int) {
	observedEventDelay = time.Duration(delay) * time.Millisecond
}
```

### ObservedElement
One element to track, file or directory (all files in it are tracked).

#### ObservedElement.Ready
Shows whether the delay between invents has passed and whether the file 
is ready before the re-invent.
```golang
func (of *ObservedElement) Ready() bool {
	return time.Since(of.LastEventTime) > observedEventDelay
}
```

### Wiretap
An object that monitors files and triggers an event when one of them is saved.

The global algorithm works as follows:

* The user sets up the elements to be monitored.
* The elements to be monitored are processed and start listening using `syscall.SYS_INOTIFY_ADD_WATCH`.
* A cycle is started that waits for `InotifyEvent` and processes it.
* When `InotifyEvent` occurs, the trigger is activated.

#### Wiretap.OnStart
Starts every time during the start of the wiretapping.<br>
The method is usually needed for preliminary initialization before monitoring begins.
```golang
func (w *Wiretap) OnStart(fn func()) {
	w.onStartFn = fn
}
```

#### Wiretap.OnTrigger
A function that will be executed each time the trigger is executed.
```golang
func (w *Wiretap) OnTrigger(fn func(filePath string)) {
	w.onTriggerFn = fn
}
```

#### Wiretap.SetDirs
Sets the directories to be listened to.<br>
It is important to specify that it is the directory and all files in it that is listened to.<br>
One directory and all files in it is one `ObservedElement`.<br>
Subdirectories are already considered new `ObservedElement`.
```golang
func (w *Wiretap) SetDirs(dirs []string) {
	w.dirs = dirs
}
```

#### Wiretap.SetExcludeDirs
Excludes the directory and absolutely all subdirectories from listening.
```golang
func (w *Wiretap) SetExcludeDirs(dirs []string) {
	w.excludeDirs = dirs
}
```

#### Wiretap.SetFiles
Adds individual files to the wiretap.
```golang
func (w *Wiretap) SetFiles(files []string) {
	w.files = files
}
```

#### Wiretap.Start
Starts listening. Performs some initialization actions.
```golang
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
```