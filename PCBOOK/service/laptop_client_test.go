package service_test

import (
	"context"
	"net"
	"testing"

	"com.wlq/pcbook/pb"
	"com.wlq/pcbook/sample"
	"com.wlq/pcbook/serializer"
	"com.wlq/pcbook/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopStore := service.NewInMemoryLaptopStore()
	serverAddress := startTestLaptopServer(t, laptopStore, nil, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedId := laptop.Id
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedId, res.Id)

	other, err := laptopStore.Find(res.Id)
	require.NoError(t, err)
	require.NotNil(t, other)

	requireSampleLaptop(t, laptop, other)

}

// 序列化对象不能直接比较,因为其中有protobuf生成的特字段,先转换成JSON再比较
func requireSampleLaptop(t *testing.T, laptop *pb.Laptop, other *pb.Laptop) {
	json1, err := serializer.ProtobufToJSON(laptop)
	require.NoError(t, err)
	json2, err := serializer.ProtobufToJSON(laptop)
	require.NoError(t, err)

	require.Equal(t, json1, json2)
}

// 服务端
func startTestLaptopServer(t *testing.T,
	laptopStore service.LaptopStore,
	imageStore service.ImageStore,
	ratingStore service.RatingStore,
) string {
	laptopServer := service.NewLaptopServer(
		laptopStore,
		imageStore,
		ratingStore,
	)

	// 注册服务
	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	// 随机端口
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	// 监听
	go grpcServer.Serve(listener)

	return listener.Addr().String()
}

// 客户端
func newTestLaptopClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}
