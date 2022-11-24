mkdir -p ./gen/proto
rm ./gen/proto/*

protoc --proto_path=./protos --go_out=./gen/proto --go_opt=paths=source_relative \
    --go-grpc_out=./gen/proto --go-grpc_opt=paths=source_relative \
    protos/main.proto