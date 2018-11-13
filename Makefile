.PHONY: build
build: rice test
	# Incrustate the static files
	cd pkg/handler ; rice embed-go
	go build -a -ldflags "-s -w" -o urlshortener cmd/urlshortener/main.go

.PHONY: rice
	go get github.com/GeertJohan/go.rice/rice

.PHONY: test
test: lint
	go test ./...

.PHONY: lint
lint:
	gometalinter ./...
