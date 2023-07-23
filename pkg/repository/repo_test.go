package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/grip211/crud/pkg/commands"
	"github.com/grip211/crud/pkg/database"
	"github.com/grip211/crud/pkg/database/mysql"
	"github.com/grip211/crud/pkg/xrand"
)

func TestRepo_Create(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	type args struct {
		command *commands.CreateCommand
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, ctx context.Context, repo *Repo, command *commands.CreateCommand)
	}{
		{
			name: "successfully create, get and delete record",
			args: args{
				command: &commands.CreateCommand{
					Model:       xrand.RandStringBytesMask(30),
					Company:     xrand.RandStringBytesMask(30),
					Quantity:    10,
					Price:       20,
					CPU:         30,
					Memory:      40,
					DisplaySize: 50,
					Camera:      60,
				},
			},
			check: func(t *testing.T, ctx context.Context, repo *Repo, command *commands.CreateCommand) {
				id, err := repo.Create(ctx, command)
				require.NoError(t, err)

				product, err := repo.ReadOne(ctx, id)
				require.NoError(t, err)

				require.Equal(t, product.ID, id)
				require.Equal(t, product.Model, command.Model)
				require.Equal(t, product.Company, command.Company)

				err = repo.Delete(ctx, &commands.DeleteCommand{
					ID: id,
				})
				require.NoError(t, err)
			},
		},

		// ... other tests cases
	}

	conn, err := mysql.New(ctx, &database.Opt{
		Host:               os.Getenv("DB_Host"),
		User:               os.Getenv("DB_USER"),
		Password:           os.Getenv("DB_PASS"),
		Name:               os.Getenv("DB_NAME"),
		Dialect:            "mysql",
		MaxConnMaxLifetime: time.Minute * 5,
		MaxOpenConns:       10,
		MaxIdleConns:       9,
		Debug:              true,
	})
	require.NoError(t, err)

	repo := &Repo{
		db: conn,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, ctx, repo, tt.args.command)
		})
	}
}

func TestRepo_Update(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	type args struct {
		command *commands.UpdateCommand
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, ctx context.Context, repo *Repo, command *commands.UpdateCommand)
	}{
		{
			name: "successfully update, get and delete record",
			args: args{
				command: &commands.UpdateCommand{
					Model:       xrand.RandStringBytesMask(30),
					Company:     xrand.RandStringBytesMask(30),
					Quantity:    10,
					Price:       20,
					CPU:         30,
					Memory:      40,
					DisplaySize: 50,
					Camera:      60,
				},
			},
			check: func(t *testing.T, ctx context.Context, repo *Repo, command *commands.UpdateCommand) {
				id := repo.Update(ctx, command)
				require.NoError(t, id)

				_, err := repo.Read(ctx)
				require.NoError(t, err)

				err = repo.Delete(ctx, &commands.DeleteCommand{
					ID: command.ID,
				})
				require.NoError(t, err)
			},
		},
		// ... other tests cases
	}

	conn, err := mysql.New(ctx, &database.Opt{
		Host:               os.Getenv("DB_Host"),
		User:               os.Getenv("DB_USER"),
		Password:           os.Getenv("DB_PASS"),
		Name:               os.Getenv("DB_NAME"),
		Dialect:            "mysql",
		MaxConnMaxLifetime: time.Minute * 5,
		MaxOpenConns:       10,
		MaxIdleConns:       9,
		Debug:              true,
	})
	require.NoError(t, err)

	repo := &Repo{
		db: conn,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, ctx, repo, tt.args.command)
		})
	}
}

func TestRepo_Feature(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	type args struct {
		command *commands.FeatureCommand
	}
	tests := []struct {
		name  string
		args  args
		check func(t *testing.T, ctx context.Context, repo *Repo, command *commands.FeatureCommand)
	}{
		{
			name: "successfully feature, get and delete record",
			args: args{
				command: &commands.FeatureCommand{
					Model:   xrand.RandStringBytesMask(30),
					Company: xrand.RandStringBytesMask(30),
				},
			},
			check: func(t *testing.T, ctx context.Context, repo *Repo, command *commands.FeatureCommand) {
				err := repo.Feature(ctx, command)
				//require.NoError(t, err)

				//company := repo.Feature(ctx, command)
				//require.NoError(t, company)

				_, err = repo.Read(ctx)
				require.NoError(t, err)

				require.Equal(t, command.Model, command.Model)
				require.Equal(t, command.Company, command.Company)

				err = repo.Delete(ctx, &commands.DeleteCommand{
					ID: command.ID,
				})
				require.NoError(t, err)
			},
		},
		// ... other tests cases
	}

	conn, err := mysql.New(ctx, &database.Opt{
		Host:               os.Getenv("DB_Host"),
		User:               os.Getenv("DB_USER"),
		Password:           os.Getenv("DB_PASS"),
		Name:               os.Getenv("DB_NAME"),
		Dialect:            "mysql",
		MaxConnMaxLifetime: time.Minute * 5,
		MaxOpenConns:       10,
		MaxIdleConns:       9,
		Debug:              true,
	})
	require.NoError(t, err)

	repo := &Repo{
		db: conn,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, ctx, repo, tt.args.command)
		})
	}
}
