BROKER_BINARY=brokerApp
TABELOGSPIDER_BINARY=tabelogspiderApp

## up: start docker compose
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## up_build: start docker compose with build
up_build: build_tabelogspider build_broker 
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## build_broker: build broker binary
build_broker:
	@echo "Building broker binary..."
	cd ../broker && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_tabelogspider: build tabelogspider binary
build_tabelogspider:
	@echo "Building tabelogspider binary..."
	cd ../tabelog-spider-service && env GOOS=linux CGO_ENABLED=0 go build -o ${TABELOGSPIDER_BINARY} ./cmd/api
	@echo "Done!"

## start front-end
start_front:
	@echo "Starting front-end..."
	cd ../front-end && go run main.go
	@echo "Done!"

## push dbml
push_dbml:
	@echo "Pushing dbml..."
	cd ../db/dbml && dbdocs build db.dbml
	@echo "Done!"

## build sql from dbml
build_sql:
	@echo "Building sql from dbml..."
	cd ../db/dbml && dbml2sql --postgres -o schema.sql db.dbml
	@echo "Done!"

## build db
build_db:
	@echo "Building db..."
	docker exec -it project-postgres-1 createdb --username=postgres tabelogo
	@echo "Done!"

## drop db
drop_db:
	@echo "Dropping db..."
	docker exec -it project-postgres-1 dropdb --username=postgres tabelogo
	@echo "Done!"

## migrate up
migrate_up:
	@echo "Migrating up..."
	cd ../authenticate/cmd/data/migration && migrate -path . -database "postgresql://postgres:password@localhost:5432/tabelogo?sslmode=disable" up
	@echo "Done!"

## migrate down
migrate_down:
	@echo "Migrating down..."
	cd ../authenticate/cmd/data//migration && migrate -path . -database "postgresql://postgres:password@localhost:5432/tabelogo?sslmode=disable" down
	@echo "Done!"

## sqlc generate
sqlc:
	@echo "Generating sqlc..."
	cd ../authenticate/cmd/data && sqlc generate
	@echo "Done!"