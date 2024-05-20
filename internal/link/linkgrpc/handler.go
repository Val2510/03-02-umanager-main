package linkgrpc

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/internal/database"
	"gitlab.com/robotomize/gb-golang/homework/03-02-umanager/pkg/pb"
)

var _ pb.LinkServiceServer = (*Handler)(nil)

func New(linksRepository linksRepository, timeout time.Duration) *Handler {
	return &Handler{linksRepository: linksRepository, timeout: timeout}
}

type Handler struct {
	pb.UnimplementedLinkServiceServer
	linksRepository linksRepository
	timeout         time.Duration
}

func (h Handler) GetLinkByUserID(ctx context.Context, req *pb.GetLinksByUserId) (*pb.ListLinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	links, err := h.linksRepository.FindByUserID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find links by user ID: %v", err)
	}

	var response pb.ListLinkResponse
	for _, link := range links {
		response.Links = append(response.Links, &pb.Link{
			Id:     link.ID.Hex(),
			Url:    link.URL,
			Title:  link.Title,
			Tags:   link.Tags,
			Images: link.Images,
			UserId: link.UserID,
		})
	}

	return &response, nil
}

func (h Handler) CreateLink(ctx context.Context, req *pb.CreateLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	linkID := primitive.NewObjectID()

	newLink := database.CreateLinkReq{
		ID:     linkID,
		URL:    req.Url,
		Title:  req.Title,
		Tags:   req.Tags,
		Images: req.Images,
		UserID: req.UserId,
	}

	_, err := h.linksRepository.Create(ctx, newLink)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create link: %v", err)
	}

	return &pb.Empty{}, nil
}

func (h Handler) GetLink(ctx context.Context, req *pb.GetLinkRequest) (*pb.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid link ID: %v", err)
	}

	link, err := h.linksRepository.FindByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get link: %v", err)
	}

	return &pb.Link{
		Id:     link.ID.Hex(),
		Url:    link.URL,
		Title:  link.Title,
		Tags:   link.Tags,
		Images: link.Images,
		UserId: link.UserID,
	}, nil
}

func (h Handler) UpdateLink(ctx context.Context, req *pb.UpdateLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid link ID: %v", err)
	}

	updateReq := database.UpdateLinkReq{
		ID:     id,
		URL:    req.Url,
		Title:  req.Title,
		Tags:   req.Tags,
		Images: req.Images,
		UserID: req.UserId,
	}

	_, err = h.linksRepository.Update(ctx, updateReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update link: %v", err)
	}

	return &pb.Empty{}, nil
}

func (h Handler) DeleteLink(ctx context.Context, req *pb.DeleteLinkRequest) (*pb.Empty, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid link ID: %v", err)
	}

	err = h.linksRepository.Delete(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete link: %v", err)
	}

	return &pb.Empty{}, nil
}

func (h Handler) ListLinks(ctx context.Context, req *pb.Empty) (*pb.ListLinkResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	links, err := h.linksRepository.FindAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list links: %v", err)
	}

	var response pb.ListLinkResponse
	for _, link := range links {
		response.Links = append(response.Links, &pb.Link{
			Id:     link.ID.Hex(),
			Url:    link.URL,
			Title:  link.Title,
			Tags:   link.Tags,
			Images: link.Images,
			UserId: link.UserID,
		})
	}

	return &response, nil
}
