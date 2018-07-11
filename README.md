# Polly

Go Binary that generates text-to-speech using [AWS Polly]


## Installation
Requires [go] to build
Install dependancies and build binary

```sh
$ go get "github.com/aws/aws-sdk-go/aws" 		\
	"github.com/aws/aws-sdk-go/aws/session" 	\
	"github.com/aws/aws-sdk-go/service/polly" 	\
	"github.com/spf13/cobra"					\
	"github.com/fsnotify/fsnotify"				\
	"github.com/spf13/viper"
$ go build
```

AWS authentication must be set up in the $HOME directory
```
.
└──.aws
   ├── credentials 
   └── config
```

###### Sample credentials
```sh
[default]
aws_access_key_id=AKIAIOSFODNN7EXAMPL
aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

###### Sample config
```sh
[default]
region=us-west-2
output=json
```

For more info visit [here]


## Usage

```sh
$ ./polly generate
```
To generate a TTS encoding of the input text

```sh
$ ./polly serve
```
Serve starts a listener on the input text, and generates a new TTS encoding everytime the input
text is modified.
This is useful when fine-tuning an SSML file.
Caution as this may impede on AWS tier usage.

```sh
$ ./polly webserver
```
Starts a HTTP server

```sh
$ ./polly get
```
A number of AWS Polly availabilities may be retrieved from the get command.

## Configuration

A basic configuration file is provided

| Option      | Description                                                              |
| ----------- | ------------------------------------------------------------------------ |
| voice       | AWS Polly has a number of available voices to choose from.               |
| input       | Input filepath                                                           |
| outputtype  | Encoding output format:<ul><li>.mp3</li><li>.off3</li><li>.pcm</li></ul> |
| assets      | Asset Paths:<ul><li>demoPath</li><li>ttsPath</li></ul>                   |


[AWS Polly]: https://aws.amazon.com/polly/
[go]: https://golang.org/
[here]: https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html