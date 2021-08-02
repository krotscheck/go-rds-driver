
.PHONY: test vet lint sec clean

clean:
	rm -rf ./reports
	rm ./rdsdataservice_mocks_test.go

reports:
	mkdir -p ./reports

rdsdataservice_mocks_test.go:
	go generate ./...

test: rdsdataservice_mocks_test.go reports
	gotestsum --junitfile reports/unit-tests.xml -- -p 1 -covermode=atomic -coverpkg=./... -coverprofile=reports/coverage.out ./...
	go tool cover -func reports/coverage.out

vet: reports
	go vet ./...

lint: reports
	golint ./...

sec: reports
	gosec ./...

checks: test vet lint sec
