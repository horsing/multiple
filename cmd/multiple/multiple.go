package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/horsing/multiple/pkg/multi"
)

func usage(code int) {
	c := os.Args[0]
	if i := strings.LastIndex(c, string(os.PathSeparator)); i >= 0 {
		c = c[i+1:]
	}
	help := `Usage: {{.c}} [options] -t "command template ..."
Available options:
  -h|--help                                         show this help
  --in=..., --in|-i <input>                         read each element from file or stdin
  --sep=..., --sep|-s <separator>                   separator for command and its arguments
  --cpu=..., --cpu|-n <number>                      number of CPU cores to use
  --template=..., --template|-t <command template>  command template to be executed
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
	fn = func(f *os.File) { _ = f.Close() }
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

func (c config) String() string {
	return "in: [" + c.in.Name() + "], sep: [" + c.sep + "], number: [" + strconv.FormatInt(c.number, 10) + "], command: [" + c.command + "]"
}

var cfg = &config{os.Stdin, " ", int64(runtime.NumCPU()) - 1, "echo {{.self}}"}

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
		case arg == "-t" || arg == "--template" || strings.HasPrefix(arg, "--template="):
			m, v, o := parse(false, "-t", "--template", arg, os.Args[i+1:]...)
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

	tpl, err := template.New("").Funcs(map[string]interface{}{
		"add": func(l, r int) int { return l + r },
		"sub": func(l, r int) int { return l - r },
		"mul": func(l, r int) int { return l * r },
		"div": func(l, r int) int { return l / r },
		"trim": func(args ...string) (r string) {
			size := len(args)
			switch {
			case size == 1:
				// trim s
				r = args[0]
				r = strings.TrimSpace(r)
			case size == 2:
				// trim pat s
				r = args[1]
				r = strings.TrimPrefix(r, args[0])
				r = strings.TrimSuffix(r, args[0])
			case size >= 3:
				// trim prefix suffix s
				r = args[2]
				r = strings.TrimPrefix(r, args[0])
				r = strings.TrimSuffix(r, args[1])
			default:
				r = ""
			}
			return
		},
	}).Parse(cfg.command)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Starting job: %s\n", cfg.String())

	var mover = multi.NewMover(int(cfg.number), cfg.command, cfg.sep, tpl)
	defer mover.Stop()

	scanner := bufio.NewScanner(cfg.in)
	for scanner.Scan() {
		mover.Submit(scanner.Text())
	}
}