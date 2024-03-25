package main

import (
	"bufio"
	"html/template"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/horsing/multiple/pkg/multi"
)

func usage(code int) {
	c := os.Args[0]
	help := `Usage: {{.c}} [options] "command args template ..."
Available options:
  -h|--help                                show this help
  --in=<input>, --in|-i <input>            read each element from file or stdin
  --sep=<separator>, --sep|-s <separator>  separator for command and its arguments
  --cpu=<number>, --cpu|-n <number>        number of CPU cores to use
`
	e := template.Must(template.New("help").Parse(help)).Execute(os.Stdout, map[string]string{"c": c})
	if e != nil {
		panic(e)
	}
	os.Exit(code)
}

func open(f string) (file *os.File, fn func(*os.File)) {
	file, err := os.OpenFile(f, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	fn = func(f *os.File) { f.Close() }
	return
}

func parse(flag bool, s, l string, arg string, args ...string) (match bool, v *string, offset int) {
	if arg == s || arg == l {
		if flag {
			return true, nil, 0
		} else {
			return true, &args[0], 1
		}
	} else if strings.HasPrefix(arg, l+"=") {
		v := arg[len(l)+1:]
		return true, &v, 1
	}
	return false, nil, 0
}

type config struct {
	in      *os.File
	sep     string
	number  int64
	command string
}

var cfg = &config{os.Stdin, " ", int64(runtime.NumCPU()) - 1, "echo {{.file}}"}

func main() {
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case arg == "-h" || arg == "--help":
			usage(0)
		case arg == "-i" || arg == "--in" || strings.HasPrefix(arg, "--in="):
			m, v, o := parse(false, "-i", "--in", arg, os.Args[i+1:]...)
			if !m {
				usage(1)
			}
			i += o
			switch *v {
			case "-", "stdin", "/dev/stdin", "0", "/dev/0":
				cfg.in = os.Stdin
				break
			default:
				in, fn := open(arg)
				cfg.in = in
				defer fn(in)
			}
			break
		case arg == "-s" || arg == "--sep" || strings.HasPrefix(arg, "--sep="):
			m, v, o := parse(false, "-s", "--sep", arg, os.Args[i+1:]...)
			if !m {
				usage(1)
			}
			i += o
			cfg.sep = *v
			break
		case arg == "-n" || arg == "--cpu" || strings.HasPrefix(arg, "--cpu="):
			m, v, o := parse(false, "-n", "--cpu", arg, os.Args[i+1:]...)
			if !m {
				usage(1)
			}
			i += o
			cfg.number, _ = strconv.ParseInt(*v, 10, 64)
		case arg == "-t" || arg == "--tpl" || strings.HasPrefix(arg, "--tpl="):
			m, v, o := parse(false, "-t", "--tpl", arg, os.Args[i+1:]...)
			if !m {
				usage(1)
			}
			i += o
			cfg.command = *v
			break
		default:
			println("Unknown argument", arg)
		}
	}

	var mover = multi.NewMover(int(cfg.number), cfg.command, cfg.sep)
	defer mover.Stop()

	if mover != nil {
		scanner := bufio.NewScanner(cfg.in)
		for scanner.Scan() {
			mover.Submit(scanner.Text())
		}
	}
}
