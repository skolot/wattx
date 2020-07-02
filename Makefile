.PHONY: start stop

start:
	docker-compose up --build --detach

stop:
	docker-compose down