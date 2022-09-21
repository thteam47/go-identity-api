package grpcauth

import (
	"context"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/thteam47/go-identity-api/pkg/models"
	"github.com/thteam47/go-identity-api/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	JwtKey string
}

func NewAuthInterceptor(jwtKey string) *AuthInterceptor {
	return &AuthInterceptor{JwtKey: jwtKey}
}

func (interceptor *AuthInterceptor) Authentication(ctx context.Context, ctxRequest *pb.Context, privilege string, action string) (*models.UserContext, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	accessToken := ""
	authorization := md["authorization"]
	if len(authorization) < 1 {
		if ctxRequest.AccessToken != "" {
			accessToken = ctxRequest.AccessToken
		} else {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}
	}
	if accessToken == "" {
		accessToken = strings.TrimPrefix(authorization[0], "Bearer ")
		if accessToken == "undefined" {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}
	}

	claims, err := interceptor.VerifyToken(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	if !claims.PermissionAll {
		return &claims.UserContext, nil
	}
	for _, permission := range claims.Permissions {
		if permission.Privilege == privilege {
			if contains(permission.Actions, action) {
				return &claims.UserContext, nil
			}
		}
	}
	return nil, status.Error(codes.PermissionDenied, "Fobbiden!")

}

func (interceptor *AuthInterceptor) VerifyToken(accessToken string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&models.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return []byte(interceptor.JwtKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil

}
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
