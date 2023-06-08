init:
	export PATH=$PATH:$(go env GOPATH)/bin
	# go install golang.org/x/lint/golint
	# go install golang.org/x/tools/cmd/godoc

run:
	go run main.go

fmt:
	gofmt -w .

lint:
	go vet
	golint ./...

generate-wallet:
	go run cmd/bot/main.go generate-wallet