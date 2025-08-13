### Reloader
An object that is responsible for restarting the application. It builds the application into a binary file and launches it; after the trigger, it repeats this process.

It uses [Wiretap](/server/livereload/wiretap) for its work.

#### Reloader.Start
Starts listening for files to change and restarts the server.
```golang
func (r *Reloader) Start() error {
	r.wiretap.OnStart(func() {
		r.onStart()
	})
	r.wiretap.OnTrigger(func(filePath string) {
		r.onTrigger()
	})
	if err := r.wiretap.Start(); err != nil {
		return err
	}
	return nil
}
```

#### Reloader.onStart
Actions that are performed at the start. 
Here the application binary file is built, and then this file is running.
```golang
func (r *Reloader) onStart() {
	cmd := exec.Command("go", "build", "-o", "myapp", r.serverEntryPointPath)
	var stderrBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(&stderrBuf, os.Stderr)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", stderrBuf.String())
	}
	r.serverProcess = exec.Command("./myapp")
	r.serverProcess.Stdout = os.Stdout
	r.serverProcess.Stderr = os.Stderr
	if err := r.serverProcess.Start(); err != nil {
		panic(err)
	}
}
```

#### Reloader.onTrigger
Action during a file reload.<br>
The server stops and then the `onStart` method is called, 
which starts the rebuilt application again.
```golang
func (r *Reloader) onTrigger() {
	if err := r.serverProcess.Process.Kill(); err != nil {
		fmt.Println(err)
	}
	if _, err := r.serverProcess.Process.Wait(); err != nil {
		fmt.Println("err")
	}
	r.onStart()
}
```