.PHONY: run

run:
	go mod download
	go build -buildvcs=false -o app
		./app
