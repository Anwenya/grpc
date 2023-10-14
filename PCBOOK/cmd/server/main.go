package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"com.wlq/pcbook/pb"
	"com.wlq/pcbook/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// 用于测试
func seedUsers(userStore service.UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}

	return createUser(userStore, "user1", "secret", "user")
}

func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(user)
}

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

// 访问权限映射
func accessibleRoles() map[string][]string {
	const laptopServicePath = "/pcbook.LaptopService/"
	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadImage":  {"admin"},
		laptopServicePath + "RateLaptop":   {"admin", "user"},
	}
}

// 加载服务端证书文件
func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// 加载为客户端签名的CA证书和私钥
	pemClientCA, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	// 加载服务端的证书和私钥
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	// 创建证书配置 双端认证
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		// 不认证客户端的证书
		// ClientAuth: tls.NoClientCert,
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  certPool,
	}

	return credentials.NewTLS(config), nil

}

func main() {
	port := flag.Int("port", 0, "the server port")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")

	flag.Parse()
	log.Printf("start server on port %d, TLS = %t", *port, *enableTLS)

	// authServer
	userStore := service.NewInMemoryUserStore()
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users")
	}
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)

	// laptopServer
	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskIamgeStore("./image")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)

	// 拦截器
	interceptor := service.NewAuthInterceptor(
		jwtManager,
		accessibleRoles(),
	)

	// 服务器配置项
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	}

	// 启用TLS
	if *enableTLS {
		// 加载证书
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			log.Fatal("cannot load TLS credentials: ", err)
		}

		// 添加tls配置项
		serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	}

	// 创建服务器 添加配置项
	grpcServer := grpc.NewServer(serverOptions...)

	// 注册服务
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	// 支持反射
	reflection.Register(grpcServer)

	// 开启监听
	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
