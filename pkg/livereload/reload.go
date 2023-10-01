package livereload

import (
	"fmt"
	"github.com/uwine4850/foozy/internal/interfaces"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type Reload struct {
	pathToServerFile string
	dirs             []string
	wiretap          interfaces.IWiretap
	exitCh           chan os.Signal
}

func NewReload(pathToServerFile string, dirs []string, wiretap interfaces.IWiretap) *Reload {
	return &Reload{pathToServerFile: pathToServerFile, dirs: dirs, wiretap: wiretap, exitCh: make(chan os.Signal, 1)}
}

func (r *Reload) onStart() {
	cmdB := exec.Command("go", "build", r.pathToServerFile)
	err := cmdB.Run()
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("./main")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	r.wiretap.SetUserParams("cmd", cmd)

	signal.Notify(r.exitCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-r.exitCh:
			c, _ := r.wiretap.GetUserParams("cmd")
			cmd := c.(*exec.Cmd)
			err := cmd.Process.Kill()
			if err != nil {
				panic(fmt.Sprintf("Error killing process: %v\n", err))
			}
		}
	}()
	cmd.Run()
}

func (r *Reload) onTrigger() {
	c, _ := r.wiretap.GetUserParams("cmd")
	cmd := c.(*exec.Cmd)
	err := cmd.Process.Kill()
	if err != nil {
		panic(err)
	}
	log.Println("Server stopped.")
	go r.wiretap.GetOnStartFunc()()
}

func (r *Reload) Start() {
	r.wiretap.SetDirs(r.dirs)
	r.wiretap.OnStart(func() {
		r.onStart()
	})
	r.wiretap.OnTrigger(func(filePath string) {
		r.onTrigger()
	})

	go func() {
		err := r.wiretap.Start()
		if err != nil {
			panic(err)
		}
	}()

	// Waiting for the program to complete through signals.
	r.exitCh = make(chan os.Signal, 1)
	signal.Notify(r.exitCh, syscall.SIGINT, syscall.SIGTERM)
	<-r.exitCh
}
