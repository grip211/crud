package commands

import "strconv"

type CreateCommand struct {
	Model   string
	Company string
	Price   float32
}

func NewCreteCommand(model, company, price string) (*CreateCommand, error) {
	val, err := strconv.ParseFloat(price, 32)
	if err != nil {
		return nil, err
	}
	return &CreateCommand{
		Model:   model,
		Company: company,
		Price:   float32(val),
	}, nil
}

type UpdateCommand struct {
	ID      int
	Model   string
	Company string
	Price   float32
}

func NewUpdateCommand(id, model, company, price string) (*UpdateCommand, error) {
	val, err := strconv.ParseFloat(price, 32)
	if err != nil {
		return nil, err
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return &UpdateCommand{
		ID:      i,
		Model:   model,
		Company: company,
		Price:   float32(val),
	}, nil
}

type DeleteCommand struct {
	ID int
}

func NewDeleteCommand(id string) (*DeleteCommand, error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return &DeleteCommand{
		ID: i,
	}, nil
}
