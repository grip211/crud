package repository

import (
	"context"
	"errors"
	builder "github.com/doug-martin/goqu/v9"
	"github.com/grip211/crud/pkg/models"

	"github.com/grip211/crud/pkg/commands"
	"github.com/grip211/crud/pkg/database"
)

var ErrNotFound = errors.New("not found")

type Repo struct {
	db database.Pool
}

func New(db database.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Create(ctx context.Context, command *commands.CreateCommand) error {
	_, err := r.db.Builder().
		Insert("productdb.Products").
		Rows(builder.Record{
			"model":   command.Model,
			"company": command.Company,
			"price":   command.Price,
		}).
		Executor().
		ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) Read(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Builder().
		Select(
			builder.C("id"),
			builder.C("company"),
			builder.C("model"),
			builder.C("price"),
		).
		From("productdb.Products").
		ScanStructsContext(ctx, &products)

	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *Repo) ReadOne(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product
	found, err := r.db.Builder().
		Select(
			builder.C("id"),
			builder.C("company"),
			builder.C("model"),
			builder.C("price"),
		).
		From("productdb.Products").
		Where(
			builder.C("id").Eq(id),
		).
		ScanStructContext(ctx, &product)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrNotFound
	}
	return &product, nil
}

func (r *Repo) Update(ctx context.Context, command *commands.UpdateCommand) error {
	_, err := r.db.Builder().
		Update("productdb.Products").
		Set(builder.Record{
			"model":   command.Model,
			"company": command.Company,
			"price":   command.Price,
		}).
		Where(builder.C("id").Eq(command.ID)).
		Executor().
		ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) Delete(ctx context.Context, command *commands.DeleteCommand) error {
	_, err := r.db.Builder().
		Delete("productdb.Products").
		Where(builder.C("id").Eq(command.ID)).
		Executor().
		ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
