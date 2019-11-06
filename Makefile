deps:
	go mod vendor

build: deps
	go build -o webcrawler

run: build
	./webcrawler crawl

long_run: build
	./webcrawler crawl -d https://www.monzo.com -w 20 -l debug

test:
	go test -v ./...