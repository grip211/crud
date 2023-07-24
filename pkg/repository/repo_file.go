package repository

import (
	"context"

	"github.com/grip211/crud/pkg/commands"
)

type FileRepo struct {
	fileName string
}

func (f FileRepo) Create(ctx context.Context, command *commands.CreateCommand) (int, error) {
	//  TODO implement me

	panic("implement me")
	// err = file.read(fileName)
}
