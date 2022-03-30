include ./config/.envrc

# --------------------------------------------------------------------------- #
# HELPERS
# --------------------------------------------------------------------------- #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

# --------------------------------------------------------------------------- #
# DEVELOPMENT
# --------------------------------------------------------------------------- #

## run: run the ./cmd/urlshortener application
.PHONY:
run:
	go run ./cmd/urlshortener -db=${URLSHORTENER_DB_DSN}

## db/migrations/up: migrate up all the migration files in ./migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${URLSHORTENER_DB_DSN} up

# --------------------------------------------------------------------------- #
# QUALITY CONTROL
# --------------------------------------------------------------------------- #

## audit: tidy and dependencies, format, vet, test
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	@echo 'Running test...'
	go test -cover ./...


## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring all dependencies...'
	go mod vendor

# --------------------------------------------------------------------------- #
# BUILD
# --------------------------------------------------------------------------- #

## build: build the ./cmd/urlshortener application
.PHONY: build/urlshortener
build:
	@echo 'Build ./cmd/urlshortener...'
	go build -o=./bin/urlshortener ./cmd/urlshortener
	
