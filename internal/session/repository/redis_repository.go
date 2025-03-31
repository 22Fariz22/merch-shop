package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/22Fariz22/merch-shop/config"
	"github.com/22Fariz22/merch-shop/internal/models"
	"github.com/22Fariz22/merch-shop/internal/session"
	"github.com/22Fariz22/merch-shop/pkg/logger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const (
	basePrefix = "api-session:"
)

// Session repository
type sessionRepo struct {
	redisClient *redis.Client
	basePrefix  string
	cfg         *config.Config
	logger      logger.Logger
}

// Session repository constructor
func NewSessionRepository(redisClient *redis.Client, cfg *config.Config, log logger.Logger) session.SessRepository {
	return &sessionRepo{redisClient: redisClient, basePrefix: basePrefix, cfg: cfg, logger: log}
}

// Create session in redis
func (s *sessionRepo) CreateSession(ctx context.Context, sess *models.Session, expire int) (string, error) {
	s.logger.Debug("here session redisRepo CreateSession")
	sess.SessionID = uuid.New().String()
	sessionKey := s.createKey(sess.SessionID)

	sessBytes, err := json.Marshal(&sess)
	if err != nil {
		return "", errors.WithMessage(err, "sessionRepo.CreateSession.json.Marshal")
	}
	if err = s.redisClient.Set(ctx, sessionKey, sessBytes, time.Second*time.Duration(expire)).Err(); err != nil {
		return "", errors.Wrap(err, "sessionRepo.CreateSession.redisClient.Set")
	}
	return sessionKey, nil
}

// Get session by id
func (s *sessionRepo) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	sessBytes, err := s.redisClient.Get(ctx, sessionID).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "sessionRep.GetSessionByID.redisClient.Get")
	}

	sess := &models.Session{}
	if err = json.Unmarshal(sessBytes, &sess); err != nil {
		return nil, errors.Wrap(err, "sessionRepo.GetSessionByID.json.Unmarshal")
	}
	return sess, nil
}

func (s *sessionRepo) createKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", s.basePrefix, sessionID)
}
