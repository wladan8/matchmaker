.PHONY: unit-tests build docker_build docker_run

unit-tests:
	go test -count=1 -race ./...

build:
	go build -o matchmaker ./cmd

docker_build:
	docker build -t matchmaker --no-cache .

docker_run:
	docker  run -e SERVER_PORT=8080 -e GROUP_SIZE=10 -e DIFF_SKILL=10 -e DIFF_LATENCY=100  -e TICKER_FREQUENCY=1 -p 8080:8080 -t matchmaker
