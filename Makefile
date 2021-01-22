help:           ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

example: build  ## Run challenge example.
	./example.sh

build: test     ## Build /bin/matchseq binary.
	go build -o bin/matchseq cmd/matchseq/main.go

test:           ## Run tests.
	go test -coverprofile /tmp/rpm-test-coverage ./...
