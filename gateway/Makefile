run:
	@go run cmd/main.go --config=../config/gateway.yml

gen-docs:
	@swag init -g cmd/main.go --parseDependency --parseInternal
