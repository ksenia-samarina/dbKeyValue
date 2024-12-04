package grpc

import (
	"context"
	protos "github.com/ksenia-samarina/protoFiles/gen/go/lsmt"
	"google.golang.org/grpc"
)

type serverAPI struct {
	protos.UnimplementedLSMTServer
	lsmt LSMT
}

type LSMT interface {
	Get(ctx context.Context, key string) ([]byte, bool)
	Put(ctx context.Context, key string, value []byte) bool
	Delete(ctx context.Context, key string) bool
}

func (s *serverAPI) Get(ctx context.Context, in *protos.GetReq) (*protos.GetRes, error) {
	// TODO: validation, errors
	var err error
	b, ok := s.lsmt.Get(ctx, in.Key)
	if !ok {
		return nil, err
	}
	return &protos.GetRes{Value: b}, nil
}

func (s *serverAPI) Put(ctx context.Context, in *protos.PutReq) (*protos.PutRes, error) {
	// TODO: validation, errors
	var err error
	ok := s.lsmt.Put(ctx, in.Key, in.Value)
	if !ok {
		return &protos.PutRes{Ok: false}, err
	}
	return &protos.PutRes{Ok: true}, nil
}

func (s *serverAPI) Delete(ctx context.Context, in *protos.DeleteReq) (*protos.DeleteRes, error) {
	// TODO: validation, errors
	var err error
	ok := s.lsmt.Delete(ctx, in.Key)
	if !ok {
		return &protos.DeleteRes{Ok: false}, err
	}
	return &protos.DeleteRes{Ok: true}, nil
}

func Register(grpcServer *grpc.Server, lsmt LSMT) {
	protos.RegisterLSMTServer(grpcServer, &serverAPI{lsmt: lsmt})
}
