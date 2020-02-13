package main

import (
	"net"

	log "github.com/sirupsen/logrus"

	//"gitlab.com/sdce/service/currency/pkg/config"

	"gitlab.com/sdce/service/currency/pkg/rpc"

	pb "gitlab.com/sdce/protogo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	mgo "gitlab.com/sdce/exlib/mongo"

	"gitlab.com/sdce/service/currency/pkg/repository"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

//Service ...
type Service struct {
	rpc       *rpc.Server
	healthSvr *health.Server
}

//NewService ...
func NewService(db *mgo.Database) *Service {

	repo, err := repository.CreateRepository(db)
	if err != nil {
		log.Fatalln(err.Error())
	}
	s := &Service{
		// Session: sess,
		rpc:       rpc.NewRPCServer(nil, ":8027", repo),
		healthSvr: health.NewServer(),
	}
	s.healthSvr.SetServingStatus("grpc.health.v1.currency", grpc_health_v1.HealthCheckResponse_SERVING)
	return s
}

//Run ...
func (s *Service) Run() {
	lis, err := net.Listen("tcp", ":8029")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("starting rpc server, listening on %s", lis.Addr())
	gs := grpc.NewServer()
	if s.rpc == nil {
		log.Fatalln("Please create a service first, OK?")
		return
	}
	pb.RegisterCurrencyServiceServer(gs, s.rpc)
	grpc_health_v1.RegisterHealthServer(gs, s.healthSvr)
	// Register reflection service on gRPC server.
	reflection.Register(gs)
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
