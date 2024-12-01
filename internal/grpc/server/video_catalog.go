package server

import (
	"context"

	cons "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/grpc/client/user"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/aws"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/database"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/log"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/model"
	usrpb "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/proto/user"
	vcpb "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/proto/video_catalog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm/clause"
)

type videoCatalogServer struct {
	vcpb.VideoCatalogServiceServer
}

func (v *videoCatalogServer) FindAll(ctx context.Context, data *vcpb.FindAllRequest) (*vcpb.FindAllResponse, error) {
	rows := []model.Video{}
	videos := []*vcpb.Video{}

	//@TODO: implement pagination
	err := database.DB.
		Select("id", "title", "description", "thumbnail", "published_at", "duration", "resolution").
		Order("created_at desc").
		Find(&rows).Error

	if err != nil {
		log.Error("Unable to query data: %v", err)

		return nil, status.Errorf(codes.Internal, cons.MessageInternalServerError)
	}

	for _, v := range rows {
		thumbnailUrl, err := aws.CreateGetObjectPresignedUploadUrl(v.Thumbnail)

		if err != nil {
			log.Error("Unable to create thumbnail url: %v", err)

			return nil, status.Errorf(codes.Internal, cons.MessageInternalServerError)
		}

		videos = append(videos, &vcpb.Video{
			Id:           int32(v.Id),
			Title:        v.Title,
			Description:  v.Description,
			ThumbnailUrl: thumbnailUrl,
			PublishedAt:  v.PublishedAt.String(),
			Duration:     int32(v.Duration),
			Resolution:   v.Resolution,
		})
	}

	response := &vcpb.FindAllResponse{
		Message: cons.MessageOK,
		Data: &vcpb.FindAllResponseData{
			Videos: videos,
		},
	}

	return response, nil
}

func (v *videoCatalogServer) FindById(ctx context.Context, data *vcpb.FindByIdRequest) (*vcpb.FindByIdResponse, error) {
	row := new(model.Video)

	err := database.DB.
		Select("id", "title", "description", "thumbnail", "published_at", "duration", "resolution").
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).
		First(&row).Error

	if err != nil {
		log.Error("Unable to query data: %v", err)

		return nil, status.Errorf(codes.NotFound, cons.MessageNotFound)
	}

	thumbnailUrl, err := aws.CreateGetObjectPresignedUploadUrl(row.Thumbnail)

	if err != nil {
		log.Error("Unable to create thumbnail url: %v", err)

		return nil, status.Errorf(codes.Internal, cons.MessageInternalServerError)
	}

	clientResponse, err := user.User.FindById(&usrpb.FindByIdRequest{
		Id: int32(row.Id),
	})

	if err != nil {
		return nil, err
	}

	user := clientResponse.Data.User

	response := &vcpb.FindByIdResponse{
		Message: cons.MessageOK,
		Data: &vcpb.FindByIdResponseData{
			Video: &vcpb.Video{
				Id:           int32(row.Id),
				Title:        row.Title,
				Description:  row.Description,
				ThumbnailUrl: thumbnailUrl,
				PublishedAt:  row.PublishedAt.String(),
				Duration:     int32(row.Duration),
				Resolution:   row.Resolution,
				User: &vcpb.User{
					Id:    user.Id,
					Name:  user.Name,
					Image: user.Image,
				},
			},
		},
	}

	return response, nil
}
