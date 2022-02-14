
.PHONY: fmt test vet lint sec clean

fmt:
	go fmt ./...

clean:
	rm -rf ./reports
	rm ./rdsdataservice_mocks_test.go

reports:
	mkdir -p ./reports

rdsdataservice_mocks_test.go:
	go generate ./...

test_deps:
	go install gotest.tools/gotestsum@v1.7.0
	go install github.com/golang/mock/mockgen@v1.6.0

test: test_deps reports/coverage.xml reports/html/index.html

vet: reports
	go vet ./...

lint: reports
	golint ./...

sec: reports
	gosec ./...

reports/coverage.out: reports rdsdataservice_mocks_test.go
	gotestsum --junitfile reports/unit-tests.xml -- -p 1 -covermode=atomic -coverpkg=./... -coverprofile=reports/coverage.out ./...
	go tool cover -func reports/coverage.out

reports/html/index.html: reports/coverage.out
	mkdir -p ./reports/html
	gocov convert ./reports/coverage.out | gocov-html > ./reports/html/index.html

reports/coverage.xml: reports/coverage.out
	gocov convert ./reports/coverage.out | gocov-xml > ./reports/coverage.xml

checks: test vet lint sec
