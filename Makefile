NAME:=vc-payout
DC=docker-compose -f ./docker/docker-compose.yaml --env-file ./.env
tidy:
	rm -f go.sum; go mod tidy

vet:
	go vet ./...

stop:
	$(DC) stop

test:
	go test ./... --count=1

clean:
	$(DC) down --rmi local --remove-orphans -v
	$(DC) rm -f -v

build: stop
	$(DC) build vc-payout

migrate: db-up
	sleep 5;
	go run ./cmd/migrate/main.go

db-down:
	$(DC) down postgres

db-up: 
	$(DC) up -d postgres

swag:
	swag init -g ./cmd/restapi/main.go -o ./docs

run: stop
	$(DC) up -d postgres
	sleep 3;
	$(DC) up vc-payout

dev: migrate  
	go run ./cmd/restapi/main.go
