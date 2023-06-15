build:
	go build -o build/nomad-demo cmd/nomad-demo/main.go

clean:
	rm -rf build

run:
	go run cmd/nomad-demo/main.go

unit_test:
	go test -race -v ./...
