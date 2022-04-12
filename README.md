# calc-grpc-me
# just a example from https://dev.to/greenteabiscuit/mini-grpc-project-creating-a-simple-increment-api-on-go-6cn
# just cannot reproduct -devel suffix for protoc-gen-go...

go run src/backend/main.go && go run src/frontend/main.go && curl "http://localhost:8080/increment?val=10"

protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:.  --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative --openapiv2_out=.  -I . -I /home/workspace/shiyuelc/proto/include calc.proto