#!/bin/bash

go run cmd/server/main.go &
SERVER_PID=$!

sleep 1

go run cmd/client/main.go -name=myapp -local=http://localhost:3000 -server=http://localhost:8080

kill $SERVER_PID
