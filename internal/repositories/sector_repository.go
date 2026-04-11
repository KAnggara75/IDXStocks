package repositories

import (
	"context"
	"fmt"

	"github.com/KAnggara75/IDXStocks/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type SectorRepository interface {
	UpsertSectors(ctx context.Context, sectors []models.PasardanaSector) ([]models.SectorResponse, error)
}

type sectorRepository struct {
	pool *pgxpool.Pool
}

func NewSectorRepository(pool *pgxpool.Pool) SectorRepository {
	return &sectorRepository{
		pool: pool,
	}
}

func (r *sectorRepository) UpsertSectors(ctx context.Context, sectors []models.PasardanaSector) ([]models.SectorResponse, error) {
	if len(sectors) == 0 {
		return make([]models.SectorResponse, 0), nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO idxstock.sectors (id, code, name, name_en, description, last_modified)
		VALUES ($1, $2, $3, $4, $5, now())
		ON CONFLICT (id) DO UPDATE SET
			code = EXCLUDED.code,
			name = EXCLUDED.name,
			name_en = EXCLUDED.name_en,
			description = EXCLUDED.description,
			last_modified = now()
		RETURNING code, name
	`

	batch := &pgx.Batch{}
	for _, s := range sectors {
		batch.Queue(query, s.Id, s.Code, s.Name, s.NameEn, s.Description)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	return r.executeUpsertBatch(br, len(sectors), tx)
}

func (r *sectorRepository) executeUpsertBatch(br pgx.BatchResults, count int, tx pgx.Tx) ([]models.SectorResponse, error) {
	updatedSectors := make([]models.SectorResponse, 0)
	for i := 0; i < count; i++ {
		rows, err := br.Query()
		if err != nil {
			return nil, fmt.Errorf("failed to execute batch statement %d: %w", i, err)
		}
		
		for rows.Next() {
			var sr models.SectorResponse
			if err := rows.Scan(&sr.Code, &sr.Name); err == nil {
				updatedSectors = append(updatedSectors, sr)
			}
		}
		rows.Close()
	}

	if err := br.Close(); err != nil {
		return nil, fmt.Errorf("failed to close batch result: %w", err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Debugf("Successfully upserted %d sectors", len(updatedSectors))
	return updatedSectors, nil
}
