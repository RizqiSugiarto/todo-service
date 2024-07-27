package handler

import (
	"context"

	"github.com/digisata/todo-service/internal/entity"
	textPB "github.com/digisata/todo-service/stubs/text"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TextHandler struct {
	textPB.UnimplementedTextServiceServer
	textUseCase TextUseCase
}

func NewText(textUseCase TextUseCase) *TextHandler {
	return &TextHandler{
		textUseCase: textUseCase,
	}
}

func (h *TextHandler) Create(ctx context.Context, req *textPB.CreateTextRequest) (*textPB.TextBaseResponse, error) {
	payload := entity.CreateTextRequest{
		ActivityID: req.GetActivityId(),
		Text:       req.Text,
	}

	err := h.textUseCase.CreateText(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &textPB.TextBaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (g *TextHandler) Update(ctx context.Context, req *textPB.UpdateTextByIDRequest) (*textPB.TextBaseResponse, error) {
	payload := entity.UpdateTextRequest{
		ID:   req.GetId(),
		Text: req.Text,
	}

	err := g.textUseCase.UpdateText(ctx, payload)
	if err != nil && err.Error() == "data not found" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &textPB.TextBaseResponse{
		Message: "Success",
	}

	return res, nil
}

func (h *TextHandler) Get(ctx context.Context, req *textPB.GetTextByIDRequest) (*textPB.GetTextByIDResponse, error) {
	data, err := h.textUseCase.GetText(ctx, req.GetId())
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &textPB.GetTextByIDResponse{
		Id:         data.ID,
		ActivityId: data.ActivityID,
		Text:       data.Text,
		CreatedAt:  timestamppb.New(data.CreatedAt),
		UpdatedAt:  timestamppb.New(data.UpdatedAt),
	}

	return res, nil
}

func (g *TextHandler) GetAllByUserID(ctx context.Context, req *textPB.GetAllTextByActivityIDRequest) (*textPB.GetAllTextByActivityIDResponse, error) {
	payload := entity.GetAllTextRequest{
		ActivityID:   req.GetActivityId(),
		Search:       req.Search,
		Page:         req.Page,
		Limit:        req.Limit,
		IsNewest:     req.IsNewest,
		IsOldest:     req.IsOldest,
		IsAscending:  req.IsAscending,
		IsDescending: req.IsDescending,
	}

	data, paging, err := g.textUseCase.GetAllTextByActivityID(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &textPB.GetAllTextByActivityIDResponse{
		Message: "Success",
		Texts:   []*textPB.GetTextByIDResponse{},
		Paging: &textPB.TextPaging{
			CurrentPage: paging.CurrentPage,
			TotalPage:   paging.TotalPage,
			Count:       paging.Count,
		},
	}
	for _, text := range data {
		data := &textPB.GetTextByIDResponse{
			Id:         text.ID,
			ActivityId: text.ActivityID,
			Text:       text.Text,
			CreatedAt:  timestamppb.New(text.CreatedAt),
			UpdatedAt:  timestamppb.New(text.UpdatedAt),
		}

		res.Texts = append(res.Texts, data)
	}

	return res, nil
}

func (g *TextHandler) Delete(ctx context.Context, req *textPB.DeleteTextByIDRequest) (*textPB.TextBaseResponse, error) {
	err := g.textUseCase.DeleteText(ctx, req.GetId())
	if err != nil && err.Error() == "data not found" {
		return nil, status.Errorf(codes.NotFound, "data for userId: %v", req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server error: %v", err)
	}

	res := &textPB.TextBaseResponse{
		Message: "Success",
	}

	return res, nil
}
