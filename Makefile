PROJECTNAME="Executor"

all: clean deps linux

clean:
	echo "clean all"
	rm -rf dist/

deps:
	go mod vendor

linux:
	mkdir -p dist/linux/

	go build -o dist/linux/skywalking_executor cmd/main.go


