package grpc

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/nextlag/shortenerURL/internal/configuration"
	"github.com/nextlag/shortenerURL/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/nextlag/shortenerURL/proto"
)

// LinksServer является интерфейсом объединяющим методы grpc.
type LinksServer struct {
	pb.UnimplementedLinksServer
	DB *usecase.UseCase
}

var pgErr *pgconn.PgError

// Get ищет ссылку по короткому адресу и отдает длинную ссылку.
func (s *LinksServer) Get(ctx context.Context, in *pb.ShortenLink) (*pb.ShortenLinkResponse, error) {
	var response pb.ShortenLinkResponse
	longLink, deleteStatus, err := s.DB.DoGet(ctx, in.ShortenLink)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error finding link")
	}
	if deleteStatus {
		return &response, nil
	}

	response.LongLink = longLink
	response.DeleteStatus = deleteStatus

	return &response, nil
}

// Save записывает длинную ссылку в хранилище и отдает короткую ссылку и ошибку.
func (s *LinksServer) Save(ctx context.Context, in *pb.LongLink) (*pb.LongLinkResponse, error) {
	var response pb.LongLinkResponse
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	returnedShortLink, err := s.DB.DoPut(ctx, in.LongLink, "", userID)
	response.ShortenLink = returnedShortLink
	if err != nil {
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return &response, nil
		}
		return &response, status.Errorf(codes.Internal, "error posting link")
	}

	return &response, nil
}

// GetAll получает ID клиента из куки и возвращает все ссылки отправленные им.
func (s *LinksServer) GetAll(ctx context.Context, in *pb.Empty) (*pb.ListShortenLinks, error) {
	cfg, err := configuration.Load()
	if err != nil {
		return nil, err
	}
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	var response pb.ListShortenLinks

	userURLs, err := s.DB.DoGetAll(ctx, userID, cfg.BaseURL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting links")
	}

	response.UserLinks = string(userURLs)
	return &response, nil
}

// Del удаляет ссылки, отправленные клиентом при том условии, что он их загрузил.
func (s *LinksServer) Del(ctx context.Context, in *pb.ListShortenLinksToDelete) (*pb.Empty, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	var result pb.Empty

	go func() {
		for _, alias := range in.UserLinks {
			s.DB.DoDel(userID, []string{alias})
		}
	}()

	return &result, nil
}

func getUserID(ctx context.Context) (int, error) {
	var id int
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return id, status.Errorf(codes.DataLoss, "no metadata")
	}

	values := md.Get("userID")
	if len(values) == 0 {
		return id, status.Errorf(codes.Unauthenticated, "invalid access token")
	}

	id, err := strconv.Atoi(values[0])
	if err != nil {
		return id, status.Errorf(codes.Internal, "can't convert to int")
	}

	return id, nil
}
