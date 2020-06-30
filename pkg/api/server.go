package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"

	pb "github.com/whywaita/satelit/api/satelit"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/internal/logger"
	"github.com/whywaita/satelit/pkg/datastore"
	"github.com/whywaita/satelit/pkg/europa"
	"github.com/whywaita/satelit/pkg/ganymede"
	"github.com/whywaita/satelit/pkg/ipam"
)

// A SatelitServer is definition of Satlite API Server
type SatelitServer struct {
	pb.SatelitServer

	Datastore datastore.Datastore

	Europa   europa.Europa
	IPAM     ipam.IPAM
	Ganymede ganymede.Ganymede
}

// Run start gRPC Server
func (s *SatelitServer) Run() error {
	logger.Logger.Info(fmt.Sprintf("Run satelit server, listen on %s", config.GetValue().API.Listen))
	lis, err := net.Listen("tcp", config.GetValue().API.Listen)
	if err != nil {
		return err
	}
	opts := []grpc_zap.Option{
		grpc_zap.WithMessageProducer(grpc_zap.DefaultMessageProducer),
	}
	grpc_zap.ReplaceGrpcLoggerV2(logger.Logger)
	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.PayloadUnaryServerInterceptor(logger.Logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
				return true
			}),
			grpc_zap.UnaryServerInterceptor(logger.Logger, opts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.PayloadStreamServerInterceptor(logger.Logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
				return true
			}),
			grpc_zap.StreamServerInterceptor(logger.Logger, opts...),
		),
	)
	pb.RegisterSatelitServer(grpcServer, s)

	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

// parseRequestUUID return uuid.UUID from gRPC request string
func (s *SatelitServer) parseRequestUUID(reqUUID string) (uuid.UUID, error) {
	u := uuid.FromStringOrNil(reqUUID)
	if u == uuid.Nil {
		return uuid.Nil, fmt.Errorf("failed to parse uuid from string (uuid: %s)", reqUUID)
	}

	return u, nil
}

var pool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 1024*64))
	},
}

type meta struct {
	name        string
	description string
}

// header.Size is byte
func sanitizeImageSize(headerSize uint64) int {
	const (
		BYTE float64 = 1 << (10 * iota)
		KILOBYTE
		MEGABYTE
		GIGABYTE
	)

	sizeGB := float64(headerSize) / GIGABYTE
	gb := math.Trunc(sizeGB)

	if sizeGB == gb {
		// not digit loss
		return int(gb)
	}

	// if occurred digit loss, disk capacity need to extend + 1GB
	return int(gb) + 1
}

func (s *SatelitServer) receiveImage(stream pb.Satelit_UploadImageServer, w io.Writer) (meta, error) {
	m := meta{}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return meta{}, fmt.Errorf("failed to recv file: %w", err)
		}

		if mt := resp.GetMeta(); mt != nil {
			m.name = mt.Name
			m.description = mt.Description
		}
		if chunk := resp.GetChunk(); chunk != nil {
			_, err := w.Write(chunk.Data)
			if err != nil {
				return meta{}, fmt.Errorf("failed to write chunk data: %w", err)
			}
		}
	}

	return m, nil
}
