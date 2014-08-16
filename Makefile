.PHONY: all

GO=go
GOOS=linux
GOARCH=arm
GOARM=5

all:
	${GO} build

arm:
	GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} ${GO} build

clean:
	rm -f ${GOPATH}/pkg/*/${PACKAGE}.a

fmt:
	go fmt ./...

watch:
	watchmedo shell-command --patterns="*.go" --recursive --wait \
          --command="make arm && echo \"...complete\""
