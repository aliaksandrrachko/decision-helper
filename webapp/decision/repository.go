package decision

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
	"io.github.aliaksandrrachko/decision-helper/webapp/entities"
	"io.github.aliaksandrrachko/decision-helper/webapp/pkg/dbcontext"
)

type Repository interface {
	Get(ctx context.Context, id int64) (entities.Decision, error)
	Query(ctx context.Context, offset, limit int) ([]entities.Decision, error)
	Count(ctx context.Context) (int, error)
	Save(ctx context.Context, decision entities.Decision) (int64, error)
	Delete(ctx context.Context, id int64) (int, error)
}

type repository struct {
	db     *dbcontext.DB
	logger logrus.Logger
}

func NewRepository(db *dbcontext.DB, logger logrus.Logger) Repository {
	return repository{db: db, logger: logger}
}

func (r repository) Get(ctx context.Context, id int64) (entities.Decision, error) {
	var dec entities.Decision

	row := r.db.With(ctx).QueryRowContext(ctx,
		`SELECT 
	 	dec.dec_id,
		dec.title,
		dec.navi_user,
		dec.navi_date
		FROM decision dec 
		WHERE dec_id = $1`,
		id,
	)

	if err := row.Scan(&dec.Id, &dec.Title, &dec.NavigationUser, &dec.NavigationDate); err != nil {
		if err == sql.ErrNoRows {
			return dec, fmt.Errorf("decision %d: no such decision: %v", id, err)
		}
		return dec, fmt.Errorf("get %d: %v", id, err)
	}
	return dec, nil
}

func (r repository) Query(ctx context.Context, offset, limit int) (decisions []entities.Decision, err error) {
	rows, err := r.db.With(ctx).QueryContext(ctx,
		fmt.Sprintf(
			`SELECT 
	 		dec.dec_id,
			dec.title,
			dec.navi_user,
			dec.navi_date
			FROM decision dec 
			limit %d offset %d`, limit, offset),
	)
	if err != nil {
		return nil, fmt.Errorf("decisions not found: %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var dec entities.Decision
		if err := rows.Scan(&dec.Id, &dec.Title, &dec.NavigationUser, &dec.NavigationDate); err != nil {
			return nil, fmt.Errorf("decisions not found: %v", err)
		}
		decisions = append(decisions, dec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("decisions not found: %v", err)
	}

	return decisions, nil
}

func (r repository) Count(ctx context.Context) (count int, err error) {
	if err := r.db.With(ctx).QueryRowContext(ctx,
		`SELECT 
	 	count(*)
		FROM decision dec`,
	).Scan(count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r repository) Save(ctx context.Context, decision entities.Decision) (int64, error) {
	if decision.Id == 0 {
		return r.create(ctx, decision)
	} else {
		r.update(ctx, decision)
		return decision.Id, nil
	}
}

func (r repository) Delete(ctx context.Context, id int64) (int, error) {
	result, err := r.db.With(ctx).ExecContext(ctx,
		`DELETE decision 
		WHERE dec_id = $1`,
		id,
	)

	if err != nil {
		return 0, fmt.Errorf("delete decision: %v", err)
	}

	if count, err := result.RowsAffected(); err != nil {
		return 0, fmt.Errorf("delete decision: %v", err)
	} else {
		return int(count), nil
	}
}

func (r repository) create(ctx context.Context, decision entities.Decision) (int64, error) {
	var id int64
	row := r.db.With(ctx).QueryRowContext(ctx,
		`INSERT INTO decision 
		(dec_id, title, navi_user, navi_date) 
		VALUES 
		(nextval('dec_seq'), $1, $2, NOW()) returning dec_id`,
		decision.Title, decision.NavigationUser,
	)

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("decision %d: %v", id, err)
	}

	return id, nil

}

func (r repository) update(ctx context.Context, decision entities.Decision) (int, error) {
	result, err := r.db.With(ctx).ExecContext(ctx,
		`UPDATE decision 
		SET title = $1,
		navi_user = $2,
		navi_date = NOW()
		WHERE dec_id = $3`,
		decision.Title, decision.NavigationUser, decision.Id,
	)

	if err != nil {
		return 0, fmt.Errorf("update decision: %v", err)
	}

	if id, err := result.RowsAffected(); err != nil {
		return int(id), fmt.Errorf("update decision: %v", err)
	} else {
		return int(id), nil
	}
}
