package grpcserver

import (
	"google.golang.org/grpc"

	"github.com/andreysavvinov/petproject/internal/grpcjson"
)

func New() *grpc.Server {
	return grpc.NewServer(grpc.ForceServerCodec(grpcjson.Codec{}))
}
