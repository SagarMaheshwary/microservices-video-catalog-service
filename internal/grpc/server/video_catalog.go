package server

import (
	"context"
	"path"

	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/config"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/grpc/client/user"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/aws"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/database"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-video-catalog-service/internal/model"
	userpb "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/proto/user"
	videocatalogpb "github.com/sagarmaheshwary/microservices-video-catalog-service/internal/proto/video_catalog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type videoCatalogServer struct {
	videocatalogpb.VideoCatalogServiceServer
}

func (v *videoCatalogServer) FindAll(ctx context.Context, data *videocatalogpb.FindAllRequest) (*videocatalogpb.FindAllResponse, error) {
	rows := []model.Video{}
	videos := []*videocatalogpb.Video{}

	//@TODO: implement pagination
	err := database.Conn.
		Select("id", "title", "description", "thumbnail", "published_at", "duration", "resolution").
		Order("created_at desc").
		Find(&rows).Error

	if err != nil {
		logger.Error("Unable to query data: %v", err)

		return nil, status.Errorf(codes.Internal, constant.MessageInternalServerError)
	}

	for _, v := range rows {
		thumbnailURL, err := aws.CreateGetObjectPresignedUploadUrl(v.Thumbnail)

		if err != nil {
			logger.Error("Unable to create thumbnail url: %v", err)

			return nil, status.Errorf(codes.Internal, constant.MessageInternalServerError)
		}

		videos = append(videos, &videocatalogpb.Video{
			Id:           int32(v.Id),
			Title:        v.Title,
			Description:  v.Description,
			ThumbnailUrl: thumbnailURL,
			PublishedAt:  v.PublishedAt.String(),
			Duration:     int32(v.Duration),
			Resolution:   v.Resolution,
		})
	}

	response := &videocatalogpb.FindAllResponse{
		Message: constant.MessageOK,
		Data: &videocatalogpb.FindAllResponseData{
			Videos: videos,
		},
	}

	return response, nil
}

func (v *videoCatalogServer) FindById(ctx context.Context, data *videocatalogpb.FindByIdRequest) (*videocatalogpb.FindByIdResponse, error) {
	video := new(model.Video)

	err := database.Conn.
		Select("id", "title", "description", "thumbnail", "path", "published_at", "duration", "resolution").
		Where(&model.Video{Id: uint(data.Id)}).
		First(&video).Error

	if err != nil {
		logger.Error("Unable to query data: %v", err)

		return nil, status.Errorf(codes.NotFound, constant.MessageNotFound)
	}

	u, err := user.User.FindById(&userpb.FindByIdRequest{
		Id: int32(video.Id),
	})

	if err != nil {
		return nil, err
	}

	thumbnailURL, err := aws.CreateGetObjectPresignedUploadUrl(video.Thumbnail)

	if err != nil {
		logger.Error("Unable to create thumbnail url: %v", err)

		return nil, status.Errorf(codes.Internal, constant.MessageInternalServerError)
	}

	manifestURL := path.Join(config.Conf.AWS.CloudFrontURL, video.Path, constant.MPEGDASHManifestFile)

	response := &videocatalogpb.FindByIdResponse{
		Message: constant.MessageOK,
		Data: &videocatalogpb.FindByIdResponseData{
			ManifestUrl: manifestURL,
			Video: &videocatalogpb.Video{
				Id:           int32(video.Id),
				Title:        video.Title,
				Description:  video.Description,
				ThumbnailUrl: thumbnailURL,
				PublishedAt:  video.PublishedAt.String(),
				Duration:     int32(video.Duration),
				Resolution:   video.Resolution,
				User: &videocatalogpb.User{
					Id:    u.Data.User.Id,
					Name:  u.Data.User.Name,
					Image: u.Data.User.Image,
				},
			},
		},
	}

	return response, nil
}
