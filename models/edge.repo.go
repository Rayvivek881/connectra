package models

import (
	"context"

	"github.com/uptrace/bun"
)

type EdgesStruct struct {
	PgDbClient *bun.DB
}

func EdgesRepository(db *bun.DB) EdgesSvcRepo {
	return &EdgesStruct{
		PgDbClient: db,
	}
}

type EdgesSvcRepo interface {
	CreateEdges(edges []*ModelEdges) error
	UpdateNodeDegree(source string) error
}

func (e *EdgesStruct) CreateEdges(edges []*ModelEdges) error {
	_, err := e.PgDbClient.NewInsert().
		Model(&edges).
		Exec(context.Background())

	return err
}

func (e *EdgesStruct) UpdateNodeDegree(source string) error {
	_, err := e.PgDbClient.NewUpdate().
		Model((*ModelJobNodes)(nil)).
		Set("degree = degree - 1").
		Where("uuid IN (?)",
			e.PgDbClient.NewSelect().
				Model((*ModelEdges)(nil)).
				Column("target").
				Where("source = ?", source),
		).
		Exec(context.Background())

	return err
}
