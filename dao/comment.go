package dao

import (
	"blog/dto"
	"blog/model"
	"errors"
	"os"
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

func (m *CommentDao) GetCommentDetail(CommentID uint, uid uint) (model.Comment, error) {
	var detail model.Comment
	err := m.Orm.Model(&model.Comment{}).Where("id = ?", CommentID).First(&detail).Error
	if err != nil {
		return model.Comment{}, err
	}
	if !m.CheckEditCommentPermission(uid, *detail.OwnerID) {
		return model.Comment{}, errors.New("no permission")
	}
	return detail, nil
}

func (m *CommentDao) CheckEditCommentPermission(uid uint, ownerID uint) bool {
	return uid == ownerID
}
func (m *CommentDao) GetMyLikeCount(commentID uint) (int64, error) {
	var comment model.Comment
	err := m.Orm.Model(&model.Comment{}).Where("id = ?", commentID).First(&comment).Error
	if err != nil {
		return 0, err
	}
	return comment.LikeCount, nil
}

func (m *CommentDao) CheckQuestionIdExist(questionID uint) error {
	var question model.Question
	err := m.Orm.Model(&model.Question{}).Where("id = ?", questionID).First(&question).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *CommentDao) CheckCommentExist(uid uint) error {
	var comment model.Comment
	err := m.Orm.Model(&model.Comment{}).Where("id = ?", uid).First(&comment).Error
	if err != nil {
		return err
	}
	if comment.PictureID != nil {
		var pic model.UploadFile
		err = m.Orm.Model(&model.UploadFile{}).Where("id = ?", comment.PictureID).First(&pic).Error
		if err != nil {
			return err
		}
		err = os.Remove(pic.FilePath)
		if err != nil {
			return err
		}
		err = m.Orm.Model(&model.UploadFile{}).Where("id = ?", pic.ID).Delete(&pic).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *CommentDao) SetCommentPicture(fileID uint, commentID uint) error {
	return m.Orm.Model(&model.Comment{}).Where("id = ?", commentID).Update("picture_id", fileID).Error
}

func (m *CommentDao) DeleteCommentPicture(commentID uint) error {
	var iFile model.UploadFile
	err := m.Orm.Model(&model.UploadFile{}).Where("id = ?", commentID).Delete(&iFile).Error
	if err != nil {
		return err
	}
	return nil
}
