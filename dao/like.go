package dao

var likeDao *LikeDao

type LikeDao struct {
	BaseDao
}

func NewLikeDao() *LikeDao {
	if likeDao == nil {
		likeDao = &LikeDao{
			NewBaseDao(),
		}
	}
	return likeDao
}
