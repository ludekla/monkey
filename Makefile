
vet:
	go vet ./...

fmt:
	go fmt ./...

test:
	go test ./pkg/...

clean:
	rm bin/*