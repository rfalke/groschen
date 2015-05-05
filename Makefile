all: ensure-deps fmt test build run

ensure-deps: dependencies/src/gopkg.in/alecthomas/kingpin.v1/app.go

dependencies/src/gopkg.in/alecthomas/kingpin.v1/app.go:
	mkdir -p dependencies
	export GOPATH=$$(pwd)/dependencies ; cd dependencies && go get gopkg.in/alecthomas/kingpin.v1
	
build:
	export GOPATH=$$(pwd)/dependencies:$$(pwd) ; cd src/groschen ; go build -o groschen.exe main/groschen.go

fmt:
	export GOPATH=$$(pwd)/dependencies:$$(pwd) ; cd src/groschen ; go fmt

test:
	export GOPATH=$$(pwd)/dependencies:$$(pwd) ; cd src/groschen ; go test -v

run:
	./src/groschen/groschen.exe --r3 -o _out http://www.google.com

clean:
	rm -f src/groschen/groschen.exe src/groschen/*~
	rm -rf _out

mrproper: clean
	rm -rf dependencies
