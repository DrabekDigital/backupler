build:
	go build

build-osx-intel:
	GOOS=darwin GOARCH=amd64 go build -o bin/app-amd64-darwin .

build-osx-apple:
	GOOS=darwin GOARCH=arm64 go build -o bin/app-arm64-darwin .

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/app-amd64-linux .


tests:
	go test -v ./... # recursively

VER=`grep "version=" extras/spk/staging/INFO |cut -d'"' -f 2`

package-synology-amd64: build-linux
	rm -rf extras/spk/staging/
	mkdir -p extras/spk/staging/
	cp -r extras/spk/template/* extras/spk/staging
	mkdir extras/spk/staging/package/bin
	cp bin/app-amd64-linux extras/spk/staging/package/bin/backupler
	sh extras/spk/dir2spk extras/spk/staging
	mv staging-$(VER).spk backupler-$(VER).spk
	