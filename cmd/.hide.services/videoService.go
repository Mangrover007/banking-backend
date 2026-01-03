package services

import "github.com/Mangrover007/banking-backend/cmd/entities"

type VideoService interface {
	FindAll() []entities.Video
	Save(entities.Video) entities.Video
}

type videoService struct {
	videos []entities.Video
}

func New() VideoService {
	return &videoService{}
}

func (v *videoService) FindAll() []entities.Video {
	return v.videos
}

func (v *videoService) Save(video entities.Video) entities.Video {
	v.videos = append(v.videos, video)
	return video
}
