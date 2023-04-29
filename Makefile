GOCMD=go
BINARY_NAME=gaze
OUT=dist
CMD=cmd/gaze/main.go

build:
	go build -ldflags "-s -w" -v ${CMD}
build-all: build-macos-amd build-macos-arm build-windows build-linux
build-all-amd: build-macos-amd build-windows build-linux
build-macos-amd:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ${OUT}/macos_amd/$(BINARY_NAME) -v ${CMD}
build-macos-arm:
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o ${OUT}/macos_arm/$(BINARY_NAME) -v ${CMD}
build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ${OUT}/windows/$(BINARY_NAME).exe -v ${CMD}
build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${OUT}/linux/$(BINARY_NAME) -v ${CMD}
ut:
	go test github.com/wtetsu/gaze/pkg/...
e2e:
	go build -ldflags "-s -w" -o test/e2e -v ${CMD}
	cd test/e2e && sh test_all.sh
	
cov:
	go test -coverprofile=coverage.txt -covermode=atomic github.com/wtetsu/gaze/pkg/...
clean:
	go clean ${CMD}
	rm -f ${OUT}/macos-amd/$(BINARY_NAME)
	rm -f ${OUT}/macos-arm/$(BINARY_NAME)
	rm -f ${OUT}/windows/$(BINARY_NAME)
	rm -f ${OUT}/linux/$(BINARY_NAME)
