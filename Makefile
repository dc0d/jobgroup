.PHONY: test
test:
	clear
	go test -p 1 -count=1 -timeout 30s -coverprofile=./cover-profile.out -covermode=set -coverpkg=./... ./...; \
	go tool cover -html=./cover-profile.out -o ./cover-report.html

lint:
	golangci-lint run ./...
