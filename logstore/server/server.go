package server

import (
	"context"
	"logstore/internal/log/proto"
)

type CommitLog interface {
	Append(*proto.Record) (uint64, error)
	Read(uint64) (*proto.Record, error)
}

type Config struct {
	CommitLog CommitLog
}

type grpcServer struct {
	*Config
}

func newGrpcServer(config *Config) (srv *grpcServer, err error) {
	srv = &grpcServer{
		Config: config,
	}
	return srv, nil
}

func (s *grpcServer) Append(ctx context.Context, req *proto.AppendRequest) (
	*proto.AppendResponse, error) {
	off, err := s.CommitLog.Append(req.Record)
	if err != nil {
		return nil, err
	}
	return &proto.AppendResponse{Offset: off}, nil
}

func (s *grpcServer) Read(ctx context.Context, req *proto.ReadRequest) (
	*proto.ReadResponse, error) {
	record, err := s.CommitLog.Read(req.Offset)
	if err != nil {
		return nil, err
	}
	return &proto.ReadResponse{Record: record}, nil
}