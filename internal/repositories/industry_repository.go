package repositories

import (
	"context"
	"fmt"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type IndustryRepository interface {
	UpsertIndustries(ctx context.Context, industries []models.Industry) ([]models.BasicResponse, error)
	UpsertSubIndustries(ctx context.Context, subIndustries []models.SubIndustry) ([]models.BasicResponse, error)
}

type industryRepository struct {
	pool *pgxpool.Pool
}

func NewIndustryRepository(pool *pgxpool.Pool) IndustryRepository {
	return &industryRepository{
		pool: pool,
	}
}

func (r *industryRepository) UpsertIndustries(ctx context.Context, industries []models.Industry) ([]models.BasicResponse, error) {
	if len(industries) == 0 {
		return make([]models.BasicResponse, 0), nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO idxstock.industry (id, name, last_modified)
		VALUES ($1, $2, now())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			last_modified = now()
		RETURNING id, name
	`

	batch := &pgx.Batch{}
	for _, ind := range industries {
		batch.Queue(query, ind.Id, ind.Name)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	updated := make([]models.BasicResponse, 0)
	for i := 0; i < len(industries); i++ {
		var res models.BasicResponse
		err := br.QueryRow().Scan(&res.Id, &res.Name)
		if err != nil {
			// Skip if RETURNING didn't return anything (though for UPSERT it usually should)
			continue
		}
		updated = append(updated, res)
	}

	if err := br.Close(); err != nil {
		return nil, fmt.Errorf("failed to close batch result: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Debugf("Successfully upserted %d industries", len(updated))
	return updated, nil
}

func (r *industryRepository) UpsertSubIndustries(ctx context.Context, subIndustries []models.SubIndustry) ([]models.BasicResponse, error) {
	if len(subIndustries) == 0 {
		return make([]models.BasicResponse, 0), nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO idxstock.sub_industry (id, name, industry_id, last_modified)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			industry_id = EXCLUDED.industry_id,
			last_modified = now()
		RETURNING id, name
	`

	batch := &pgx.Batch{}
	for _, sub := range subIndustries {
		batch.Queue(query, sub.Id, sub.Name, sub.IndustryId)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	updated := make([]models.BasicResponse, 0)
	for i := 0; i < len(subIndustries); i++ {
		var res models.BasicResponse
		err := br.QueryRow().Scan(&res.Id, &res.Name)
		if err != nil {
			continue
		}
		updated = append(updated, res)
	}

	if err := br.Close(); err != nil {
		return nil, fmt.Errorf("failed to close batch result: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Debugf("Successfully upserted %d sub-industries", len(updated))
	return updated, nil
}
