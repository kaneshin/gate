# Gate

Gate posts a message in Slack and LINE.

## Installation

```shell
go get github.com/kaneshin/gate/cmd/...
```

## Setup

Gate loads its configuration in ~/.config/gate/config.yml as a default. You need to create it then setup your channels' configuration:

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
    service-1: '[YOUR-ACCESS-TOKEN]'
```

## Usage

### gate

Run the gate server

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
