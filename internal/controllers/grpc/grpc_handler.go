package grpc

import (
	"context"
	"errors"
	"fmt"
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

// LinksServer is a gRPC server that implements the Links service.
type LinksServer struct {
	pb.UnimplementedLinksServer
	DB *usecase.UseCase
}

var pgErr *pgconn.PgError

// Get retrieves a long link by its short link.
func (s *LinksServer) Get(ctx context.Context, in *pb.ShortenLink) (*pb.ShortenLinkResponse, error) {
	var response pb.ShortenLinkResponse
	url, err := s.DB.DoGet(ctx, in.ShortenLink)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error finding link")
	}

	if url == nil {
		return nil, status.Errorf(codes.NotFound, "Link not found")
	}

	response.LongLink = url.URL
	response.DeleteStatus = url.IsDeleted

	return &response, nil
}

// Save stores a long link and returns the corresponding short link.
func (s *LinksServer) Save(ctx context.Context, in *pb.LongLink) (*pb.LongLinkResponse, error) {
	var response pb.LongLinkResponse
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	shortLink, err := s.DB.DoPut(ctx, in.LongLink, "", userID)
	if err != nil {
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return &response, nil
		}
		return nil, status.Errorf(codes.Internal, "Error saving link")
	}

	response.ShortenLink = shortLink
	return &response, nil
}

// GetAll retrieves all links for a user.
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
	urls, err := s.DB.DoGetAll(ctx, userID, cfg.BaseURL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting links")
	}

	for _, url := range urls {
		response.UserLinks = append(response.UserLinks, &pb.UserLink{
			LongLink:  url.URL,
			ShortLink: url.Alias,
		})
	}

	return &response, nil
}

// Del deletes links for a user with correlation_id.
func (s *LinksServer) Del(ctx context.Context, in *pb.ListShortenLinksToDelete) (*pb.Empty, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	var result pb.Empty

	go func() {
		for _, alias := range in.UserLinks {
			s.DB.DoDel(ctx, userID, []string{alias})
		}
	}()

	return &result, nil
}

// BatchShorten processes multiple URLs in a batch and returns their shortened versions.
func (s *LinksServer) BatchShorten(ctx context.Context, in *pb.BatchShortenRequest) (*pb.BatchShortenResponse, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	var response pb.BatchShortenResponse

	for _, item := range in.Items {
		if item.OriginalUrl == "" {
			continue
		}

		alias, err := s.DB.DoPut(ctx, item.OriginalUrl, "", userID)
		if err != nil {
			if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				// Skip URLs that fail due to integrity constraints
				continue
			}
			return nil, status.Errorf(codes.Internal, "Error shortening URL")
		}

		response.Items = append(response.Items, &pb.BatchShortenResponseItem{
			CorrelationId: item.CorrelationId,
			ShortUrl:      fmt.Sprintf("%s/%s", "http://localhost:3200", alias),
		})
	}

	return &response, nil
}

func getUserID(ctx context.Context) (int, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Errorf(codes.DataLoss, "No metadata")
	}

	values := md.Get("userID")
	if len(values) == 0 {
		return 0, status.Errorf(codes.Unauthenticated, "No userID in metadata")
	}

	id, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, status.Errorf(codes.Internal, "Invalid userID format")
	}

	return id, nil
}
