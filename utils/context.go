package utils

import (
	"github.com/google/uuid"
	"net/http"
)

func HasRole(r *http.Request, requiredRole string) bool {
	claims, ok := r.Context().Value("user").(*CustomClaims)
	if !ok {
		return false
	}
	return claims.Role == requiredRole
}

//
//func IsLoggedIn(r *http.Request) bool {
//	claims, ok := r.Context().Value("user").(*CustomClaims)
//	if !ok {
//		return false
//	}
//	return HasRole(r, claims.UserID)
//}

func GetClaims(r *http.Request) (*CustomClaims, bool) {
	claims, ok := r.Context().Value("user").(*CustomClaims)
	return claims, ok
}

func GetUserID(r *http.Request) (uuid.UUID, bool) {
	claims, ok := GetClaims(r)
	if !ok {
		return uuid.Nil, false
	}

	userUUID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, false
	}

	return userUUID, true
}
