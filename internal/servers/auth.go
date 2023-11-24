package gapi

// import (
// 	"context"

// 	"github.com/iKayrat/auth-grpc/internal/services/pb"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// func (server *Server) Login(ctx context.Context, req *pb.LoginRequestMessage) (*pb.LoginResponseMessage, error) {
// 	user, err := auth.LoginUser(req.GetUsername(), req.GetPassword(), server.dbCollection.Users, context.TODO())
// 	if err != nil {
// 		if err == auth.ErrUserNotFound {
// 			return nil, status.Errorf(codes.NotFound, err.Error())
// 		}
// 		if err == auth.ErrInvalidCredentials {
// 			return nil, status.Errorf(codes.Unauthenticated, err.Error())
// 		}
// 		return nil, status.Errorf(codes.Internal, err.Error())
// 	}
// 	//

// 	resp := auth.ConvertUserObjectToUser(user)

// 	token, err := auth.CreateToken(server.tokenMaker, resp.Username, int64(resp.Id), server.config.AccessTokenDuration)
// 	if err != nil {q
// 		return nil, status.Errorf(codes.Internal, err.Error())
// 	}

// 	return &pb.LoginResponseMessage{
// 		User:        resp,
// 		AccessToken: token,
// 	}, nil
// }

// import (
// 	"context"
// 	"encoding/base64"
// 	"strings"

// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/credentials"
// 	"google.golang.org/grpc/metadata"
// 	"google.golang.org/grpc/peer"
// 	"google.golang.org/grpc/status"
// )

// type userMDKey struct{}

// // UserMetadata contains metadata about a user.
// type UserMetadata struct {
// 	ID string
// }

// // Authenticator exposes a function for authenticating requests.
// type Authenticator struct {
// 	Username string
// 	Password string
// 	Token    string
// }

// // Authenticate checks that a token exists and is valid. It stores the user
// // metadata in the returned context and removes the token from the context.
// func (a Authenticator) Authenticate(ctx context.Context) (newCtx context.Context, err error) {
// 	defer func() {
// 		if err == nil {
// 			// Store user metadata
// 			userMD := &UserMetadata{
// 				ID: a.Username,
// 			}
// 			newCtx = context.WithValue(newCtx, userMDKey{}, userMD)
// 		}
// 	}()
// 	err = a.tryTLSAuth(ctx)
// 	if err == nil {
// 		return ctx, nil
// 	}

// 	newCtx, err = a.tryTokenAuth(ctx)
// 	if err == nil {
// 		return newCtx, nil
// 	}

// 	return a.tryBasicAuth(ctx)
// }

// func (a Authenticator) tryTLSAuth(ctx context.Context) error {
// 	p, ok := peer.FromContext(ctx)
// 	if !ok {
// 		return status.Error(codes.Unauthenticated, "no peer found")
// 	}

// 	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)
// 	if !ok {
// 		return status.Error(codes.Unauthenticated, "unexpected peer transport credentials")
// 	}

// 	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
// 		return status.Error(codes.Unauthenticated, "could not verify peer certificate")
// 	}

// 	if tlsAuth.State.VerifiedChains[0][0].Subject.CommonName != a.Username {
// 		return status.Error(codes.Unauthenticated, "invalid subject common name")
// 	}

// 	return nil
// }

// func (a Authenticator) tryBasicAuth(ctx context.Context) (context.Context, error) {
// 	auth, err := extractHeader(ctx, "authorization")
// 	if err != nil {
// 		return ctx, err
// 	}

// 	const prefix = "Basic "
// 	if !strings.HasPrefix(auth, prefix) {
// 		return ctx, status.Error(codes.Unauthenticated, `missing "Basic " prefix in "Authorization" header`)
// 	}

// 	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
// 	if err != nil {
// 		return ctx, status.Error(codes.Unauthenticated, `invalid base64 in header`)
// 	}

// 	cs := string(c)
// 	s := strings.IndexByte(cs, ':')
// 	if s < 0 {
// 		return ctx, status.Error(codes.Unauthenticated, `invalid basic auth format`)
// 	}

// 	user, password := cs[:s], cs[s+1:]
// 	if user != a.Username || password != a.Password {
// 		return ctx, status.Error(codes.Unauthenticated, "invalid user or password")
// 	}

// 	// Remove token from headers from here on
// 	return purgeHeader(ctx, "authorization"), nil
// }

// func (a Authenticator) tryTokenAuth(ctx context.Context) (context.Context, error) {
// 	auth, err := extractHeader(ctx, "authorization")
// 	if err != nil {
// 		return ctx, err
// 	}

// 	const prefix = "Bearer "
// 	if !strings.HasPrefix(auth, prefix) {
// 		return ctx, status.Error(codes.Unauthenticated, `missing "Bearer " prefix in "Authorization" header`)
// 	}

// 	if strings.TrimPrefix(auth, prefix) != a.Token {
// 		return ctx, status.Error(codes.Unauthenticated, "invalid token")
// 	}

// 	// Remove token from headers from here on
// 	return purgeHeader(ctx, "authorization"), nil
// }

// func extractHeader(ctx context.Context, header string) (string, error) {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return "", status.Error(codes.Unauthenticated, "no headers in request")
// 	}

// 	authHeaders, ok := md[header]
// 	if !ok {
// 		return "", status.Error(codes.Unauthenticated, "no header in request")
// 	}

// 	if len(authHeaders) != 1 {
// 		return "", status.Error(codes.Unauthenticated, "more than 1 header in request")
// 	}

// 	return authHeaders[0], nil
// }

// func purgeHeader(ctx context.Context, header string) context.Context {
// 	md, _ := metadata.FromIncomingContext(ctx)
// 	mdCopy := md.Copy()
// 	mdCopy[header] = nil
// 	return metadata.NewIncomingContext(ctx, mdCopy)
// }

// // GetUserMetadata can be used to extract user metadata stored in a context.
// func GetUserMetadata(ctx context.Context) (*UserMetadata, bool) {
// 	userMD := ctx.Value(userMDKey{})

// 	switch md := userMD.(type) {
// 	case *UserMetadata:
// 		return md, true
// 	default:
// 		return nil, false
// 	}
// }
