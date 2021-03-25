.PHONY: build

# makes sure of :https://github.com/cosmtrek/air/ for go livereload.
air: link
	@if command -v "air" ; then "air"  ; else echo \
		">>> Air not installed\n use: https://github.com/cosmtrek/air/" ;fi

build:
	@echo ">>> build binary"
	@go build -o turl.to cmd/turl.to/main.go
	@chmod +x turl.to

build-docker: build
	@docker build -t turl.to:latest .

clean:
	@if [ -x "turl.to" ]; then rm "turl.to"; fi
	@if [ -d "tmp" ]; then rm -rf "tmp"; fi
	@if [ -f "main.go" ]; then test -L "main.go" && rm "main.go"; fi

deps:
	@go get -v -t -d ./...
	@go mod vendor

link: clean
	@ln -s cmd/turl.to/main.go main.go

test:
	@go test -v --race ./...
