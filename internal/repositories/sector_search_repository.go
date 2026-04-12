package repositories

import (
	"context"
	"fmt"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type SectorSearchRepository interface {
	UpsertNewSectors(ctx context.Context, sectors []models.SectorNew) ([]models.BasicResponse, error)
	UpsertNewSubSectors(ctx context.Context, subSectors []models.SubSector) ([]models.BasicResponse, error)
}

type sectorSearchRepository struct {
	pool *pgxpool.Pool
}

func NewSectorSearchRepository(pool *pgxpool.Pool) SectorSearchRepository {
	return &sectorSearchRepository{
		pool: pool,
	}
}

func (r *sectorSearchRepository) UpsertNewSectors(ctx context.Context, sectors []models.SectorNew) ([]models.BasicResponse, error) {
	if len(sectors) == 0 {
		return make([]models.BasicResponse, 0), nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO idxstock.sector (id, name, last_modified)
		VALUES ($1, $2, now())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			last_modified = now()
		WHERE (idxstock.sector.name IS DISTINCT FROM EXCLUDED.name)
		RETURNING id, name
	`

	batch := &pgx.Batch{}
	for _, s := range sectors {
		batch.Queue(query, s.Id, s.Name)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	updated := make([]models.BasicResponse, 0)
	for i := 0; i < len(sectors); i++ {
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

	logrus.Debugf("Successfully upserted %d new sectors", len(updated))
	return updated, nil
}

func (r *sectorSearchRepository) UpsertNewSubSectors(ctx context.Context, subSectors []models.SubSector) ([]models.BasicResponse, error) {
	if len(subSectors) == 0 {
		return make([]models.BasicResponse, 0), nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO idxstock.sub_sector (id, name, sector_id, last_modified)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			sector_id = EXCLUDED.sector_id,
			last_modified = now()
		WHERE (idxstock.sub_sector.name IS DISTINCT FROM EXCLUDED.name OR
		       idxstock.sub_sector.sector_id IS DISTINCT FROM EXCLUDED.sector_id)
		RETURNING id, name
	`

	batch := &pgx.Batch{}
	for _, sub := range subSectors {
		batch.Queue(query, sub.Id, sub.Name, sub.SectorId)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	updated := make([]models.BasicResponse, 0)
	for i := 0; i < len(subSectors); i++ {
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

	logrus.Debugf("Successfully upserted %d new sub-sectors", len(updated))
	return updated, nil
}
