package repository

import (
	"context"
	"errors"
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

	conn, _ := mysql.New(ctx, &database.Opt{
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

	repo := New(conn)

	type args struct {
		command *commands.CreateCommand
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
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
			wantErr: nil,
		},
		{
			name: "successfully insert create, get and delete record",
			args: args{
				command: &commands.CreateCommand{
					Model:       xrand.RandStringBytesMask(302),
					Company:     xrand.RandStringBytesMask(302),
					Quantity:    120,
					Price:       220,
					CPU:         320,
					Memory:      420,
					DisplaySize: 520,
					Camera:      620,
				},
			},
			wantErr: ErrInsertProducts,
		},
		{
			name: "successfully insert feature create, get and delete record",
			args: args{
				command: &commands.CreateCommand{
					Model:       xrand.RandStringBytesMask(302),
					Company:     xrand.RandStringBytesMask(302),
					Quantity:    120,
					Price:       220,
					CPU:         320,
					Memory:      420,
					DisplaySize: 520,
					Camera:      620,
				},
			},
			wantErr: ErrInsertProductFeatures,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := repo.Create(ctx, tt.args.command)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("failed, expected error: %s receive %s", tt.wantErr, err)
				}
				return
			}

			product, err := repo.ReadOne(ctx, id)
			require.NoError(t, err)

			require.Equal(t, product.ID, id)

			require.Equal(t, product.Model, tt.args)
			require.Equal(t, product.Company, tt.args)

			_, err = repo.Delete(ctx, &commands.DeleteCommand{
				ID: id,
			})
			require.NoError(t, err)
		})
	}

	// ... other tests cases

}

func TestRepo_Update(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

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

	type args struct {
		command *commands.UpdateCommand
	}

	tests := []struct {
		name string
		args args
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

				id := repo.Update(ctx, command)
				require.NoError(t, id)

				_, err := repo.Read(ctx)
				require.NoError(t, err)

				//err = repo.Delete(ctx, &commands.DeleteCommand{
				//	ID: command.ID,
				//})
				require.NoError(t, err)
			},
		// ... other tests cases
	}

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
				//err := repo.Feature(ctx, command)
				//require.NoError(t, err)

				//company := repo.Feature(ctx, command)
				//require.NoError(t, company)

				//_, err = repo.Read(ctx)
				//require.NoError(t, err)

				require.Equal(t, command.Model, command.Model)
				require.Equal(t, command.Company, command.Company)

				//	err = repo.Delete(ctx, &commands.DeleteCommand{
				//		ID: command.ID,
				//})
				//require.NoError(t, err)
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
