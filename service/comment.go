package service

import (
	"blog/dao"
	"blog/dto"
	"blog/model"
	"blog/utils"
)

var commentService *CommentService

type CommentService struct {
	BaseService
	Dao *dao.CommentDao
}

func NewCommentService() *CommentService {
	if commentService == nil {
		commentService = &CommentService{
			Dao: dao.NewCommentDao(),
		}
	}
	return commentService
}

func (m *CommentService) GetContentList(iCommentList *dto.CommentListDTO) ([]model.Comment, int64, error) {
	if err := iCommentList.Validate(); err != nil {
		return nil, 0, utils.ParseValidateError(err)
	}
	if iCommentList.PageNum <= 0 {
		iCommentList.PageNum = 1
	}
	if iCommentList.PageSize <= 0 {
		iCommentList.PageSize = 10
	}
	list, nTotal, err := m.Dao.GetCommentList(iCommentList)
	if err != nil {
		return nil, 0, err
	}
	return list, nTotal, nil
}

func (m *CommentService) PostComment(iPostComment *dto.PostComment, iUserAuthDTO *dto.UserAuthDTO) (model.Comment, error) {
	if err := iPostComment.Validate(); err != nil {
		return model.Comment{}, utils.ParseValidateError(err)
	}

	comment, err := m.Dao.PostComment(iPostComment, iUserAuthDTO.Uid)
	if err != nil {
		return model.Comment{}, err
	}
	return comment, nil
}

func (m *CommentService) CommentPostComment(iComment *dto.CommentPostComment, iUserAuthDTO *dto.UserAuthDTO) (model.Comment, error) {
	if err := iComment.Validate(); err != nil {
		return model.Comment{}, utils.ParseValidateError(err)
	}

	comment, err := m.Dao.CommentPostComment(iComment, iUserAuthDTO.Uid)
	if err != nil {
		return model.Comment{}, err
	}
	return comment, nil
}

func (m *CommentService) GetMyComment(iGetMyCommentList *dto.CommentListDTO, iUserAuth *dto.UserAuthDTO) ([]model.Comment, int64, error) {
	var list []model.Comment
	var nTotal int64
	if iGetMyCommentList.PageNum <= 0 {
		iGetMyCommentList.PageNum = 1
	}
	if iGetMyCommentList.PageSize <= 0 {
		iGetMyCommentList.PageSize = 10
	}
	list, nTotal, err := m.Dao.GetMyComment(iGetMyCommentList, iUserAuth.Uid)
	if err != nil {
		return nil, 0, err
	}
	return list, nTotal, nil
}

func (m *CommentService) DeleteComment(iDeleteCommentDTO *dto.DeleteCommentDTO, iUserAuthDTO *dto.UserAuthDTO) error {
	if err := iDeleteCommentDTO.Validate(); err != nil {
		return utils.ParseValidateError(err)
	}
	return m.Dao.DeleteComment(iDeleteCommentDTO, iUserAuthDTO.Uid)
}
