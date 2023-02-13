package decision

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
	"io.github.aliaksandrrachko/decision-helper/webapp/entities"
	"io.github.aliaksandrrachko/decision-helper/webapp/pkg/dbcontext"
)

type Service interface {
	Get(ctx context.Context, id int64) (Decision, error)
	Query(ctx context.Context, offset, limit int) ([]Decision, error)
	Create(ctx context.Context, input Decision) (int64, error)
	Update(ctx context.Context, id int64, input Decision) error
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repo                  Repository
	transactionalReadOnly dbcontext.TransactionFunc
	transactional         dbcontext.TransactionFunc
	logger                logrus.Logger
}

func NewService(repo Repository, dbcontext *dbcontext.DB, logger logrus.Logger) Service {
	return service{
		repo,
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return dbcontext.Transactional(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: true}, f)
		},
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return dbcontext.Transactional(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false}, f)
		},
		logger,
	}
}

type Decision struct {
	entities.Decision
}

func (s service) Get(ctx context.Context, id int64) (decision Decision, err error) {
	s.transactionalReadOnly(ctx, func(ctx context.Context) error {
		var dec entities.Decision
		dec, err = s.repo.Get(ctx, id)
		decision = Decision{dec}
		return err
	})
	return decision, err
}

func (s service) Query(ctx context.Context, offset, limit int) (decisions []Decision, err error) {
	s.transactionalReadOnly(ctx, func(ctx context.Context) error {
		var items []entities.Decision
		items, err = s.repo.Query(ctx, offset, limit)

		decisions = make([]Decision, 0)
		for _, item := range items {
			decisions = append(decisions, Decision{item})
		}
		return err
	})
	return decisions, err
}

func (s service) Create(ctx context.Context, input Decision) (id int64, err error) {
	s.transactional(ctx, func(ctx context.Context) error {
		id, err = s.repo.Save(ctx, input.Decision)
		return err
	})
	return id, err
}

func (s service) Update(ctx context.Context, id int64, input Decision) (err error) {
	return s.transactional(ctx, func(ctx context.Context) error {
		_, err = s.repo.Save(ctx, input.Decision)
		return err
	})
}

func (s service) Delete(ctx context.Context, id int64) error {
	return s.transactional(ctx, func(ctx context.Context) (err error) {
		_, err = s.repo.Delete(ctx, id)
		return err
	})
}
