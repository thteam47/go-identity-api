package models

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserContext   UserContext  `bson:"user_context,omitempty"`
	PermissionAll bool         `bson:"permission_all,omitempty"`
	Role          string       `bson:"role,omitempty"`
	Permissions   []Permission `bson:"permissions,omitempty"`
	jwt.StandardClaims
}
