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

Please see below lists of constants used in app:
const yymmdd = "20060102" // date format constant
const ddmmyy = "02012006" // date format constant
const jwtSecret = "your-secret-key" // string for JWT secret
const Limit = 30 // limit rows for db queries results

Usage
To run this application, execute:
go run cmd/server/main.go
You should be able to access this application at http://127.0.0.1:7540

To build Docker Container, execute:
docker build --tag my-app:v1 .
To run your Docker Container, execute:
docker run -d --rm -p 7540:7540 my-app:v1

To run tests, execute:
go test ./tests
