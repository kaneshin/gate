# Gate

[![GoDoc](https://godoc.org/github.com/kaneshin/gate?status.svg)](https://godoc.org/github.com/kaneshin/gate)
[![Build Status](https://travis-ci.org/kaneshin/gate.svg?branch=master)](https://travis-ci.org/kaneshin/gate)
[![codecov](https://codecov.io/gh/kaneshin/gate/branch/master/graph/badge.svg)](https://codecov.io/gh/kaneshin/gate)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaneshin/gate)](https://goreportcard.com/report/github.com/kaneshin/gate)

## Installation

```shell
go get -d github.com/kaneshin/gate/cmd/...
```

## Usage

### gate

```shell
gate -config=/path/to/conf.toml -port=8080
```

conf.toml

```toml
[slack]
incoming_url = "[your-incoming-url]"
channel      = "general"
username     = "gate"
icon_emoji   = ":ghost:"

[line]
access_token = "[your-access-token]"

[facebook.messenger]
id = "[sender-id]"
access_token = "[page-access-token]"
```

### gatecli

```shell
echo "foobar" | gatecli -host=http://localhost:8080
```

## License

[The MIT License (MIT)](http://kaneshin.mit-license.org/)

## Author

Shintaro Kaneko <kaneshin0120@gmail.com>
