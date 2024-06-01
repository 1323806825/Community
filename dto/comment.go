package dto

import (
	"blog/model"
	"github.com/go-playground/validator/v10"
)

type CommentDTO struct {
	Content string `json:"content" form:"content" validate:"required,min=1,max=1024"`
}

type CommentListDTO struct {
	QuestionID uint `json:"question_id" form:"question_id" validate:"required,number"`
	Paginate
	CommentList []CommentDTO `json:"comment_list"`
}

func (m *CommentListDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

type PostComment struct {
	QuestionID uint   `json:"question_id" form:"question_id" validate:"required,number"`
	Content    string `json:"content" form:"content" validate:"required,min=1,max=2048"`
}

func (m *PostComment) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

func (m *PostComment) ConveyToModel(comment *model.Comment, uid *uint) {
	comment.QuestionID = m.QuestionID
	comment.Content = m.Content
	comment.OwnerID = uid
}

type CommentPostComment struct {
	ParentCommentID uint   `json:"parent_comment_id" form:"parent_comment_id" validate:"required,number"`
	Content         string `json:"content" form:"content" validate:"required,min=1,max=2048"`
}

func (m *CommentPostComment) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

func (m *CommentPostComment) ConveyToModel(iComment *model.Comment, uid *uint) {
	iComment.ParentCommentID = m.ParentCommentID
	iComment.OwnerID = uid
	iComment.Content = m.Content
}

type DeleteCommentDTO struct {
	CommentID uint `json:"comment_id" form:"comment_id" validate:"required,number"`
}

func (m *DeleteCommentDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

func (m *DeleteCommentDTO) ConveyToModel(iComment *model.Comment, uid *uint) {
	iComment.QuestionID = m.CommentID
	iComment.OwnerID = uid
}
