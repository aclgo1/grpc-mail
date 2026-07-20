package grpcauth

import (
	"context"
	"crypto/rsa"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/aclgo/grpc-mail/config"
	jwt "github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GrpcAuth struct {
	publicKey *rsa.PublicKey
}

func NewGrpcAuth(cfg *config.Config) *GrpcAuth {
	g := GrpcAuth{}

	pubKeyData, err := os.ReadFile(cfg.PathPublicPem)
	if err != nil {
		log.Fatalf("NewGrpcAuth:os.ReadFile: %v", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if err != nil {
		log.Fatalf("jwt.ParseRSAPublicKeyFromPEM: %v", err)
	}

	g.publicKey = publicKey

	return &g
}

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrMissingMetadata   = errors.New("missing metadata")
	ErrTokenNotProvided  = errors.New("token not provided")
	ErrUnexpectedtMethod = errors.New("invalid signing method")
)

func (g *GrpcAuth) AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	}

	tokens := md["authorization"]

	if len(tokens) == 0 {
		return nil, status.Error(codes.Unauthenticated, ErrTokenNotProvided.Error())
	}

	ttkStr := strings.TrimPrefix(tokens[0], "Bearer ")

	ttk, err := jwt.Parse(ttkStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, status.Error(codes.Unauthenticated, ErrUnexpectedtMethod.Error())
		}

		return g.publicKey, nil
	})

	if err != nil || !ttk.Valid {
		return nil, status.Error(codes.Unauthenticated, ErrInvalidToken.Error())
	}

	return handler(ctx, req)
}
