# go-tmsh

[![Build Status](https://travis-ci.org/yukirii/go-tmsh.svg?branch=master)](https://travis-ci.org/yukirii/go-tmsh) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/yukirii/go-tmsh/master/LICENSE)

go-tmsh is a library providing functions to operate the BIG-IP Traffic Management Shell (TMSH) via SSH.

go-tmsh is also a command-line tool that will operate TMSH using these functions.


## Tested versions of BIG-IP

Currently the following versions of BIG-IP are tested.

* v11.5.3
* v11.2.1

## Install

```bash
$ go get github.com/yukirii/go-tmsh/...
```

## Usage

### Using the go-tmsh library

```go
import "github.com/yukirii/go-tmsh"
```

Please refer to the [examples directory](https://github.com/yukirii/go-tmsh/tree/master/examples) for an example source code.

### Using the tmsh command line tool

`tmsh` is single command-line application. This application then takes subcommands. To check the all available commands,

```
$ tmsh help
```

To get help for any specific subcommand, run it with the `-h` flag,

```
$ tmsh node -h
```


## Licence

[MIT](https://github.com/yukirii/go-tmsh/blob/master/LICENSE)
