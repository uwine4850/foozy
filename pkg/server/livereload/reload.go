package livereload

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/uwine4850/foozy/pkg/config"
	"github.com/uwine4850/foozy/pkg/interfaces"
)

type Reload struct {
	pathToServerFile string
	wiretap          interfaces.IWiretap
	wg               sync.WaitGroup
	cmd              *exec.Cmd
	cmdStart         *exec.Cmd
}

func NewReload(pathToServerFile string, wiretap interfaces.IWiretap) *Reload {
	return &Reload{pathToServerFile: pathToServerFile, wiretap: wiretap}
}

func (r *Reload) Start() {
	if !config.LoadedConfig().Default.Debug.Debug {
		fmt.Println("Debug is disabled, reloader is not working.")
		return
	}
	r.wiretap.OnStart(func() {
		r.onStart()
	})
	r.wiretap.OnTrigger(func(filePath string) {
		r.onTrigger()
	})
	err := r.wiretap.Start()
	if err != nil {
		panic(err)
	}
}

func (r *Reload) onStart() {
	r.cmd = exec.Command("go", "build", "-p", "4", r.pathToServerFile)
	var stderrBuf bytes.Buffer
	r.cmd.Stderr = io.MultiWriter(&stderrBuf, os.Stderr)
	if err := r.cmd.Run(); err != nil {
		fmt.Println("Error:", stderrBuf.String())
	}
	binaryFileName := strings.Split(filepath.Base(r.pathToServerFile), ".")[0]
	r.cmdStart = exec.Command("./" + binaryFileName)
	r.cmdStart.Stdout = os.Stdout
	r.cmdStart.Stderr = os.Stderr
	r.wg.Add(1)
	go r.runServer()
}

func (r *Reload) onTrigger() {
	if r.cmdStart.Process == nil {
		return
	}
	err := r.cmdStart.Process.Kill()
	if err != nil && err.Error() != "os: process already finished" {
		if err.Error() == "os: process already finished" {
			fmt.Println("Stop server")
			r.wg.Wait()
			r.wiretap.GetOnStartFunc()()
			return
		}
		panic(err)
	}
	if r.cmdStart.Process != nil {
		err = r.cmdStart.Wait()
		if err != nil && err.Error() != "exec: Wait was already called" {
			if err.Error() == "signal: killed" {
				r.wg.Wait()
				r.wiretap.GetOnStartFunc()()
				return
			}
			panic(err)
		}
	}
}

func (r *Reload) runServer() {
	defer r.wg.Done()
	err := r.cmdStart.Start()
	if err != nil {
		if err.Error() == "signal: killed" {
			fmt.Println("STOP SERVER")
			return
		}
		panic(err)
	}
}
