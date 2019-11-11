# Concurrent Webcrawler

## Running the Program

build using `make build`

Flags
1. `-d` followed by a string representing the domain URL to crawl
1. `-w` followed by the number of concurrent workers
1. `-l` followed by "DEBUG" or "INFO" to set the log level.
1. `-h` shows the help menu

To run tests:
`make test`

## The Flow

1. Setup goroutine responsible for publishing urls into a channel for workers to consume from.
    - Maintains a counter that indicates how many pending links need to be processed. A Breadth First Search analogy would be `while queue is not empty`
1. Create a worker pool that consumes from a string channel
    - If a page has been processed, it is ignored.
    - Workers handle url by making a GET request and then parsing the HTML to obtain the `Title` of the page and the values in the `href` attribute of `<a></a>` tags.
    - Workers then add the `Page` data into a sitegraph, a map data structure that maps a string URL -> `Page` struct, that has the `Title` and `Links` information.
    - Any links that were parsed from the HTML is sent back to the URL publisher.
1. When there are no more pending links to crawl, close the publishing channel to gracefully shut down workers in the worker pool.

## Test Coverage and Strategy
Because of the highly dynamic nature of webpages, it is easier and more predictable to test the functionality of the crawler with HTML fixtures.
The html files can be found in the `internal/crawler/test/html` folder.
An http testserver is created to serve the web pages when the webcrawler makes GET requests to the domain when run using `go test ./...`

1. 100% `graph.go`
1. 94.1% `parse.go`
1. 81.8% `crawler.go`

Happy path has test coverage.

## Performance
* With 1 worker, time to crawl www.monzo.com -> ~15m
* With 20 workers, time to crawl www.monzo.com -> 1m 40s
