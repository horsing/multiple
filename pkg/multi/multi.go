package multi

import (
	"bytes"
	"html/template"
	"math"
	"os/exec"
	"runtime"
	"strings"

	"github.com/horsing/multiple/pkg/log"
)

var logger = log.New("multi")

type Executor interface {
	Submit(files ...string)
	Stop()
}

type mover struct {
	size     int
	channels []chan string
	pos      int
	admin    chan int
}

func NewMover(n int, cmd, sep string) Executor {
	workers := int(math.Min(float64(n), float64(runtime.NumCPU()-1)))
	channels := make([]chan string, workers)
	admin := make(chan int, workers)

	tpl, err := template.New("").Funcs(map[string]interface{}{
		"add": func(l, r int) int { return l + r },
		"sub": func(l, r int) int { return l - r },
		"mul": func(l, r int) int { return l * r },
		"div": func(l, r int) int { return l / r },
	}).Parse(cmd)
	if err != nil {
		panic(err)
	}

	for i := 0; i < workers; i++ {
		channels[i] = make(chan string, 100)
		go processor(i, channels[i], admin, tpl, sep)
		logger.Info("new worker", "id", i)
	}

	return &mover{
		size:     workers,
		channels: channels,
		pos:      0,
		admin:    admin,
	}
}

func processor(i int, biz chan string, admin chan int, t *template.Template, sep string) {
	for {
		file := <-biz
		switch file {
		case ":exit":
			admin <- i
			return
		}

		logger.Info("process", "id", i, "file", file)

		buf := bytes.Buffer{}
		err := t.Execute(&buf, map[string]interface{}{"file": file})
		if err != nil {
			logger.Error(err, "render cmd", "id", i)
			continue
		}
		raw := buf.String()
		cmd := strings.Split(raw, sep)
		c := exec.Command(cmd[0], cmd[1:]...)
		if err = c.Start(); err != nil {
			logger.Error(err, "exec", "id", i, "cmd", raw)
			panic(err)
		}

		logger.Info("finish", "id", i, "command", raw)
	}
}

func (m *mover) Submit(files ...string) {
	for _, file := range files {
		m.channels[m.pos] <- file
		logger.Info("submit", "id", m.pos, "file", file)
		m.pos = (m.pos + 1) % len(m.channels)
	}
}

func (m *mover) Stop() {
	for _, ch := range m.channels {
		ch <- ":exit"
	}
	stopped := make([]int, 0, m.size)
	for {
		i := <-m.admin
		logger.Info("Worker stopped", "worker", i)
		stopped = append(stopped, i)
		if len(stopped) >= m.size {
			break
		}
	}
	logger.Info("All workers have been stopped. Goodbye!")
}
