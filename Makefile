
.PHONY: test vet lint sec clean

clean:
	rm -rf ./reports
	rm ./rdsdataservice_mocks_test.go

reports:
	mkdir -p ./reports

rdsdataservice_mocks_test.go:
	go generate ./...

test: reports/coverage.xml reports/html/index.html

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
