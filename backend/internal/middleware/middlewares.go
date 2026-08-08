package middleware

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/prashanth30-n/Pubo/internal/server"
)

type MiddleWares struct {
	Global          *GlobalMiddleWares
	Auth            *AuthMiddleware
	ContextEnhancer *ContextEnhancer
	Tracing         *TracingMiddleware
	RateLimit       *RateLimitMiddleware
}

func NewMiddlewares(s *server.Server) *MiddleWares {
	var nrApp *newrelic.Application
	if s.LoggerService != nil {
		nrApp = s.LoggerService.GetApplication()
	}
	return &MiddleWares{
		Global:          NewGlobalMiddleWares(s),
		Auth:            NewAuthMiddleWare(s),
		ContextEnhancer: NewContextEnhancer(s),
		Tracing:         NewTracingMiddleware(s, nrApp),
		RateLimit:       NewRateLimitMiddleware(s),
	}
}
