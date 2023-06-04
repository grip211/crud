package commands

import (
	"strconv"
)

// в этом файле будем описывать структуры для более удобной манипуляции с входными и выходными данными

type CreateCommand struct {
	Model    string
	Company  string
	Quantity int
	Price    float32

	CPU         int
	Memory      int
	DisplaySize int
	Camera      int
}

func NewCreteCommand(
	model, company, quantity, price, cpu, memory, display, camera string,
) (*CreateCommand, error) {
	val, err := strconv.ParseFloat(price, 32)
	if err != nil {
		return nil, err
	}
	v, err := strconv.Atoi(quantity)
	if err != nil {
		return nil, err
	}

	vCPU, err := strconv.Atoi(cpu)
	if err != nil {
		return nil, err
	}

	return &CreateCommand{
		Model:    model,
		Company:  company,
		Quantity: v,
		Price:    float32(val),
		CPU:      vCPU,
	}, nil
}

type UpdateCommand struct {
	ID       int
	Model    string
	Company  string
	Quantity int
	Price    float32

	CPU         int
	Memory      int
	DisplaySize int
	Camera      int
}

func NewUpdateCommand(id, model, company, quantity, price string) (*UpdateCommand, error) {
	val, err := strconv.ParseFloat(price, 32)
	if err != nil {
		return nil, err
	}

	v, err := strconv.Atoi(quantity)
	if err != nil {
		return nil, err
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return &UpdateCommand{
		ID:       i,
		Model:    model,
		Company:  company,
		Quantity: int(v),
		Price:    float32(val),
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

type FeatureCommand struct {
	//ID int
	Model   string
	Company string
}

func NewFeatureCommand(model, company string) (*FeatureCommand, error) {
	return &FeatureCommand{
		Model:   model,
		Company: company,
	}, nil
}
