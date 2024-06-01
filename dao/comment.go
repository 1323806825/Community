package dao

import (
	"blog/dto"
	"blog/model"
)

var commentDao *CommentDao

type CommentDao struct {
	BaseDao
}

func NewCommentDao() *CommentDao {
	if commentDao == nil {
		commentDao = &CommentDao{
			NewBaseDao(),
		}
	}
	return commentDao
}

func (m *CommentDao) GetCommentList(iGetCommentList *dto.CommentListDTO) ([]model.Comment, int64, error) {
	var iCommentList []model.Comment
	var nTotal int64
	err := m.Orm.Model(&model.Comment{}).
		Where("question_id = ?", iGetCommentList.QuestionID).
		Scopes(m.Paginate(iGetCommentList.Paginate)).
		Order("created_at desc").
		Find(&iCommentList).
		Offset(-1).Limit(-1).
		Count(&nTotal).
		Error
	if err != nil {
		return nil, 0, err
	}
	return iCommentList, nTotal, nil
}

func (m *CommentDao) PostComment(iPostComment *dto.PostComment, uid uint) (model.Comment, error) {
	var iComment model.Comment
	iPostComment.ConveyToModel(&iComment, &uid)
	err := m.Orm.Save(&iComment).Error
	if err != nil {
		return model.Comment{}, err
	}
	return iComment, nil
}

func (m *CommentDao) CommentPostComment(iCommentPostComment *dto.CommentPostComment, uid uint) (model.Comment, error) {
	var iComment model.Comment
	iCommentPostComment.ConveyToModel(&iComment, &uid)
	err := m.Orm.Model(&model.Comment{}).Save(&iComment).Error
	if err != nil {
		return model.Comment{}, err
	}
	return iComment, nil
}

func (m *CommentDao) GetMyComment(iGetMyComment *dto.CommentListDTO, uid uint) ([]model.Comment, int64, error) {
	var list []model.Comment
	var nTotal int64
	err := m.Orm.Model(&model.Comment{}).Where("owner_id = ?", uid).
		Scopes(m.Paginate(iGetMyComment.Paginate)).
		Order("created_at desc").
		Find(&list).
		Offset(-1).Limit(-1).
		Count(&nTotal).Error
	if err != nil {
		return nil, 0, err
	}
	return list, nTotal, nil
}

func (m *CommentDao) DeleteComment(iDeleteCommentDTO *dto.DeleteCommentDTO, uid uint) error {
	var iComment model.Comment
	iDeleteCommentDTO.ConveyToModel(&iComment, &uid)
	err := m.Orm.Model(&model.Comment{}).Where("id = ?", iDeleteCommentDTO.CommentID).Delete(&iComment).Error
	return err
}
