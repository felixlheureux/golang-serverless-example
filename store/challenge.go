package store

import (
	"github.com/Masterminds/squirrel"
	"github.com/childrenofukiyo/odin/pkg/db"
	"github.com/childrenofukiyo/odin/pkg/domain"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
	"time"
)

type ChallengeStore interface {
	Get(ethereumAddress string) (domain.Challenge, error)
	Store(challenge domain.Challenge) (domain.Challenge, error)
	Remove(ethereumAddress string) error
}

type challengeStore struct {
	logger *zap.SugaredLogger
	db     *sqlx.DB
}

func NewChallengeStore(logger *zap.SugaredLogger, db *sqlx.DB) ChallengeStore {
	return &challengeStore{logger, db}
}

func (s *challengeStore) Get(ethereumAddressHex string) (domain.Challenge, error) {
	var result domain.Challenge

	query, args, _ := sq.Select(challengesColumns...).
		From(challengesTable).
		Where(squirrel.Eq{"ethereum_address": ethereumAddressHex}).
		OrderBy("created_at DESC").
		ToSql()

	if err := s.db.Get(&result, query, args...); err != nil {
		return result, db.QueryExecuteError(err, query, args)
	}

	return result, nil
}

func (s *challengeStore) Store(challenge domain.Challenge) (domain.Challenge, error) {
	now := time.Now()

	challenge.ChallengeID = "chl_" + ksuid.New().String()
	challenge.CreatedAt = now

	query, args, _ := sq.Insert(challengesTable).
		Columns(challengesColumns...).
		Values(
			challenge.ChallengeID,
			challenge.EthereumAddressHex,
			challenge.Challenge,
			challenge.CreatedAt,
		).
		ToSql()

	if _, err := s.db.Exec(query, args...); err != nil {
		return challenge, db.QueryExecuteError(err, query, args)
	}

	return challenge, nil
}

func (s *challengeStore) Remove(ethereumAddress string) error {
	query, args, _ := sq.Delete(challengesTable).
		Where(squirrel.Eq{"ethereum_address": ethereumAddress}).
		ToSql()

	if _, err := s.db.Exec(query, args...); err != nil {
		return db.QueryExecuteError(err, query, args)
	}

	return nil
}
