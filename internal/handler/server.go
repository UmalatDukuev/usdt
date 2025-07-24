package handler

import (
	"context"
	pb "usdt/internal/handler/pb"

	"time"
	"usdt/internal/service"
)

type Server struct {
	pb.UnimplementedRateServiceServer
	Service service.RateServiceInterface
}

// GetRates calls service to get current rates and stores them, then returns the response protobuf.
func (s *Server) GetRates(ctx context.Context, _ *pb.Empty) (*pb.RateResponse, error) {
	rate, err := s.Service.GetAndStoreRates(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.RateResponse{
		Ask:       rate.Ask,
		Bid:       rate.Bid,
		Timestamp: rate.Timestamp.Format(time.RFC3339),
	}, nil
}

// HealthCheck returns a simple status indicating service health.
func (s *Server) HealthCheck(ctx context.Context, _ *pb.Empty) (*pb.HealthStatus, error) {
	return &pb.HealthStatus{Status: "OK"}, nil
}
