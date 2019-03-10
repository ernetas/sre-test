PHONY: build

build:
	docker build -t ernestas/sre-test:$(shell git rev-parse HEAD) code/
	docker push ernestas/sre-test:$(shell git rev-parse HEAD)
