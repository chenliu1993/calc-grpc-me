# calc-grpc-me
# just a example from https://dev.to/greenteabiscuit/mini-grpc-project-creating-a-simple-increment-api-on-go-6cn
# just cannot reproduct -devel suffix for protoc-gen-go...

go run src/backend/main.go && go run src/frontend/main.go && curl "http://localhost:8080/increment?val=10"