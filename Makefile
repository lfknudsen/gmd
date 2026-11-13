all: build
	install
	

build: cmd/gmd/main.go pkg/lib.go
	go build -C cmd/gmd -o ../..
