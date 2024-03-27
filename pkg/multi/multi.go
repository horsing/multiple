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
	tpl      *template.Template
}

func NewMover(n int, cmd, sep string, tpl *template.Template) Executor {
	workers := int(math.Min(float64(n), float64(runtime.NumCPU()-1)))
	channels := make([]chan string, workers)
	admin := make(chan int, workers)

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
		tpl:      tpl,
	}
}

func processor(i int, biz chan string, admin chan int, t *template.Template, sep string) {
	for {
		item := <-biz
		switch item {
		case ":exit":
			admin <- i
			return
		}

		buf := bytes.Buffer{}
		err := t.Execute(&buf, map[string]interface{}{"self": item})
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

		logger.Info("finish", "id", i, "cmd", raw)
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