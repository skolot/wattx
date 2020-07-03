.PHONY: start stop test-start test-stop test run-test

start:
	docker-compose up --build --detach

stop:
	docker-compose down

test-start:
	COMPOSE_PROJECT_NAME=wattx_test docker-compose -f docker-compose-test.yml up --build --detach

test-stop:
	COMPOSE_PROJECT_NAME=wattx_test docker-compose -f docker-compose-test.yml down

test: test-start run-test test-stop

run-test:
	cd ./test && \
		go test -v -count=1 main_test.go
