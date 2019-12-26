# Gate

[![GoDoc](https://godoc.org/github.com/kaneshin/gate?status.svg)](https://godoc.org/github.com/kaneshin/gate)
[![codecov](https://codecov.io/gh/kaneshin/gate/branch/master/graph/badge.svg)](https://codecov.io/gh/kaneshin/gate)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaneshin/gate)](https://goreportcard.com/report/github.com/kaneshin/gate)

## Installation

```shell
go get github.com/kaneshin/gate/cmd/...
```

## Usage

Create `~/.config/gate/config.yml` to be loaded by tools.

```yml
---
default_target: 'slack.incoming.channel-1'

env:
  host: 'http://127.0.0.1'
  port: 8080

slack:
  incoming:
    channel-1: '[YOUR-INCOMING-URL]'
    channel-2: '[YOUR-INCOMING-URL]'

line:
  notify:
    room-1: '[YOUR-ACCESS-TOKEN]'
```

### gate

```shell
$ gate

# gate -config=/path/to/config.yml
```

### gatecli

```shell
$ echo "foobar" | gatecli

# echo "foobar" | gatecli -config=/path/to/config.yml -target=slack.incomiong.channel-2
```

## License

[The MIT License (MIT)](http://kaneshin.mit-license.org/)

## Author

Shintaro Kaneko <kaneshin0120@gmail.com>
