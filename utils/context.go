package utils

import (
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
