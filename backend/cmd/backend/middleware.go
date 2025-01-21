package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const connectionIDKey contextKey = "connectionID"

// Middleware to generate unique connection ID
func connectionIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a new UUID
		connID := uuid.New().String()
		connID = connID[:8]

		// Add the connection ID to the request context
		ctx := context.WithValue(r.Context(), connectionIDKey, connID)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
