
# protobuf 文件定义
# syntax = "proto3";
# option go_package = "path;name";
# path 表示生成的go文件的存放地址,会自动生成目录的。
# name 表示生成的go文件所属的包名
# option go_package="./;proto";
# 定义包名
# package proto;

# -I:import目录,.proto文件中的import命令将从该目录下开始找,默认值是当前执行编译命令的目录.
# --go_out:生成的protobuf的GO文件存放路径
# --go-grpc_out:冒号前边是配置项,后边是生成的gRPC的GO文件存放路径
# 最后一个参数是要编译的.proto文件的路径及通配符
gen-protobuf-grpc:
	protoc -I=proto --go_out=./pb --go-grpc_out=require_unimplemented_servers=false:./pb proto/*.proto
	
# 启动服务端
start-server:
	go run ./cmd/server/main.go -port=50051
# 启动客户端
start-client:
	go run ./cmd/client/main.go -address=0.0.0.0:50051
# 启动服务端并开启TLS
start-server-tls:
	go run ./cmd/server/main.go -port=50051 -tls
# 启动客户端并开启TLS
start-client-tls:
	go run ./cmd/client/main.go -address=0.0.0.0:50051 -tls
start-server-1:
	go run ./cmd/server/main.go -port=50051
start-server-2:
	go run ./cmd/server/main.go -port=50052
start-server-tls-1:
	go run ./cmd/server/main.go -port=50051 -tls
start-server-tls-2:
	go run ./cmd/server/main.go -port=50052 -tls

start-client-nginx:
	go run ./cmd/client/main.go -address=0.0.0.0:8080

start-client-nginx-tls:
	go run ./cmd/client/main.go -address=0.0.0.0:8080 -tls

# .PHONY: gen-protobuf-grpc start-server start-client start-server-tls start-client-tls start_server1 start_server-2 start_client_nginx

# evans -r -p 8080



# openssl req -newkey rsa:4096 -keyout client-key.pem -out client-req.pem -nodes -subj "/C=CN/ST=WLQNET/L=grpc-client/O=grpc-client/OU=grpc-client/CN=grpc-client/emailAddress=grpc-client@gmail.com"

# openssl x509 -req -in client-req.pem -days 90 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.cnf