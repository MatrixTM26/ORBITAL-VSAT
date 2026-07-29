BINARY  = demon
CMD     = ./cmd/demon
OUT     = bin/$(BINARY)

build:
	go build -ldflags="-s -w" -o $(OUT) $(CMD)

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(OUT)-linux-amd64 $(CMD)

build-arm:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(OUT)-linux-arm64 $(CMD)

tidy:
	go mod tidy

clean:
	rm -rf bin/

run:
	go run $(CMD)

.PHONY: build build-linux build-arm tidy clean run
