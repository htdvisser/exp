package signature

import (
	"context"

	"htdvisser.dev/exp/rjs"
)

func NewSignatureService(client rjs.Client) *Service {
	return &Service{client: client}
}

type Service struct {
	client rjs.Client
}

func (s *Service) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
	var res AddResponse
	if err := s.client.Call(ctx, AddURI, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Service) Delete(ctx context.Context, req *DeleteRequest) (*DeleteResponse, error) {
	var res DeleteResponse
	if err := s.client.Call(ctx, DeleteURI, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Service) List(ctx context.Context, req *ListRequest) (ListResponse, error) {
	var res ListResponse
	if err := s.client.Call(ctx, ListURI, req, &res); err != nil {
		return nil, err
	}
	// TODO: sort.Sort(res)
	return res, nil
}
