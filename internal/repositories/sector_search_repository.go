package repositories

import (
	"context"
	"fmt"

	"github.com/KAnggara75/IDXStock/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type SectorSearchRepository interface {
	UpsertNewSectors(ctx context.Context, sectors []models.SectorNew) ([]models.BasicResponseWithCode, error)
	UpsertNewSubSectors(ctx context.Context, subSectors []models.SubSector) ([]models.BasicResponseWithCode, error)
}

type sectorSearchRepository struct {
	pool *pgxpool.Pool
}

func NewSectorSearchRepository(pool *pgxpool.Pool) SectorSearchRepository {
	return &sectorSearchRepository{
		pool: pool,
	}
}

func (r *sectorSearchRepository) UpsertNewSectors(ctx context.Context, sectors []models.SectorNew) ([]models.BasicResponseWithCode, error) {
	if len(sectors) == 0 {
		return make([]models.BasicResponseWithCode, 0), nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO idxstock.sector (id, code, name, name_en, description, last_modified)
		VALUES ($1, $2, $3, $4, $5, now())
		ON CONFLICT (id) DO UPDATE SET
			code = EXCLUDED.code,
			name = EXCLUDED.name,
			name_en = EXCLUDED.name_en,
			description = EXCLUDED.description,
			last_modified = now()
		WHERE (idxstock.sector.name IS DISTINCT FROM EXCLUDED.name OR
		       idxstock.sector.code IS DISTINCT FROM EXCLUDED.code OR
		       idxstock.sector.description IS DISTINCT FROM EXCLUDED.description OR
		       idxstock.sector.name_en IS DISTINCT FROM EXCLUDED.name_en)
		RETURNING id, code, name
	`

	batch := &pgx.Batch{}
	for _, s := range sectors {
		batch.Queue(query, s.Id, s.Code, s.Name, s.NameEn, s.Description)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	updated := make([]models.BasicResponseWithCode, 0)
	for i := 0; i < len(sectors); i++ {
		var res models.BasicResponseWithCode
		err := br.QueryRow().Scan(&res.Id, &res.Code, &res.Name)
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

func (r *sectorSearchRepository) UpsertNewSubSectors(ctx context.Context, subSectors []models.SubSector) ([]models.BasicResponseWithCode, error) {
	if len(subSectors) == 0 {
		return make([]models.BasicResponseWithCode, 0), nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO idxstock.sub_sector (id, code, name, name_en, description, sector_id, last_modified)
		VALUES ($1, $2, $3, $4, $5, $6, now())
		ON CONFLICT (id) DO UPDATE SET
			code = EXCLUDED.code,
			name = EXCLUDED.name,
			name_en = EXCLUDED.name_en,
			description = EXCLUDED.description,
			sector_id = EXCLUDED.sector_id,
			last_modified = now()
		WHERE (idxstock.sub_sector.name IS DISTINCT FROM EXCLUDED.name OR
		       idxstock.sub_sector.code IS DISTINCT FROM EXCLUDED.code OR
		       idxstock.sub_sector.description IS DISTINCT FROM EXCLUDED.description OR
		       idxstock.sub_sector.name_en IS DISTINCT FROM EXCLUDED.name_en OR
		       idxstock.sub_sector.sector_id IS DISTINCT FROM EXCLUDED.sector_id)
		RETURNING id, code, name
	`

	batch := &pgx.Batch{}
	for _, sub := range subSectors {
		batch.Queue(query, sub.Id, sub.Code, sub.Name, sub.NameEn, sub.Description, sub.SectorId)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	updated := make([]models.BasicResponseWithCode, 0)
	for i := 0; i < len(subSectors); i++ {
		var res models.BasicResponseWithCode
		err := br.QueryRow().Scan(&res.Id, &res.Code, &res.Name)
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
