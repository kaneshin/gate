# Gate

Gate posts a message in Slack and LINE.

## Installation

To compile the Gate binaries from source, clone the Gate repository. Then, navigate to the new directory.

```shell
$ git clone https://github.com/kaneshin/gate.git
$ cd gate
```

Compile the `gate` and `gatecli` which will be stored it in $GOPATH/bin.

```shell
$ go install
```

Finally, make sure that the gate and gatecli binaries are available on your PATH.

## Setup

Gate loads its configuration in ~/.config/gate/config.json as a default. You need to create it then setup your channels' configuration:

```json
{
  "gate": {
    "scheme": "http",
    "host": "0.0.0.0",
    "port": 5731,
    "client": {
      "default": "slack.channel-1"
    }
  },
  "platforms": {
    "slack": {
      "channel-1": "[YOUR-INCOMING-URL]",
      "channel-2": "[YOUR-INCOMING-URL]"
    },
    "line": {
      "service-1": "[YOUR-ACCESS-TOKEN]"
    }
  }
}
```

## Usage

### gate

Run the gate server. If neeeded, run the command with argument of configuration path.

```shell
$ gate -config=/path/to/config.json
```

### gatecli

First, you need to install the configuration of gate to run `gatecli` via network. You will get the partial configuration of config.json.

```shell
$ curl -sL http://0.0.0.0:5731/config/cli.json > ~/.config/gate/cli.json
```

Then, run `gatecli` to post the platforms.

```shell
$ echo "foobar" | gatecli
```

If you'd like to parse another cli.json and notify specified target. Please run the command with arguments like the below.

```
$ echo "foobar" | gatecli -config=/path/to/cli.json -target=slack.channel-2
```

## License

[The MIT License (MIT)](http://kaneshin.mit-license.org/)

## Author

Shintaro Kaneko <kaneshin0120@gmail.com>
