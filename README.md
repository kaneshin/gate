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
    },
    "pixela": {
      "username/graph-id": "[YOUR-TOKEN]"
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
$ echo "foobar" | gatecli slack.channel-1
$ echo "foobar" | gatecli # You can send a message without argument for the default target in cli.json
```

If you'd like to parse another cli.json. Please run the command with arguments like the below.

```
$ echo "foobar" | gatecli -config=/path/to/cli.json slack.channel-2
```

#### Slack

You need to get an incoming webhook url what channel you want to post message before run it.

```shell
echo "Hello world!" | gatecli slack.channel-1
```

#### LINE

You need to create a service what channel you want to post message before run it.

```shell
echo "Hello world!" | gatecli line.service-1
```

#### Pixela

You need to create a graph what you want to post quantity before run it.

Input a quantity from STDIN, then pass it to `gatecli`.

```shell
echo 5 | gatecli pixela.username/graph-id
```

## License

[The MIT License (MIT)](http://kaneshin.mit-license.org/)

## Author

Shintaro Kaneko <kaneshin0120@gmail.com>
