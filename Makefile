export GO111MODULE=on

default: test

ci: depsdev test

test:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic

lint:
	golangci-lint run ./...

depsdev:
	go install github.com/Songmu/ghch/cmd/ghch@latest

credits:
	go install github.com/Songmu/gocredits/cmd/gocredits@v1.0.0
	gocredits . > CREDITS

prerelease:
	git pull origin main --tag
	go mod tidy
	ghch -w -N ${VER}
	$(MAKE) credits
	git add CHANGELOG.md CREDITS go.mod go.sum
	git commit -m'Bump up version number'
	git tag ${VER}

prerelease_for_tagpr:
	$(MAKE) credits
	git add CHANGELOG.md CREDITS go.mod go.sum

release:
	git push origin main --tag

.PHONY: default test credits
