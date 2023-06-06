package repository

import (
	"context"
	"errors"
	builder "github.com/doug-martin/goqu/v9"
	"github.com/grip211/crud/pkg/models"

	"github.com/grip211/crud/pkg/commands"
	"github.com/grip211/crud/pkg/database"
)

// Repo пишем структуру которая будет реализовывать все нужные нам методы для работы с базой данных
// в нашем случае Create, Read, Update, Delete

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
	result, err := r.db.Builder().
		Insert("productdb.Products").
		Rows(builder.Record{
			"model":    command.Model,
			"company":  command.Company,
			"quantity": command.Quantity,
			"price":    command.Price,
		}).
		Executor().
		ExecContext(ctx)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	_, err = r.db.Builder().
		Insert("productdb.ProductsFeatures").
		Rows(builder.Record{
			"product_id":   id,
			"cpu":          command.CPU,
			"memory":       command.Memory,
			"display_size": command.DisplaySize,
			"camera":       command.Camera,
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
			builder.C("quantity"),
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
			builder.C("quantity"),
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

func (r *Repo) ReadOneWithFeatures(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product
	found, err := r.db.Builder().
		Select(
			builder.I("Products.id").As("id"),
			builder.C("company"),
			builder.C("model"),
			builder.C("quantity"),
			builder.C("price"),
			builder.I("ProductsFeatures.cpu").As(builder.C("features.cpu")),
			builder.I("ProductsFeatures.memory").As(builder.C("features.memory")),
			builder.I("ProductsFeatures.display_size").As(builder.C("features.display")),
			builder.I("ProductsFeatures.camera").As(builder.C("features.camera")),
		).
		From("productdb.Products").
		LeftJoin(
			builder.T("ProductsFeatures"),
			builder.On(builder.Ex{
				"Products.id": builder.I("ProductsFeatures.product_id")}),
		).
		Where(
			builder.I("Products.id").Eq(id),
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
			"model":    command.Model,
			"company":  command.Company,
			"quantity": command.Quantity,
			"price":    command.Price,
		}).
		Where(
			builder.C("id").Eq(command.ID),
		).
		Executor().
		ExecContext(ctx)
	if err != nil {
		return err
	}

	_, err = r.db.Builder().
		Insert("productdb.ProductsFeatures").
		Rows(builder.Record{
			"product_id":   command.ID,
			"cpu":          command.CPU,
			"memory":       command.Memory,
			"display_size": command.DisplaySize,
			"camera":       command.Camera,
		}).
		OnConflict(builder.DoUpdate("key", builder.Record{
			"cpu":          command.CPU,
			"memory":       command.Memory,
			"display_size": command.DisplaySize,
			"camera":       command.Camera,
		})).
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

func (r *Repo) Feature(ctx context.Context, command *commands.FeatureCommand) error {
	_, err := r.db.Builder().
		Insert("productdb.Products").
		Rows(builder.Record{
			"model":   command.Model,
			"company": command.Company,
		}).
		Executor().
		ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
