gen-proto:
	@name=$(name);
	@protoc --proto_path=common/api/$(name) \
	  --go_out=common/api/$(name) --go_opt=paths=source_relative \
		--go-grpc_out=common/api/$(name) --go-grpc_opt=paths=source_relative \
		common/api/$(name)/$(name).proto

run-sso:
	@go run sso/cmd/main.go