gen-proto:
	@name=$(name);
	@protoc --proto_path=common/api/$(name) \
	  --go_out=common/api/$(name) --go_opt=paths=source_relative \
		--go-grpc_out=common/api/$(name) --go-grpc_opt=paths=source_relative \
		common/api/$(name)/$(name).proto

dev:
	@docker compose -f compose.dev.yml up -d

stop-dev:
	@docker compose -f compose.dev.yml down

start:
	@docker compose -f compose.yml up -d

stop:
	@docker compose -f compose.yml down