run:
	@export $$(cat .env | xargs) && go run cmd/moscode-cdn-fiber/main.go

build:
	@export $$(cat .env | xargs) && go build cmd/moscode-cdn-fiber/main.go

docker:
	docker build -t smbrine/moscode-cdn .
	docker push smbrine/moscode-cdn