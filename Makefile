

CGO := 0
SOURCES := $(shell find . -type f -name "*.go")
.PHONE: build

all: build

build: pod-scraper

pod-scraper: $(SOURCES)
	CGO_ENABLED=$(CGO) go build -v -o pod-scraper ./cmd/pod-scraper
