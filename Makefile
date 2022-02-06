dependencies:
	cd src/ && go get

build: dependencies
	export CGO_ENABLED=1
	go build -o build/manifold src/*.go

install:
	cp build/manifold /usr/bin/manifold

clear:
	rm -rf build/
