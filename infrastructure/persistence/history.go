package persistence

import (
	"context"
	"database/sql"

	"github.com/nao1215/sqluv/config"
	"github.com/nao1215/sqluv/domain/model"
	"github.com/nao1215/sqluv/domain/repository"
	"github.com/nao1215/sqluv/infrastructure"
)

// _ interface implementation check
var _ repository.HistoryTableCreator = (*historyTableCreator)(nil)

type historyTableCreator struct {
	db *sql.DB
}

// NewHistoryTableCreator return new HistoryTableCreator.
func NewHistoryTableCreator(db config.HistoryDB) repository.HistoryTableCreator {
	return &historyTableCreator{
		db: db,
	}
}

// CreateTable create a DB table for sqluv shell history
func (h *historyTableCreator) CreateTable(ctx context.Context) error {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := "CREATE TABLE IF NOT EXISTS `history` (id INTEGER PRIMARY KEY AUTOINCREMENT, request TEXT)"
	_, err = tx.ExecContext(ctx, q)
	if err != nil {
		return err
	}

	q = "CREATE INDEX IF NOT EXISTS `history_id_index` ON `history`(`id`)"
	_, err = tx.ExecContext(ctx, q)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// _ interface implementation check
var _ repository.HistoryCreator = (*historyCreator)(nil)

type historyCreator struct {
	db *sql.DB
}

// NewHistoryCreator return new HistoryCreator.
func NewHistoryCreator(db config.HistoryDB) repository.HistoryCreator {
	return &historyCreator{
		db: db,
	}
}

// Create set history record in DB
func (h *historyCreator) Create(ctx context.Context, t *model.Table) error {
	if err := t.Valid(); err != nil {
		return err
	}

	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, v := range t.Records() {
		if _, err := tx.ExecContext(ctx, infrastructure.GenerateInsertStatement(t.Name(), v)); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// _ interface implementation check
var _ repository.HistoryLister = (*historyLister)(nil)

type historyLister struct {
	db *sql.DB
}

// NewHistoryLister return new HistoryLister.
func NewHistoryLister(db config.HistoryDB) repository.HistoryLister {
	return &historyLister{
		db: db,
	}
}

// List get sql shell all history.
func (h *historyLister) List(ctx context.Context) (model.Histories, error) {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx,
		"SELECT `id`, `request` FROM `history` ORDER BY `id` ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id int
	var request string
	histories := model.Histories{}
	for rows.Next() {
		if err := rows.Scan(&id, &request); err != nil {
			return nil, err
		}
		histories = append(histories, model.NewHistory(id, request))
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return histories, nil
}
