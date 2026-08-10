BINARY  = vsat
CMD     = ./cmd/vsat
OUT     = bin/$(BINARY)

build:
	@mkdir -p bin
	go build -ldflags="-s -w" -o $(OUT) $(CMD)

build-arm:
	@mkdir -p bin
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(OUT)-arm64 $(CMD)

run:
	go run $(CMD)

clean:
	rm -rf bin/

.PHONY: build build-arm run clean
