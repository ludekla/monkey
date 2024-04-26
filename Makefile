
vet:
	go vet ./...

fmt:
	go fmt ./...

test:
	go test ./...

vtest:
	go test -v ./...

clean:
	rm bin/*