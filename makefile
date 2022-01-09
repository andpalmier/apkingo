all: apkingo

apkingo: build
	@go build -o build/apkingo cmd/*.go

install:
	@go install ./cmd

build:
	@mkdir -p build

clean:
	@rm -rf build
