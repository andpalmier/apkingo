NAME=apkingo
SOURCE=cmd/*.go
BUILD_FOLDER=build

all: apkingo

apkingo: build
	@go build -o $(BUILD_FOLDER)/$(NAME) $(SOURCE)

build:
	@mkdir -p $(BUILD_FOLDER)

clean:
	@rm -rf $(BUILD_FOLDER)
