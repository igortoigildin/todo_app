## Introduction

A simple todolist application written in Go.

#### Requirements

SQLite installed

Go installed

#### Please use the following config to run tests:

var Port = 7540

var DBFile = "../scheduler.db"

var FullNextDate = false

var Search = true

var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.UHVCm6mMM4NlVujjwVPVmP6hwq4n31MUd7Z-MFW2yao`

## Usage

To run this application, execute:

```shell
go run cmd/server/main.go
```

You should be able to access this application at http://127.0.0.1:7540

#### To build Docker Container, execute:

<!-- start:code block -->

```shell
docker build --tag my-app:v1 .
```

<!-- end:code block -->

#### To run Docker Container, execute:

<!-- start:code block -->

```shell
docker run -d -p 7540:7540 my-app:v1
```

<!-- end:code block -->

#### To run tests, execute:

<!-- start:code block -->

```shell
go test ./tests
```

<!-- end:code block -->
