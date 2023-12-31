BROKER_BINARY=brokerApp
TABELOGSPIDER_BINARY=tabelogspiderApp
AUTHENTICATE_BINARY=authenticateApp
GOOGLEMAP_BINARY=googleMapApp
LOGGER_BINARY=loggerServiceApp
MAIL_BINARY=mailServiceApp
LISTENER_BINARY=listenerApp
FRONT_BINARY=frontEndApp
BROKER_BINARY_X86=brokerAppX86
TABELOGSPIDER_BINARY_X86=tabelogspiderAppX86
AUTHENTICATE_BINARY_X86=authenticateAppX86
GOOGLEMAP_BINARY_X86=googleMapAppX86
LOGGER_BINARY_X86=loggerServiceAppX86
MAIL_BINARY_X86=mailServiceAppX86
LISTENER_BINARY_X86=listenerAppX86
FRONT_BINARY_X86=frontEndAppX86

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
up_build: build_broker build_authenticate build_google-map build_tabelogspider build_logger build_mail build_listener
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

build_broker_x86:
	@echo "Building broker binary..."
	cd ../broker && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ${BROKER_BINARY_X86} ./cmd/api
	@echo "Done!"

## build_authenticate: build authenticate binary
build_authenticate:
	@echo "Building authenticate binary..."
	cd ../authenticate && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTHENTICATE_BINARY} ./cmd/api
	@echo "Done!"

build_authenticate_x86:
	@echo "Building authenticate binary..."
	cd ../authenticate && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ${AUTHENTICATE_BINARY_X86} ./cmd/api
	@echo "Done!"

## build_googlemap: build googlemap binary
build_google-map:
	@echo "Building google-map binary..."
	cd ../google-map && env GOOS=linux CGO_ENABLED=0 go build -o ${GOOGLEMAP_BINARY} ./cmd/api
	@echo "Done!"

build_google-map_x86:
	@echo "Building google-map binary..."
	cd ../google-map && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ${GOOGLEMAP_BINARY_X86} ./cmd/api
	@echo "Done!"

## build_tabelogspider: build tabelogspider binary
build_tabelogspider:
	@echo "Building tabelogspider binary..."
	cd ../tabelog-spider-service && env GOOS=linux CGO_ENABLED=0 go build -o ${TABELOGSPIDER_BINARY} ./cmd/api
	@echo "Done!"

build_tabelogspider_x86:
	@echo "Building tabelogspider binary..."
	cd ../tabelog-spider-service && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ${TABELOGSPIDER_BINARY_X86} ./cmd/api
	@echo "Done!"

## build logger: build logger binary
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

build_logger_x86:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ${LOGGER_BINARY_X86} ./cmd/api
	@echo "Done!"

## build mail: build mail binary
build_mail:
	@echo "Building mail binary..."
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Done!"

## build listener: build listener binary
build_listener:
	@echo "Building listener binary..."
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Done!"

build_listener_x86:
	@echo "Building listener binary..."
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ${LISTENER_BINARY_X86} .
	@echo "Done!"

build_front:
	@echo "Building front-end..."
	cd ../front-end && env GOOS=linux CGO_ENABLED=0 go build -o ${FRONT_BINARY} .
	@echo "Done!"

build_front_x86:
	@echo "Building front-end..."
	cd ../front-end && env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ${FRONT_BINARY_X86} .
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
create_sql:
	@echo "Creating sql from dbml..."
	cd ../db/dbml && dbml2sql --postgres -o schema.sql db.dbml
	@echo "Done!"

## create db
create_db:
	@echo "Creating db..."
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