package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kholodihor/cows-shelter-backend/storage"
)

type contextKey string

const storageKey contextKey = "storage"

// StorageMiddleware injects the storage service into the request context
func StorageMiddleware(store storage.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add the storage service to the request context
			ctx := context.WithValue(r.Context(), storageKey, store)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetStorage retrieves the storage service from the request context
func GetStorage(ctx context.Context) storage.Service {
	if store, ok := ctx.Value(storageKey).(storage.Service); ok {
		return store
	}
	return nil
}

// GinStorageMiddleware creates a Gin middleware that injects the storage service
func GinStorageMiddleware(store storage.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add the storage service to the request context
		ctx := context.WithValue(c.Request.Context(), storageKey, store)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
