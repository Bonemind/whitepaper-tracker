# AWS Whitepaper read tracker

Simple tracker that fetches a list of AWS whitepapers, and allows you to track when you read them.

## Usage

The project is split up in a `backend` and `frontend` folder. The backend is a go app that will both serve the frontend, and fetch the aws whitepapers, store them in a sqlite db.

Running the app assumes you've built the frontend, and copied the contents of `frontend/public` to `backend/frontend`

### Server usage

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

#### Env vars

```shell
DB_LOCATION
PORT
SKIPLOAD
TEST_FETCH
```

#### Config file

```
skipload=false
db_location=papers.db
test_fetch=true
port=3333
```

Then point to the config folder with `whitepaper_tracker -config papers.conf`

### Building the backend

```bash
cd backend && go build .
```

### Building the frontend

```bash
cd frontend && npm i && npm run build
```

For development, you can use `npm run dev` to get a livereloading variant for the frontend.

During development and building, you can point to the backend using the `API_BASE_PATH` env var.

For example:

```bash
export API_BASE_PATH=http://localhost:3000/api
```
If you have the backend running using the default configuration, this should point you to that backend.

## Docker image

A Docker image is available that will build the back and frontend, and will run the app when started.

### Building

From the root of the repository:

```bash
sudo docker build . -t whitepaper-tracker:latest
```

### Running the app

```bash
sudo docker run -it --rm -p 3000:3000 -v $PWD:/data -e DB_LOCATION="/data/papers.db" whitepaper-tracker:latest
```