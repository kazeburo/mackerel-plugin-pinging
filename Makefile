VERSION=0.0.4
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

all: mackerel-plugin-pinging

.PHONY: mackerel-plugin-pinging

bundle:
	dep ensure

update:
	dep ensure -update

mackerel-plugin-pinging: main.go
	go build $(LDFLAGS) -o mackerel-plugin-pinging

linux: main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-pinging

check:
	go test ./...

fmt:
	go fmt ./...

tag:
	git tag v${VERSION}
	git push origin v${VERSION}
	git push origin master
	goreleaser --rm-dist
