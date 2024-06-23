Introduction

A simple todolist application written in Go

 
Requirements
SQLite installed
Go installed


Task of increased complexity performed:

- env variables;
- /api/tasks?search="";
- authentication with JWT;
- Dockerfile;


Please use the following config to run tests:

var Port = 7540

var DBFile = "../scheduler.db"

var FullNextDate = false

var Search = true

var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.UHVCm6mMM4NlVujjwVPVmP6hwq4n31MUd7Z-MFW2yao`


Usage

To run this application, execute:

go run cmd/server/main.go

You should be able to access this application at http://127.0.0.1:7540
