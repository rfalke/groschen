all: fmt test build run

build:
	go build -o groschen.exe main/groschen.go

fmt:
	go fmt

test:
	go test -v

run:
	./groschen.exe -o _out http://www.google.com

clean:
	rm -f groschen.exe *~
	rm -rf _out

