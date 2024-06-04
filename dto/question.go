package dto

import (
	"blog/model"
	"github.com/go-playground/validator/v10"
)

type QuestionDTO struct {
	Title   string `json:"title" form:"title" validate:"required,min=1,max=64"`
	Content string `json:"content" form:"content" validate:"required,min=1,max=1024"`
}

func (m *QuestionDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

func (m *QuestionDTO) ConveyToModel(question *model.Question, uid *uint) {
	question.OwnerID = uid
	question.Content = m.Content
	question.Title = m.Title
}

type QuestionListDTO struct {
	QuestionList []QuestionDTO `json:"question_list"`
	Paginate
}

func (m *QuestionListDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

type UpdateTitleDTO struct {
	QuestionID uint   `json:"question_id" form:"question_id" validate:"required,number"`
	Title      string `json:"title" form:"title" validate:"required,min=1,max=64"`
}

func (m *UpdateTitleDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

type UpdateContentDTO struct {
	QuestionID uint   `json:"question_id" form:"question_id" validate:"required,number"`
	Content    string `json:"content" form:"content" validate:"required,min=1,max=1028"`
}

func (m *UpdateContentDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

type UploadQuestionPictureDTO struct {
	QuestionID uint `json:"question_id" form:"question_id" validate:"required,number"`
}

func (m *UploadQuestionPictureDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

type DeleteQuestionDTO struct {
	QuestionID uint `json:"question_id" form:"question_id" validate:"required,number"`
}

func (m *DeleteQuestionDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

type DeleteQuestionPictureDTO struct {
	QuestionID uint `json:"question_id" form:"question_id" validate:"required,number"`
}

func (m *DeleteQuestionPictureDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

type AddQuestionSubscribeDTO struct {
	QuestionID uint `json:"question_id" form:"question_id" validate:"required,number"`
}

func (m *AddQuestionSubscribeDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

func (m *AddQuestionSubscribeDTO) ConveyToModel(AddQuestionSubscribe *model.Question, uid *uint) {
	AddQuestionSubscribe.OwnerID = uid

}

type GetMySubscribeDTO struct {
	QuestionID uint `json:"question_id" form:"question_id" validate:"required,number"`
}

func (m *GetMySubscribeDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

type DeleteQuestionSubscribeDTO struct {
	QuestionID uint `json:"question_id" form:"question_id" validate:"required,number"`
}

func (m *DeleteQuestionSubscribeDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

func (m *DeleteQuestionSubscribeDTO) ConveyToModel(DeleteQuestionSubscribe *model.Question, uid *uint) {
	DeleteQuestionSubscribe.OwnerID = uid

}
