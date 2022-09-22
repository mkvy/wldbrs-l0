start-docker:
		docker-compose up

stop-docker:
		docker-compose stop

test:
		cd server-subscriber && go test ./cache && go test ./database

run-subscriber:
		cd server-subscriber && go build ./
		cd server-subscriber && server-subscriber.exe

run-publisher:
		cd publisher && go build ./
		cd publisher && publisher.exe

.PHONY: start-docker stop-docker test run-publisher run-subscriber