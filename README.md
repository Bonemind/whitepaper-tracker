# AWS Whitepaper read tracker

Simple tracker that fetches a list of AWS whitepapers, and allows you to track when you read them.

## Usage

```bash
Usage of whitepaper-tracker
  -config="": path to config file
  -db_location="papers.db": The location of the sqlite db
  -port=3000: Port to listen on
  -skipload=false: Whether to skip the initial whitepaper load, useful for testing
  -test_fetch=false: Whether to test if item fetch still works instead of starting the server
```

Variables can also be passed in via environment variables, and a config file:

Priority is: command line args > env vars > config file > defaults

### Env vars

```shell
DB_LOCATION
PORT
SKIPLOAD
TEST_FETCH
```

### Config file

```
skipload=false
db_location=papers.db
test_fetch=true
port=3333
```

Then point to the config folder with `whitepaper_tracker -config papers.conf`

## TODO

  - Write a frontend
  - Pull newest whitepaper info every day