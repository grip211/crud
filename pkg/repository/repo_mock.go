package repository

import (
	"context"

	"github.com/grip211/crud/pkg/commands"
)

type mockRepo struct{}

func (m mockRepo) Create(ctx context.Context, command *commands.CreateCommand) (int, error) {
	return 0, nil
}

func (m mockRepo) Update(ctx context.Context, commands *commands.UpdateCommand) (int, error) {
	return 0, nil
}

func (m mockRepo) Delete(ctx context.Context, command commands.DeleteCommand) (int, error) {
	return 0, nil
}

func (m mockRepo) Feature(ctx context.Context, command commands.FeatureCommand) (int, error) {
	return 0, nil
}
