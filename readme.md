# Microservice Project

This project has two store and authentication services in its backend and also has a frontend app for the admin panel.

## Usage 
#### FrontEnd
First, go to the `frontend` directory:
```bash
cd ./frontend
```
To install packages and program dependencies:
```bash
npm i
```

Run application
```bash
npm start
```

#### BackEnd
First, go to the `backend` directory:
```bash
cd ./frontend
```
Running the program as a script:
```bash
go run main.go
```
Creating an executable file for the program:
```bash
go build main
```
```bash
go build main
```

#### Run App with Docker
set your environment config, then:
```bash
cd ..
docker compose up -d # or `docker-compose up -d` for older version
