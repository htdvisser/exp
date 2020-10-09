package router

import (
	"context"
	"sort"

	"htdvisser.dev/exp/rjs/client"
)

func NewRouterService(client *client.Client) *Service {
	return &Service{client: client}
}

type Service struct {
	client *client.Client
}

func (s *Service) Claim(ctx context.Context, req *ClaimRequest) (ClaimResponse, error) {
	var res ClaimResponse
	if err := s.client.Call(ctx, ClaimURI, req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) ClaimBulk(ctx context.Context, req *ClaimBulkRequest) (ClaimResponse, error) {
	var res ClaimResponse
	if err := s.client.Call(ctx, ClaimURI, req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) Add(ctx context.Context, req *AddRequest) (AddResponse, error) {
	var res AddResponse
	if err := s.client.Call(ctx, AddURI, req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) AddBulk(ctx context.Context, req *AddBulkRequest) (AddResponse, error) {
	var res AddResponse
	if err := s.client.Call(ctx, AddURI, req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) Delete(ctx context.Context, req *DeleteRequest) (DeleteResponse, error) {
	var res DeleteResponse
	if err := s.client.Call(ctx, DeleteURI, req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) DeleteBulk(ctx context.Context, req *DeleteBulkRequest) (DeleteResponse, error) {
	var res DeleteResponse
	if err := s.client.Call(ctx, DeleteURI, req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) Setup(ctx context.Context, req *SetupRequest) (SetupResponse, error) {
	var res SetupResponse
	if err := s.client.Call(ctx, SetupURI, req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) SetupBulk(ctx context.Context, req *SetupBulkRequest) (SetupResponse, error) {
	var res SetupResponse
	if err := s.client.Call(ctx, SetupURI, req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) Info(ctx context.Context, req *InfoRequest) (InfoResponse, error) {
	var res InfoResponse
	if err := s.client.Call(ctx, InfoURI, req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) InfoBulk(ctx context.Context, req *InfoBulkRequest) (InfoResponse, error) {
	var res InfoResponse
	if err := s.client.Call(ctx, InfoURI, req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) List(ctx context.Context, req *ListRequest) (ListResponse, error) {
	var res ListResponse
	if err := s.client.Call(ctx, ListURI, req, &res); err != nil {
		return nil, err
	}
	sort.Sort(res)
	return res, nil
}

func (s *Service) Token(ctx context.Context, req *TokenRequest) (TokenResponse, error) {
	var res TokenResponse
	if err := s.client.Call(ctx, TokenURI, req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Service) TokenBulk(ctx context.Context, req *TokenBulkRequest) (TokenResponse, error) {
	var res TokenResponse
	if err := s.client.Call(ctx, TokenURI, req, &res); err != nil {
		return nil, err
	}
	return res, nil
}
