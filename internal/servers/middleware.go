package gapi

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader   = "authorization"
	authorizationBasic    = "Basic "
	authorizationUsername = "username"
	authorizationPassword = "password"
	authorizationBearer   = "bearer"
)

func (server *Server) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover:", err)
		}
	}()

	// return handler(ctx, req)
	if methodRequiresAuthentication(info.FullMethod) {
		auth, err := extractHeader(ctx, authorizationHeader)
		if err != nil {
			return ctx, err
		}

		const prefix = authorizationBasic
		if !strings.HasPrefix(auth, prefix) {
			return ctx, status.Error(codes.Unauthenticated, `missing "Basic " prefix in "Authorization" header`)
		}

		c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
		if err != nil {
			return ctx, status.Error(codes.Unauthenticated, `invalid base64 in header`)
		}

		cs := string(c)
		s := strings.IndexByte(cs, ':')
		if s < 0 {
			return ctx, status.Error(codes.Unauthenticated, `invalid basic auth format`)
		}

		username, password := cs[:s], cs[s+1:]

		user, err := server.Store.GetByUsername(ctx, username)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid user or password")
		}

		fmt.Println(username, password)
		if username != user.Username || password != user.Password {
			return ctx, status.Error(codes.Unauthenticated, "invalid user or password")
		}

		return handler(ctx, req)
	}
	return handler(ctx, req)
}

func extractHeader(ctx context.Context, header string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no headers in request")
	}

	authHeaders, ok := md[header]
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no header in request")
	}

	if len(authHeaders) != 1 {
		return "", status.Error(codes.Unauthenticated, "more than 1 header in request")
	}

	return authHeaders[0], nil
}
