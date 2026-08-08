package middleware

import "github.com/prashanth30-n/Pubo/internal/server"

type RateLimitMiddleware struct {
	server *server.Server
}

func NewRateLimitMiddleware(s *server.Server) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		server: s,
	}
}
func (r *RateLimitMiddleware) RecordRateLimitHit(endpoint string) {
	if r.server.LoggerService != nil && r.server.LoggerService.GetApplication() != nil {
		r.server.LoggerService.GetApplication().RecordCustomEvent("RatLimtiHit", map[string]interface{}{
			"endpoint": endpoint,
		})
	}
}
