// Package middleware provides HTTP and gRPC middleware for Phonic AI Calling Agent
package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/ArbajAnsari19/phonic/pkg/logger"
)

// generateID generates a random ID for tracing
func generateID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID
		return hex.EncodeToString([]byte(time.Now().Format("20060102150405")))
	}
	return hex.EncodeToString(bytes)
}

// HTTPTracing middleware adds tracing information to HTTP requests
func HTTPTracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Generate or extract trace ID
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = generateID()
		}
		
		// Generate request ID
		requestID := generateID()
		
		// Add to response headers
		w.Header().Set("X-Trace-ID", traceID)
		w.Header().Set("X-Request-ID", requestID)
		
		// Create context with tracing information
		ctx := context.WithValue(r.Context(), logger.TraceIDKey, traceID)
		ctx = context.WithValue(ctx, logger.RequestIDKey, requestID)
		
		// Create wrapped response writer to capture status code
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		// Process request
		next.ServeHTTP(wrappedWriter, r.WithContext(ctx))
		
		// Log request
		duration := time.Since(start)
		logger.WithContext(ctx).LogHTTPRequest(
			r.Method,
			r.URL.Path,
			r.UserAgent(),
			r.RemoteAddr,
			wrappedWriter.statusCode,
			duration,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// HTTPLogging middleware provides detailed HTTP request logging
func HTTPLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Get logger with context
		log := logger.WithContext(r.Context())
		
		// Log incoming request
		log.Info("HTTP request started",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.String("user_agent", r.UserAgent()),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int64("content_length", r.ContentLength),
		)
		
		// Process request
		next.ServeHTTP(w, r)
		
		// Log completion
		duration := time.Since(start)
		log.Info("HTTP request completed",
			zap.Duration("duration", duration),
		)
	})
}

// GRPCTracing interceptor adds tracing to gRPC requests
func GRPCTracingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	
	// Extract or generate trace ID
	var traceID string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ids := md.Get("trace-id"); len(ids) > 0 {
			traceID = ids[0]
		}
	}
	if traceID == "" {
		traceID = generateID()
	}
	
	// Generate request ID
	requestID := generateID()
	
	// Add to context
	ctx = context.WithValue(ctx, logger.TraceIDKey, traceID)
	ctx = context.WithValue(ctx, logger.RequestIDKey, requestID)
	
	// Add to outgoing metadata
	ctx = metadata.AppendToOutgoingContext(ctx, "trace-id", traceID, "request-id", requestID)
	
	// Process request
	resp, err := handler(ctx, req)
	
	// Log request
	duration := time.Since(start)
	logger.WithContext(ctx).LogGRPCRequest(info.FullMethod, duration, err)
	
	return resp, err
}

// GRPCLogging interceptor provides detailed gRPC request logging
func GRPCLoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	
	// Get logger with context
	log := logger.WithContext(ctx)
	
	// Log incoming request
	log.Info("gRPC request started",
		zap.String("method", info.FullMethod),
		zap.Any("request", req),
	)
	
	// Process request
	resp, err := handler(ctx, req)
	
	// Log completion
	duration := time.Since(start)
	if err != nil {
		log.Error("gRPC request failed",
			zap.Duration("duration", duration),
			zap.Error(err),
		)
	} else {
		log.Info("gRPC request completed",
			zap.Duration("duration", duration),
		)
	}
	
	return resp, err
}

// Recovery middleware handles panics gracefully
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log := logger.WithContext(r.Context())
				log.Error("HTTP handler panic",
					zap.Any("panic", err),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
				)
				
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}

// GRPCRecovery interceptor handles gRPC panics
func GRPCRecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log := logger.WithContext(ctx)
			log.Error("gRPC handler panic",
				zap.Any("panic", r),
				zap.String("method", info.FullMethod),
			)
			
			err = status.Error(codes.Internal, "Internal server error")
		}
	}()
	
	return handler(ctx, req)
}

// CORS middleware handles Cross-Origin Resource Sharing
func CORS(allowedOrigins []string, allowedMethods []string, allowedHeaders []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			
			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}
			
			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			
			// Set other CORS headers
			w.Header().Set("Access-Control-Allow-Methods", joinStrings(allowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", joinStrings(allowedHeaders, ", "))
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			
			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// joinStrings joins a slice of strings with a separator
func joinStrings(slice []string, separator string) string {
	if len(slice) == 0 {
		return ""
	}
	if len(slice) == 1 {
		return slice[0]
	}
	
	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += separator + slice[i]
	}
	return result
}
