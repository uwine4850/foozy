package livereload

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Reloader struct {
	serverEntryPointPath string
	wiretap              *Wiretap
	serverProcess        *exec.Cmd
}

func NewReloader(serverEntryPointPath string, wiretap *Wiretap) *Reloader {
	return &Reloader{
		serverEntryPointPath: serverEntryPointPath,
		wiretap:              wiretap,
	}
}

// Start starts listening for files to change and restarts the server.
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

// onStart actions that are performed at the start.
// Here the application binary file is built, and then this file is running.
func (r *Reloader) onStart() {
	cmd := exec.Command("go", "build", "-o", "myapp", r.serverEntryPointPath)
	var stderrBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(&stderrBuf, os.Stderr)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", stderrBuf.String())
	}
	binaryFileName := strings.Split(filepath.Base(r.serverEntryPointPath), ".")[0]
	r.serverProcess = exec.Command("./" + binaryFileName)
	r.serverProcess.Stdout = os.Stdout
	r.serverProcess.Stderr = os.Stderr
	if err := r.serverProcess.Start(); err != nil {
		panic(err)
	}
}

// onTrigger action during a file reload.
// The server stops and then the [onStart] method is called,
// which starts the rebuilt application again.
func (r *Reloader) onTrigger() {
	if err := r.serverProcess.Process.Kill(); err != nil {
		fmt.Println(err)
	}
	if _, err := r.serverProcess.Process.Wait(); err != nil {
		fmt.Println("err")
	}
	r.onStart()
}
