run:
	@export $$(cat .env | xargs) && go run cmd/moscode-cdn-fiber/main.go

build:
	@export $$(cat .env | xargs) && go build cmd/moscode-cdn-fiber/main.go
