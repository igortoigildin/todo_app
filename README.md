## Requirements

SQLite installed
Go installed

#### Please use the following config to run tests:

<!-- start:code block -->

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = false
var Search = true
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.UHVCm6mMM4NlVujjwVPVmP6hwq4n31MUd7Z-MFW2yao`

<!-- end:code block -->

#### Please see below lists of constants used in app:

<!-- start:code block -->

const yymmdd = "20060102" // date format constant
const ddmmyy = "02012006" // date format constant
const jwtSecret = "your-secret-key" // string for JWT secret
const Limit = 30 // limit rows for db queries results

<!-- end:code block -->

## Usage

#### To run this application, execute:

<!-- start:code block -->

go run cmd/server/main.go

<!-- end:code block -->

You should be able to access this application at http://127.0.0.1:7540

#### To build Docker Container, execute:

<!-- start:code block -->

docker build --tag my-app:v1 .

<!-- end:code block -->

#### To run your Docker Container, execute:

<!-- start:code block -->

docker run -d --rm -p 7540:7540 my-app:v1

<!-- end:code block -->

#### To run tests, execute:

<!-- start:code block -->

go test ./tests

<!-- end:code block -->
