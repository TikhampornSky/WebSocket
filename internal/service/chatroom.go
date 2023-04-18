package service

import (
	"context"
	"server/internal/domain"
	"server/internal/port"
	"time"
)

type chatroomService struct {
	port.ChatroomRepoPort
	timeout time.Duration
}

func NewChatroomService(repo port.ChatroomRepoPort) port.ChatroomServicePort {
	return &chatroomService{
		repo,
		time.Duration(2) * time.Second,
	}
}

func (s *chatroomService) CreateChatroom(ctx context.Context, req *domain.CreateChatroomReq) (*domain.CreateChatroomRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	c := &domain.Chatroom{
		Name: req.Name,
		Category: req.Category,
	}

	r, err := s.ChatroomRepoPort.CreateChatroom(ctx, c)
	if err != nil {
		return nil, err
	}

	res := &domain.CreateChatroomRes{
		ID:   r.ID,
		Name: r.Name,
		Category: r.Category,
	}

	return res, nil
}

func (s *chatroomService) JoinChatroom(ctx context.Context, req *domain.JoinLeaveChatroomReq) (*domain.JoinLeaveChatroomRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.ChatroomRepoPort.JoinChatroom(ctx, req.ID, req.ClientID)
	if err != nil {
		return nil, err
	}

	return &domain.JoinLeaveChatroomRes{
		ID:   res.ID,
		Name: res.Name,
		Clients: res.Clients,
		Category: res.Category,
	}, nil
}


func (s *chatroomService) LeaveChatroom(ctx context.Context, req *domain.JoinLeaveChatroomReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.ChatroomRepoPort.LeaveChatroom(ctx, req.ID, req.ClientID)
	if err != nil {
		return err
	}

	return nil
}

func (s *chatroomService) GetChatroomByID(ctx context.Context, req *domain.GetChatroomByIDReq) (*domain.GetChatroomByIDRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	r, err := s.ChatroomRepoPort.GetChatroomByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	res := &domain.GetChatroomByIDRes{
		ID:   r.ID,
		Name: r.Name,
		Clients: r.Clients,
		Category: r.Category,
	}
	return res, nil
}

func (s *chatroomService) UpdateChatroomName(ctx context.Context, req *domain.UpdateChatroomNameReq) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.ChatroomRepoPort.UpdateChatroomName(ctx, req.ID, req.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *chatroomService) GetAllChatrooms(ctx context.Context, userID int64) ([]*domain.Chatroom, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	r, err := s.ChatroomRepoPort.GetAllChatrooms(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := []*domain.Chatroom{}
	for _, c := range r {
		res = append(res, &domain.Chatroom{
			ID:      c.ID,
			Name:    c.Name,
			Clients: c.Clients,
			Category: c.Category,
		})
	}

	return res, nil
}