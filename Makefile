run:
	go run cmd/service/*

build:
	mkdir -p out && cd out && go build -o service.currency ../cmd/service/*	

tidy:
	go fmt ./...
	go mod tidy
