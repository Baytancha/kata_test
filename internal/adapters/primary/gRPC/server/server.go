package server

import (
	pb "kata_test/generated/protobuf/rates"
	rates "kata_test/internal/adapters/primary/gRPC/handlers/rates"
	xlogger "kata_test/internal/infrastructure/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	//healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	sleep  = time.Duration(time.Second * 5)
	system = ""
)

type GRPConfig struct {
	Addr string
	Port int
	Flag string
}

type GrpcServer struct {
	addr         string
	port         int
	flag         string
	sleep        time.Duration
	server       *grpc.Server
	healthserver *health.Server
	logger       *zap.Logger
	rates        *rates.RatesHandler
}

type Handlers struct {
	Rates *rates.RatesHandler
}

func New(cfg *GRPConfig, handlers *Handlers) *GrpcServer {
	return &GrpcServer{
		addr:   cfg.Addr,
		port:   cfg.Port,
		flag:   cfg.Flag,
		sleep:  sleep,
		logger: xlogger.Logger().Named("grpc"),
		rates:  handlers.Rates,
	}
}

func (g *GrpcServer) Shutdown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		s := <-sigCh
		g.logger.Info("got signal ,attempting graceful shutdown", zap.String("signal", s.String()))
		g.server.GracefulStop()
		wg.Done()
	}()
}

func (g *GrpcServer) Healthcheck() {
	go func() {
		next := healthpb.HealthCheckResponse_SERVING

		for {
			g.healthserver.SetServingStatus(system, next)

			if next == healthpb.HealthCheckResponse_NOT_SERVING {
				g.logger.Info("server is not serving")
			}

			time.Sleep(g.sleep)
		}
	}()
}

func (g *GrpcServer) SetupAndServe() error {

	listen, err := net.Listen("tcp", g.addr)
	if err != nil {
		g.logger.Fatal("failed to listen on port", zap.Int("port", g.port), zap.Error(err))
	}

	GrpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	healthcheck := health.NewServer()
	g.healthserver = healthcheck
	healthpb.RegisterHealthServer(GrpcServer, g.healthserver)

	g.server = GrpcServer
	pb.RegisterRatesServiceServer(g.server, g.rates)

	if g.flag == "development" {
		reflection.Register(GrpcServer)
	}

	go g.Healthcheck()

	g.logger.Info("starting api service on port", zap.Int("port", g.port))

	return GrpcServer.Serve(listen)

	// if err := GrpcServer.Serve(listen); err != nil {
	// 	g.logger.Fatal("failed to serve grpc on port", zap.Int("port", g.port), zap.Error(err))
	// }
}

func (g *GrpcServer) Stop() {
	g.server.Stop()
}
