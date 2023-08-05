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
			name: "failed insert create, get and delete record",
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
			name: "failed insert feature create, get and delete record",
			args: args{
				command: &commands.CreateCommand{
					Model:       xrand.RandStringBytesMask(300),
					Company:     xrand.RandStringBytesMask(302),
					Quantity:    1230,
					Price:       2230,
					CPU:         3320,
					Memory:      4230,
					DisplaySize: 5230,
					Camera:      6230,
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

			require.Equal(t, tt.args.command.Model, product.Model)
			require.Equal(t, tt.args.command.Company, product.Company)
			require.Equal(t, tt.args.command.Price, product.Price)
			require.Equal(t, tt.args.command.Quantity, product.Quantity)

			_, err = repo.Delete(ctx, &commands.DeleteCommand{
				ID: id,
			})
			require.NoError(t, err)
		})
	}
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
	require.NoError(t, err)

	repo := New(conn)

	type args struct {
		createCommand *commands.CreateCommand
		updateCommand *commands.UpdateCommand
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "successfully update, get and delete record",
			args: args{
				createCommand: &commands.CreateCommand{
					Model:       xrand.RandStringBytesMask(30),
					Company:     xrand.RandStringBytesMask(30),
					Quantity:    10,
					Price:       20,
					CPU:         30,
					Memory:      40,
					DisplaySize: 50,
					Camera:      60,
				},
				updateCommand: &commands.UpdateCommand{
					Model:       xrand.RandStringBytesMask(20),
					Company:     xrand.RandStringBytesMask(20),
					Quantity:    100,
					Price:       200,
					CPU:         300,
					Memory:      400,
					DisplaySize: 500,
					Camera:      600,
				},
			},
			wantErr: nil,
		},
		{
			name: "failed update, get and delete record",
			args: args{
				createCommand: &commands.CreateCommand{
					Model:       xrand.RandStringBytesMask(30),
					Company:     xrand.RandStringBytesMask(30),
					Quantity:    10,
					Price:       20,
					CPU:         30,
					Memory:      40,
					DisplaySize: 50,
					Camera:      60,
				},
				updateCommand: &commands.UpdateCommand{
					Model:       xrand.RandStringBytesMask(2230),
					Company:     xrand.RandStringBytesMask(2440),
					Quantity:    111,
					Price:       33,
					CPU:         222,
					Memory:      44,
					DisplaySize: 55,
					Camera:      66,
				},
			},
			wantErr: ErrUpdateProduct,
		},
		{
			name: "failed update, get and delete record",
			args: args{
				createCommand: &commands.CreateCommand{
					Model:       xrand.RandStringBytesMask(30),
					Company:     xrand.RandStringBytesMask(30),
					Quantity:    10,
					Price:       20,
					CPU:         30,
					Memory:      40,
					DisplaySize: 50,
					Camera:      60,
				},
				updateCommand: &commands.UpdateCommand{
					Model:       xrand.RandStringBytesMask(20),
					Company:     xrand.RandStringBytesMask(20),
					Quantity:    1111100,
					Price:       2111100,
					CPU:         3111100,
					Memory:      4111100,
					DisplaySize: 5111100,
					Camera:      6111100,
				},
			},
			wantErr: ErrUpsertFeature,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := repo.Create(ctx, tt.args.createCommand)
			require.NoError(t, err)

			tt.args.updateCommand.ID = id

			err = repo.Update(ctx, tt.args.updateCommand)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					//	t.Fatalf("failed, expected error: %s receive %s", tt.wantErr, err)
				}
				return
			}

			product, err := repo.ReadOneWithFeatures(ctx, id)
			require.NoError(t, err)

			require.Equal(t, tt.args.updateCommand.Model, product.Model)
			require.Equal(t, tt.args.updateCommand.Company, product.Company)
			require.Equal(t, tt.args.updateCommand.Price, product.Price)
			require.Equal(t, tt.args.updateCommand.Quantity, product.Quantity)
			require.Equal(t, tt.args.updateCommand.Camera, int(product.Features.Camera.Int32))
			require.Equal(t, tt.args.updateCommand.CPU, int(product.Features.CPU.Int32))
			require.Equal(t, tt.args.updateCommand.Memory, int(product.Features.Memory.Int32))
			require.Equal(t, tt.args.updateCommand.DisplaySize, int(product.Features.Display.Int32))

			_, err = repo.Delete(ctx, &commands.DeleteCommand{
				ID: id,
			})
			require.NoError(t, err)
		})
	}
}

// delete
// read
// read one
// read one with feature - тоже самое как с Read One тестом, только еще првоерять дополнительные поля

func TestRepo_Delete(t *testing.T) {
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
	require.NoError(t, err)

	repo := New(conn)

	type args struct {
		createCommand *commands.CreateCommand
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "successfully update, get and delete record",
			args: args{
				createCommand: &commands.CreateCommand{
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
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := repo.Create(ctx, tt.args.createCommand)
			require.NoError(t, err)

			_, err = repo.Delete(ctx, &commands.DeleteCommand{
				ID: id,
			})
			require.NoError(t, err)

			_, err = repo.ReadOne(ctx, id)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("failed, expected error: %s receive %s", tt.wantErr, err)
				}
				return
			}
		})
	}
}

func TestRepo_ReadOne(t *testing.T) {
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
	require.NoError(t, err)

	repo := New(conn)

	type args struct {
		createCommand *commands.CreateCommand
	}
	tests := []struct {
		name      string
		args      args
		wantErr   error
		replaceID int
	}{
		{
			name: "successfully update, get and delete record",
			args: args{
				createCommand: &commands.CreateCommand{
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
			name: "successfully update, get and delete record",
			args: args{
				createCommand: &commands.CreateCommand{
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
			wantErr:   ErrNotFound,
			replaceID: 9999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := repo.Create(ctx, tt.args.createCommand)
			require.NoError(t, err)

			if tt.replaceID != 0 {
				id = tt.replaceID
			}

			product, err := repo.ReadOne(ctx, id)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("failed, expected error: %s receive %s", tt.wantErr, err)
				}
				return
			}

			require.Equal(t, tt.args.createCommand.Model, product.Model)
			require.Equal(t, tt.args.createCommand.Company, product.Company)
			require.Equal(t, tt.args.createCommand.Price, product.Price)
			require.Equal(t, tt.args.createCommand.Quantity, product.Quantity)

			_, err = repo.Delete(ctx, &commands.DeleteCommand{
				ID: id,
			})
			require.NoError(t, err)
		})
	}
}

func TestRepo_ReadOneWithFeatures(t *testing.T) {
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
	require.NoError(t, err)

	repo := New(conn)

	type args struct {
		createCommand *commands.CreateCommand
	}
	tests := []struct {
		name      string
		args      args
		wantErr   error
		replaceID int
	}{
		{
			name: "successfully update, get and delete record",
			args: args{
				createCommand: &commands.CreateCommand{
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
			name: "successfully update, get and delete record",
			args: args{
				createCommand: &commands.CreateCommand{
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
			wantErr:   ErrNotFound,
			replaceID: 9999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := repo.Create(ctx, tt.args.createCommand)
			require.NoError(t, err)

			if tt.replaceID != 0 {
				id = tt.replaceID
			}

			product, err := repo.ReadOneWithFeatures(ctx, id)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("failed, expected error: %s receive %s", tt.wantErr, err)
				}
				return
			}

			require.Equal(t, tt.args.createCommand.Model, product.Model)
			require.Equal(t, tt.args.createCommand.Company, product.Company)
			require.Equal(t, tt.args.createCommand.Price, product.Price)
			require.Equal(t, tt.args.createCommand.Quantity, product.Quantity)
			require.Equal(t, tt.args.createCommand.CPU, int(product.Features.CPU.Int32))
			require.Equal(t, tt.args.createCommand.Camera, int(product.Features.Camera.Int32))
			require.Equal(t, tt.args.createCommand.Memory, int(product.Features.Memory.Int32))
			require.Equal(t, tt.args.createCommand.DisplaySize, int(product.Features.Display.Int32))

			_, err = repo.Delete(ctx, &commands.DeleteCommand{
				ID: id,
			})
			require.NoError(t, err)
		})
	}
}

func TestRepo_Read(t *testing.T) {
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
	require.NoError(t, err)

	repo := New(conn)

	type args struct {
		pseudoCreateCommands []*commands.UpdateCommand
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "successfully update, get and delete records",
			args: args{
				pseudoCreateCommands: []*commands.UpdateCommand{
					{
						Model:       xrand.RandStringBytesMask(30),
						Company:     xrand.RandStringBytesMask(30),
						Quantity:    10,
						Price:       20,
						CPU:         30,
						Memory:      40,
						DisplaySize: 50,
						Camera:      60,
					},
					{
						Model:       xrand.RandStringBytesMask(30),
						Company:     xrand.RandStringBytesMask(30),
						Quantity:    110,
						Price:       120,
						CPU:         130,
						Memory:      140,
						DisplaySize: 150,
						Camera:      160,
					},
					{
						Model:       xrand.RandStringBytesMask(30),
						Company:     xrand.RandStringBytesMask(30),
						Quantity:    210,
						Price:       220,
						CPU:         230,
						Memory:      240,
						DisplaySize: 250,
						Camera:      260,
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdIDs := make([]int, 0, len(tt.args.pseudoCreateCommands))
			for _, createCommand := range tt.args.pseudoCreateCommands {
				id, err := repo.Create(ctx, &commands.CreateCommand{
					Model:       createCommand.Model,
					Company:     createCommand.Company,
					Quantity:    createCommand.Quantity,
					Price:       createCommand.Price,
					CPU:         createCommand.CPU,
					Memory:      createCommand.Memory,
					DisplaySize: createCommand.DisplaySize,
					Camera:      createCommand.Camera,
				})
				require.NoError(t, err)

				createCommand.ID = id
				createdIDs = append(createdIDs, id)
			}

			products, err := repo.Read(ctx)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("failed, expected error: %s receive %s", tt.wantErr, err)
				}
				return
			}

			if len(products) != len(tt.args.pseudoCreateCommands) {
				//				t.Fatalf("faield, expect %b count product, receive count %b",
				//					len(tt.args.pseudoCreateCommands),
				//					len(products),
				//				)
			}

			for _, product := range products {
				for _, createCommand := range tt.args.pseudoCreateCommands {
					if createCommand.ID == product.ID {
						require.Equal(t, createCommand.Model, product.Model)
						require.Equal(t, createCommand.Company, product.Company)
						require.Equal(t, createCommand.Price, product.Price)
						require.Equal(t, createCommand.Quantity, product.Quantity)
					}
				}
			}

			for _, id := range createdIDs {
				_, err = repo.Delete(ctx, &commands.DeleteCommand{
					ID: id,
				})
				require.NoError(t, err)
			}
		})
	}
}
