package start

import "errors"

type Usecase struct {
}

func (u *Usecase) Execute(query UsecaseQuery) (*UsecaseResponse, error) {
	return nil, errors.New("unsupported usecase")
}
