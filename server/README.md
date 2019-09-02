[![Build Status](https://travis-ci.org/ATNU/soundscapes-of-text-webserver.svg?branch=master)](https://travis-ci.org/ATNU/soundscapes-of-text-webserver) [![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/ATNU/soundscapes-of-text-webserver) [![Go Report Card](https://goreportcard.com/badge/github.com/ATNU/soundscapes-of-text-webserver)](https://goreportcard.com/report/github.com/ATNU/soundscapes-of-text-webserver)
# Polly

Service to generate text-to-speech (tts) encodings using [AWS Polly] and serve via a RESTful API


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

Requires AWS authentication details in the $HOME directory, for more AWS credential info, visit [here]
```
.
└──.aws
   ├── credentials 
   └── config
```

###### Sample Credentials
```sh
[default]
aws_access_key_id=AKIAIOSFODNN7EXAMPL
aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

###### Sample Config
```sh
[default]
region=us-west-2
output=json
```

Requires access to AWS S3 Bucket with full access.
Bucket policy must allow public ```GET``` and ```HEAD``` method
###### Sample Bucket Policy
```sh
{
    "Version": "2012-10-17",
    "Id": "Policy1532096226680",
    "Statement": [
        {
            "Sid": "Stmt1532096224403",
            "Effect": "Allow",
            "Principal": "*",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::uk.ac.ncl.sot/*"
        }
    ]
}

```

Use provided configuration to set S3 Bucket name, SNS logging (optional), and http server settings 
```sh
"webserver": {
        "addr": "0.0.0.0:8080",
        "clientAddr": "http://localhost:4200",
        "timeout": {
            "write": 15,
            "read": 15,
            "idle": 60,
            "cancel": 60 
        }
    },
    "s3": {
        "bucketName": "uk.ac.ncl.sot",
        "outputFormat": "mp3",
        "maxRetryCount": 10
    },
    "sns": {
        "pollyTopicName": "arn:aws:sns:us-east-1:438791141487:Polly"
    }
```

## Usage

```sh
$ ./polly webserver
```
Initialises webserver with the following routes:

##### GET /languages
Query all supported languages

###### Response:
```sh
{
    {
        "Name": "[name]",
        "Code": "[code]"
    },
    ...
}
```

##### GET /voices/{languageCode}
Query a specific language code to return all available voices for that language.

###### Response:
```sh
{
    "voices": {
        "Gender": "[gender]"
        "Id": "[id]"
        "LanguageCode": "[languageCode]",
        "LanguageName": "[languageName]",
		"Name": "[name]"
    },
    ...
}
```

##### GET /demo/{voiceID}
Retreive a short demo .mp3 of the queried voice.

###### Response:
```sh
"Content-Type", "audio/mpeg"
```

##### POST /generate/?voice={voideID}
Generate a text-to-speech (tts) encoding of the request body and store in AWS S3. The object URL
is returned once the resource is available.

###### Response:
```sh
https://s3.us-east-1.amazonaws.com/uk.ac.ncl.sot/afd70890-8019-4b5e-90c3-165615727926.mp3
```

[AWS Polly]: https://aws.amazon.com/polly/
[go]: https://golang.org/
[here]: https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html

