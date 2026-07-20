package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/aclgo/grpc-mail/config"
	grpcService "github.com/aclgo/grpc-mail/internal/mail/delivery/grpc/service"
	httpService "github.com/aclgo/grpc-mail/internal/mail/delivery/http/service"
	"github.com/aclgo/grpc-mail/internal/server/interceptors"
	grpcauth "github.com/aclgo/grpc-mail/pkg/grpc_auth"
	"github.com/aclgo/grpc-mail/pkg/logger"
	"github.com/aclgo/grpc-mail/proto"

	"google.golang.org/grpc"
)

type Server struct {
	config       *config.Config
	logger       logger.Logger
	serviceHTTP  *HttpHandlerService
	servicesGRPC *grpcService.MailService
	stopFn       sync.Once
	grpcAuth     *grpcauth.GrpcAuth
}

type HttpHandlerService struct {
	endpoint string
	service  *httpService.MailService
}

func NewHttpHandlerService(
	endpoint string,
	service *httpService.MailService,
) *HttpHandlerService {
	return &HttpHandlerService{
		endpoint: endpoint,
		service:  service,
	}
}

func NewServer(cfg *config.Config,
	logger logger.Logger,
	svcHTTP *HttpHandlerService,
	svcsGRPC *grpcService.MailService,
	grpcAuth *grpcauth.GrpcAuth) *Server {
	return &Server{
		config:       cfg,
		logger:       logger,
		serviceHTTP:  svcHTTP,
		servicesGRPC: svcsGRPC,
		grpcAuth:     grpcAuth,
	}
}

func (s *Server) Run(ctxSignal context.Context) error {

	ctxHttp := context.Background()

	var (
		errHTTP = make(chan error)
		errGRPC = make(chan error)
	)

	go func() {
		// s.logger.Infof("http server init port %v", s.config.ServiceHTTPPort)
		err := s.httpRun(ctxHttp)
		if err != nil {
			s.logger.Errorf("Run:%v", err)
			errHTTP <- fmt.Errorf("Run:%v", err)
		}
	}()

	go func() {
		// s.logger.Infof("grpc server init port %v", s.config.ServiceGRPCPort)
		err := s.grpcRun()
		if err != nil {
			s.logger.Errorf("Run:%v", err)
			errGRPC <- fmt.Errorf("Run:%v", err)
		}
	}()

	select {
	case eHTTP := <-errHTTP:
		return eHTTP
	case eGRPC := <-errGRPC:
		return eGRPC
	case <-ctxSignal.Done():
		s.logger.Info("shutting down servers")
		return nil
	}
}

func (s *Server) httpRun(ctx context.Context) error {

	mux := http.NewServeMux()

	mux.HandleFunc(s.serviceHTTP.endpoint, s.serviceHTTP.service.SendService(ctx))

	s.logger.Infof("server HTTP run on port %d", s.config.ServiceHTTPPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.config.ServiceHTTPPort), mux)
	if err != nil {
		s.logger.Infof("httpRun.ListenAndServe: %v", err)
		return fmt.Errorf("httpRun.ListenAndServe: %v", err)
	}

	return nil
}

func (s *Server) grpcRun() error {

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.ServiceGRPCPort))
	if err != nil {
		s.logger.Infof("grpcRun.Listen: %v", err)
		return fmt.Errorf("grpcRun.Listen: %v", err)
	}

	interceptorGRPC := interceptors.NewinterceptorGRPC(s.logger)

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptorGRPC.Logger),
		grpc.ChainUnaryInterceptor(s.grpcAuth.AuthInterceptor),
	}

	srv := grpc.NewServer(opts...)

	proto.RegisterMailServiceServer(srv, s.servicesGRPC)

	s.logger.Infof("server GRPC run on port %d", s.config.ServiceGRPCPort)
	err = srv.Serve(l)
	if err != nil {
		s.logger.Infof("grpcRun.Serve: %v", err)
		return fmt.Errorf("grpcRun.Serve: %v", err)
	}

	return nil
}
