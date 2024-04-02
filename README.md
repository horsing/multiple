# Multiple

[![Build Status](https://github.com/horsing/multiple/actions/workflows/go.yml/badge.svg)](https://github.com/horsing/multiple/actions/workflows/ci.yml)
[![LICENSE](https://img.shields.io/github/license/horsing/multiple.svg)](https://github.com/horsing/multiple/blob/master/LICENSE)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/horsing/multiple)](https://goreportcard.com/report/github.com/horsing/multiple)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/2761/badge)](https://bestpractices.coreinfrastructure.org/projects/6232)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/horsing/multiple/badge)](https://securityscorecards.dev/viewer/?uri=github.com/horsing/multiple)
[![Codecov](https://img.shields.io/codecov/c/github/horsing/multiple?style=flat-square&logo=codecov)](https://codecov.io/gh/horsing/multiple)
[![CLOMonitor](https://img.shields.io/endpoint?url=https://clomonitor.io/api/projects/cncf/chubao-fs/badge)](https://clomonitor.io/projects/cncf/chubao-fs)
[![Release](https://img.shields.io/github/v/release/horsing/multiple.svg?color=161823&style=flat-square&logo=smartthings)](https://github.com/horsing/multiple/releases)
[![Tag](https://img.shields.io/github/v/tag/horsing/multiple.svg?color=ee8936&logo=fitbit&style=flat-square)](https://github.com/horsing/multiple/tags)

## Overview

Multiple is an open-source command line tool, which can read input and separate into multiple subtasks to accelerate
large scale jobs' execution.

## Usage

### Introduction

```text
Usage: multiple.exe [options] -t "command template ..."
Available options:
  -h|--help                                         show this help
  --in=..., --in|-i <input>                         read each element from file or stdin
  --sep=..., --sep|-s <separator>                   separator for command and its arguments
  --cpu=..., --cpu|-n <number>                      number of CPU cores to use
  --template=..., --template|-t <command template>  command template to be executed
```

### Batch rename files

```bash
find . -regex ".*\.jpg" | multiple -t 'mv {{.self}} {{trim ".jpg" .self}}.png'
```

### Batch convert files' encoding

```powershell
gci -r -fi "*.java.gbk" | multiple -t 'iconv -f GBK -t UTF8 -o {{trim .self|trim \".gbk\"}} {{trim .self}}'
```

> Former command will find all `.java.gbk` files in current directory recursively, then convert them to `.java` files
> from encoding "GBK" to "UTF8".

### Builtin functions

- `trim`: Trim the prefix or suffix of a string.

> 1. `{{trim .self}}` will trim all white spaces from a string.
> 2. `{{trim "pattern" .self}}` will trim `pattern` from a string(both head and tail).
> 3. `{{trim "prefix" "suffix" .self}}` will trim the prefix `prefix` and suffix `suffix` from a string.

- `add`/`sub`/`mul`/`div`: Arithmetic operations.

## License

Multiple is licensed under the [MIT](https://opensource.org/license/mit).
For detail see [LICENSE](LICENSE).

## Note

The master branch may be in an unstable or even broken state during development. Please use releases instead of the
master branch in order to get a stable set of binaries.

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=horsing/multiple&type=Date)](https://star-history.com/#horsing/multiple&Date)