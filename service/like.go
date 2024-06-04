package service

import "blog/dao"

var likeService *LikeService

type LikeService struct {
	BaseService
	Dao *dao.LikeDao
}

func NewLikeApi() *LikeService {
	if likeService == nil {
		likeService = &LikeService{
			Dao: dao.NewLikeDao(),
		}
	}
	return likeService
}
