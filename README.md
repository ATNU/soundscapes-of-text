# Polly

Go Binary that generates text-to-speech using [AWS Polly]

___
## Installation
Requires [go] to build
Install dependancies and build binary

```sh
$ go get "github.com/aws/aws-sdk-go/aws" \
	"github.com/aws/aws-sdk-go/aws/session" \
	"github.com/aws/aws-sdk-go/service/polly" \
	"github.com/spf13/viper"
$ go build
```

AWS authentication must be set up in the $HOME directory 
.
└── .aws
   ├── credentials
   └── config

###### Sample credentials
`[default]`
`aws_access_key_id=AKIAIOSFODNN7EXAMPL`
`aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY`

###### Sample config
`[default]`
`region=us-west-2`
`output=json`

For more info visit [here]
___
## Usage



The binary expects a single command line argument specifying the
file to encode.

```sh
$ ./polly file.txt
```

The following filetypes are accepted:

 - .txt for plaintext
 - .xml for speech synthesis markup language




___
## Configuration

A basic configuration file is provided
 Option        | Description           |
| ------------- |-------------|
| voice      | AWS Polly has a number of available voices to choose from. 
                Specify one here |
| output      | The output format: May be - .mp3 <\br> - .ogg3 <\br> - .pcm |




[AWS Polly]: https://aws.amazon.com/polly/
[go]: https://golang.org/
[here]: https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html