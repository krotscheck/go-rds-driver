
generate:
	go generate ./...

test: generate
	go test ./...

vet:
	go vet ./...

lint:
	golint ./...

sec:
	gosec ./...

checks: test vet lint sec
