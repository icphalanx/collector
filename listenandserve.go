package collector

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"

	pb "github.com/icphalanx/rpc"
)

func ListenAndServeAtAddr(addr string, c Config) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	col, err := NewServer(c)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.Creds(credentials.NewTLS(c.TLSConfig)))
	log.Println("collector serving at", addr)
	pb.RegisterPhalanxCollectorServer(grpcServer, col)
	return grpcServer.Serve(lis)
}
