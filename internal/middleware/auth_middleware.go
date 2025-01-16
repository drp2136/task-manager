package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		fields := strings.Fields(tokenString)
		accessToken := fields[1]
		if len(fields) == 2 {
			authorizationType := fields[0]
			if !strings.EqualFold(authorizationType, "Bearer") {
				fmt.Println("unsupported authorization type " + authorizationType)
				return
			}

			fmt.Println(accessToken)
		} else {
			fmt.Println("invalid authorization header format")
		}

		// token, err := auth.ParseToken(accessToken)
		// if err != nil || !token.Valid {
		// 	http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		// 	return
		// }

		next.ServeHTTP(w, r)
	})
}
