all: apkingo

apkingo: build
	@go build -o build/apkingo cmd/*.go

install: apkingo
	@cp build/apkingo /usr/bin/
	@chmod a+x /usr/bin/apkingo

build:
	@mkdir -p build

clean:
	@rm -rf build
