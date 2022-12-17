build:
	cd cmd/queue && go build

test: build
	go test -race -cover

run:
	@cd cmd/queue && ./queue -m MM1 -a 9 -s 10
