PHONY: build test

build:
	docker build -t ernestas/sre-test:$(shell git rev-parse HEAD) code/
	docker push ernestas/sre-test:$(shell git rev-parse HEAD)

test:
	cd code && go get -v -t -d && go test -v
