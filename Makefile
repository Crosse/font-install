package = github.com/Crosse/font-install

.PHONY: release

release:
	mkdir -p release
	GOOS=linux   GOARCH=amd64 go build -o release/font-install-linux-amd64   $(package)
	GOOS=linux   GOARCH=386   go build -o release/font-install-linux-386     $(package)
	GOOS=linux   GOARCH=arm   go build -o release/font-install-linux-arm     $(package)
	GOOS=linux   GOARCH=arm64 go build -o release/font-install-linux-arm64   $(package)
	GOOS=darwin  GOARCH=amd64 go build -o release/font-install-darwin-amd64  $(package)
	GOOS=windows go get .
	GOOS=windows GOARCH=amd64 go build -o release/font-install-windows-amd64 $(package)
	GOOS=windows GOARCH=386   go build -o release/font-install-windows-386   $(package)