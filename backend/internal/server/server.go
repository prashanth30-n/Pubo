package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/prashanth30-n/Pubo/internal/config"
	"github.com/prashanth30-n/Pubo/internal/database"
	"github.com/prashanth30-n/Pubo/internal/lib/job"
	loggerPkg "github.com/prashanth30-n/Pubo/internal/logger"
)

type Server struct {
	Config        *config.Config
	Logger        *zerolog.Logger
	LoggerService *loggerPkg.LoggerService
	DB            *database.Database
	Redis         *redis.Client
	httpServer    *http.Server
	Job           *job.JobService
}

func New(cfg *config.Config, logger *zerolog.Logger, loggerService *loggerPkg.LoggerService) (*Server, error) {
	db, err := database.New(cfg, logger, loggerService)
	if err != nil {
		return nil, fmt.Errorf("failed to intialise database %w", err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
	})

	if loggerService != nil && loggerService.GetApplication() != nil {
		redisClient.AddHook(nrredis.NewHook(redisClient.Options()))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var jobService *job.JobService
	if err := redisClient.Ping(ctx).Err(); err != nil {
		// Do not start Asynq without Redis. Starting it after a failed ping causes
		// several independent worker loops to retry and flood the application log.
		logger.Warn().Err(err).Msg("redis unavailable; background jobs are disabled")
	} else {
		jobService = job.NewJobService(logger, cfg)
		jobService.InitHandlers(cfg, logger)
		if err := jobService.Start(); err != nil {
			return nil, err
		}
	}
	server := &Server{
		Config:        cfg,
		Logger:        logger,
		LoggerService: loggerService,
		DB:            db,
		Redis:         redisClient,
		Job:           jobService,
	}
	return server, nil
}

func (s *Server) SetupHttpServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:         ":" + s.Config.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(s.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.Config.Server.IdleTimeout) * time.Second,
	}
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not intiliased")
	}
	s.Logger.Info().Str("port", s.Config.Server.Port).Str("env", s.Config.Primary.Env).Msg("starting server")

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}
	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection %w", err)
	}
	if s.Job != nil {
		s.Job.Stop()
	}
	return nil
}
