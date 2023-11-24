.PHONY: proto run

run:
	go run cmd/app/main.go 

p1:
	protoc --proto_path=api/service/ --go_out=./internal/services/pb \
  	--go-grpc_out=./internal/services/pb --go-grpc_opt=paths=source_relative \
  	api/service/*.proto

proto:
	protoc --proto_path=api/service/ --go_out=./internal/services/pb --go_opt=paths=source_relative \
	--go-grpc_out=./internal/services/pb --go-grpc_opt=paths=source_relative api/service/*.proto

gen:
	protoc -I=api/service/ --go_out=./internal/services/pb api/service/profile.proto