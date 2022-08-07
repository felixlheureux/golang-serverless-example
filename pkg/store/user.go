package store

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/manta-coder/golang-serverless-example/pkg/db"
	"github.com/manta-coder/golang-serverless-example/pkg/domain"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
	"time"
)

type UserStore interface {
	Get(userID string) (domain.User, error)
	FindByEthereumAddress(ethereumAddressHex string) (domain.User, error)
	Store(user domain.User) (domain.User, error)
	Update(user domain.User) (domain.User, error)
	Remove(userID string) error
}

type userStore struct {
	logger *zap.SugaredLogger
	db     *sqlx.DB
}

func NewUserStore(logger *zap.SugaredLogger, db *sqlx.DB) UserStore {
	return &userStore{logger, db}
}

func (s *userStore) Get(userID string) (domain.User, error) {
	var result domain.User

	query, args, _ := sq.Select(usersColumns...).
		From(usersTable).
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()

	if err := s.db.Get(&result, query, args...); err != nil {
		return result, db.QueryExecuteError(err, query, args)
	}

	return result, nil
}

func (s *userStore) FindByEthereumAddress(ethereumAddressHex string) (domain.User, error) {
	var result domain.User

	query, args, _ := sq.Select(usersColumns...).
		From(usersTable).
		Where(squirrel.Eq{"ethereum_address": ethereumAddressHex}).
		ToSql()

	err := s.db.Get(&result, query, args...)
	switch err {
	case nil:
		return result, nil
	case sql.ErrNoRows:
		return domain.User{}, nil
	default:
		return result, db.QueryExecuteError(err, query, args)
	}
}

func (s *userStore) Store(user domain.User) (domain.User, error) {
	now := time.Now()

	user.UserID = "usr_" + ksuid.New().String()
	user.CreatedAt = now
	user.UpdatedAt = now

	query, args, _ := sq.Insert(usersTable).
		Columns(usersColumns...).
		Values(
			user.UserID,
			user.EthereumAddressHex,
			user.Username,
			user.DefaultCharacterID,
			user.UpdatedAt,
			user.CreatedAt,
		).
		ToSql()

	if _, err := s.db.Exec(query, args...); err != nil {
		return user, db.QueryExecuteError(err, query, args)
	}

	return user, nil
}

func (s *userStore) Update(user domain.User) (domain.User, error) {
	now := time.Now()

	user.UpdatedAt = now

	query, args, _ := sq.Update(usersTable).
		Set("username", user.Username).
		Set("default_character_id", user.DefaultCharacterID).
		Set("updated_at", user.UpdatedAt).
		Where(squirrel.Eq{"user_id": user.UserID}).
		ToSql()

	if _, err := s.db.Exec(query, args...); err != nil {
		return user, db.QueryExecuteError(err, query, args)
	}

	return user, nil
}

func (s *userStore) Remove(userID string) error {
	query, args, _ := sq.Delete(usersTable).
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()

	if _, err := s.db.Exec(query, args...); err != nil {
		return db.QueryExecuteError(err, query, args)
	}

	return nil
}
