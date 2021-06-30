

SOURCES := $(shell find . -type f -name "*.go")
.PHONE: build

all: build

build: pod-scraper

pod-scraper: $(SOURCES)
	go build -o pod-scraper ./cmd/pod-scraper
