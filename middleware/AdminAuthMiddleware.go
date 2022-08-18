package middleware

import (
	"audioPhile/database/helper"
	"audioPhile/models"
	"net/http"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := helper.GetContextData(r)
		if ctx == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		signedUserRole := ctx.Role
		if signedUserRole != string(models.UserRoleAdmin) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
